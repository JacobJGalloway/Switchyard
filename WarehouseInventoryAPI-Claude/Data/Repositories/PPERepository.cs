using Microsoft.EntityFrameworkCore;
using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Data.Repositories
{
    public class PPERepository(InventoryContext context) : IPPERepository
    {
        private readonly InventoryContext _context = context;

        public async Task<IEnumerable<PPE>> GetAllAsync()
        {
            return await _context.PPE.ToListAsync();
        }

        public async Task<List<PPE>> GetBySKUIdAsync(string skuId)
        {
            return await _context.GetPPEBySKUIdsync(skuId);
        }

        public async Task<List<PPE>> GetByLocationAsync(string locationId)
        {
            return await _context.PPE
                .Where(p => p.LocationId == locationId)
                .ToListAsync();
        }

        public async Task<List<PPE>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await _context.PPE
                .Where(p => p.LocationId == locationId && p.SKUMarker == skuId)
                .ToListAsync();
        }

        public async Task<PPE> AddAsync(PPE item)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            item.LocationId = warehouseId;
            item.PartitionKey = $"{warehouseId}-{item.SKUMarker}-{Guid.NewGuid():N}";
            _context.PPE.Add(item);
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, PPE item)
        {
            var existingItems = await GetBySKUIdAsync(skuId);
            if (existingItems.Count == 0) return;

            var target = existingItems.FirstOrDefault(c => c.PartitionKey == item.PartitionKey)
                         ?? existingItems[0];

            target.RowKey = item.RowKey;
            target.SKUMarker = item.SKUMarker;
            target.UnloadedDate = item.UnloadedDate;
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            var target = await _context.PPE.FindAsync(partitionKey);
            if (target is null) return;

            if (projected.HasValue) target.Projected = projected.Value;
            if (unloadedDate.HasValue) target.UnloadedDate = unloadedDate.Value;
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            List<PPE> items = await _context.PPE.Where(p => p.SKUMarker == skuId).ToListAsync();
            if (items.Count == 0) return false;
            _context.PPE.RemoveRange(items);
            return true;
        }
        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var target = await _context.PPE.FindAsync(partitionKey);
            if (target is null) return false;

            _context.PPE.Remove(target);
            return true;
        }
    }
}