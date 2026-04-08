using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Services.Interfaces;

namespace WarehouseLogistics_Claude.Services
{
    public class BillOfLadingService(IUnitOfWork unitOfWork) : IBillOfLadingService
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        public async Task<string> CreateAsync(BillOfLading billOfLading)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            if (billOfLading.Status != "Pending")
                throw new ArgumentException("A Pending Bill of Lading is required.");

            if (string.IsNullOrWhiteSpace(billOfLading.CustomerFirstName) ||
                string.IsNullOrWhiteSpace(billOfLading.CustomerLastName) ||
                string.IsNullOrWhiteSpace(billOfLading.City) ||
                string.IsNullOrWhiteSpace(billOfLading.State))
                throw new ArgumentException("Customer information is required.");

            if (billOfLading.LineEntries == null || !billOfLading.LineEntries.Any())
                throw new ArgumentException("At least one line entry is required.");

            if (string.IsNullOrWhiteSpace(warehouseId))
                throw new ArgumentException("Originating warehouse ID is not configured.");

            billOfLading.TransactionId = Guid.NewGuid().ToString()[..8].Replace("-", "");
            billOfLading.Status = "Submitted";
            billOfLading.PartitionKey = $"{billOfLading.TransactionId}-{Guid.NewGuid().ToString().Replace("-", "")}";

            await _unitOfWork.BillsOfLading.AddAsync(billOfLading);

            await PersistLineEntriesAsync(billOfLading.TransactionId, billOfLading.LineEntries);

            await _unitOfWork.SaveChangesAsync();

            await WriteDocumentAsync(billOfLading);
            await ProcessLocationStop(billOfLading.TransactionId, warehouseId);

            return billOfLading.TransactionId;
        }

        public async Task ProcessLocationStop(string transactionId, string locationId)
        {
            var lineEntries = await _unitOfWork.LineEntries.GetLineEntriesByTransactionIdAsync(transactionId);
            var locationLineEntries = lineEntries.Where(le => le.LocationId == locationId).ToList();

            if (locationLineEntries.Count == 0) return;

            foreach (var lineEntry in locationLineEntries)
            {
                lineEntry.IsProcessed = true;
                await _unitOfWork.LineEntries.UpdateLineEntryAsync(lineEntry);
            }

            await _unitOfWork.SaveChangesAsync();
        }

        public async Task<IEnumerable<BillOfLading>> GetAllAsync()
            => await _unitOfWork.BillsOfLading.GetAllAsync();

        public async Task<BillOfLading?> GetByTransactionIdAsync(string transactionId)
        {
            var bol = await _unitOfWork.BillsOfLading.GetByTransactionIdAsync(transactionId);
            if (bol is null) return null;
            bol.LineEntries = await _unitOfWork.LineEntries.GetLineEntriesByTransactionIdAsync(transactionId);
            return bol;
        }

        public async Task<List<LineEntry>> GetLineEntriesByTransactionIdAsync(string transactionId)
            => await _unitOfWork.LineEntries.GetLineEntriesByTransactionIdAsync(transactionId);

        public async Task ReplaceLocationStopAsync(string transactionId, string oldLocationId, string newLocationId)
        {
            var allEntries = await _unitOfWork.LineEntries.GetLineEntriesByTransactionIdAsync(transactionId);
            var oldEntries = allEntries.Where(le => le.LocationId == oldLocationId).ToList();

            if (oldEntries.Count == 0)
                throw new ArgumentException($"No line entries found for location {oldLocationId} on transaction {transactionId}.");

            if (oldEntries.Any(le => le.IsProcessed))
                throw new InvalidOperationException($"Location {oldLocationId} has already been processed and cannot be replaced.");

            foreach (var entry in oldEntries)
            {
                var newEntry = new LineEntry
                {
                    TransactionId = transactionId,
                    LocationId    = newLocationId,
                    SKUMarker     = entry.SKUMarker,
                    Quantity      = entry.Quantity,
                };
                await _unitOfWork.LineEntries.AddAsync(newEntry);
            }

            await _unitOfWork.LineEntries.DeleteByLocationAsync(transactionId, oldLocationId);
            await _unitOfWork.SaveChangesAsync();
        }

        private async Task PersistLineEntriesAsync(string transactionId, List<LineEntry> lineEntries)
        {
            foreach (var lineEntry in lineEntries)
            {
                lineEntry.TransactionId = transactionId;
                await _unitOfWork.LineEntries.AddAsync(lineEntry);
            }
        }

        private static async Task WriteDocumentAsync(BillOfLading billOfLading)
        {
            const int lineWidth = 60;
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;

            string centeredWarehouseIdLine = string.Empty;
            if (!string.IsNullOrWhiteSpace(warehouseId))
            {
                string warehouseIdLine = $"Warehouse ID: {warehouseId}";
                centeredWarehouseIdLine = warehouseId[..2] == "WH"
                    ? warehouseIdLine.PadLeft((lineWidth + warehouseIdLine.Length) / 2).PadRight(lineWidth)
                    : string.Empty;
            }

            string transactionLine = $"Transaction ID: {billOfLading.TransactionId}";
            string centeredTransactionLine = transactionLine.PadLeft((lineWidth + transactionLine.Length) / 2).PadRight(lineWidth);

            int colWarehouse = Math.Max("Warehouse".Length,  billOfLading.LineEntries.Max(e => e.LocationId.Length));
            int colSKU       = Math.Max("SKU Marker".Length, billOfLading.LineEntries.Max(e => e.SKUMarker.Length));
            int colQty       = Math.Max("Quantity".Length,   billOfLading.LineEntries.Max(e => e.Quantity.ToString().Length));

            string tableTop     = $"┌{new string('─', colWarehouse + 2)}┬{new string('─', colSKU + 2)}┬{new string('─', colQty + 2)}┐";
            string tableHeader  = $"│ {"Warehouse".PadRight(colWarehouse)} │ {"SKU Marker".PadRight(colSKU)} │ {"Quantity".PadRight(colQty)} │";
            string tableDivider = $"├{new string('─', colWarehouse + 2)}┼{new string('─', colSKU + 2)}┼{new string('─', colQty + 2)}┤";
            string tableBottom  = $"└{new string('─', colWarehouse + 2)}┴{new string('─', colSKU + 2)}┴{new string('─', colQty + 2)}┘";

            var originRows = string.IsNullOrWhiteSpace(warehouseId)
                ? []
                : billOfLading.LineEntries
                    .Where(e => string.Equals(e.LocationId, warehouseId))
                    .Select(e => $"│ {e.LocationId.PadRight(colWarehouse)} │ {e.SKUMarker.PadRight(colSKU)} │ {e.Quantity.ToString().PadRight(colQty)} │")
                    .ToList();

            var remainingRows = billOfLading.LineEntries
                .Where(e => !string.Equals(e.LocationId, warehouseId))
                .OrderBy(e => e.LocationId)
                .Select(e => $"│ {e.LocationId.PadRight(colWarehouse)} │ {e.SKUMarker.PadRight(colSKU)} │ {e.Quantity.ToString().PadRight(colQty)} │")
                .ToList();

            var fileLines = new[]
            {
                centeredWarehouseIdLine,
                centeredTransactionLine,
                $"Client: {billOfLading.CustomerFirstName} {billOfLading.CustomerLastName}",
                $"City: {billOfLading.City} State: {billOfLading.State}",
                string.Empty,
                new string('─', lineWidth),
                Environment.GetEnvironmentVariable("LINE_ENTRIES_WORDING") ?? "Line Entries",
                new string('─', lineWidth),
                string.Empty,
                tableTop,
                tableHeader,
                tableDivider
            }
            .Concat(originRows)
            .Concat(remainingRows)
            .Append(tableBottom);

            string downloadsPath = Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.UserProfile), "Downloads");
            await File.WriteAllLinesAsync(Path.Combine(downloadsPath, $"BillOfLading_{billOfLading.TransactionId}.txt"), fileLines);
        }
    }
}
