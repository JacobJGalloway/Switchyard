using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data.Interfaces;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data.Repositories
{
    public class PPERepository(InventoryContext writeContext, InventoryReadContext readContext) : IPPERepository
    {
        public async Task<IEnumerable<PPE>> GetAllAsync()
        {
            return await readContext.PPE.AsNoTracking().ToListAsync();
        }

        public async Task<List<PPE>> GetBySKUIdAsync(string skuId)
        {
            return await readContext.PPE
                .Where(p => p.SKUMarker == skuId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<List<PPE>> GetByLocationAsync(string locationId)
        {
            return await readContext.PPE
                .Where(p => p.LocationId == locationId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<List<PPE>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await readContext.PPE
                .Where(p => p.LocationId == locationId && p.SKUMarker == skuId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<PPE> AddAsync(PPE item)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            item.LocationId = warehouseId;
            item.PartitionKey = $"{warehouseId}-{item.SKUMarker}-{Guid.NewGuid():N}";
            writeContext.PPE.Add(item);
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, PPE item)
        {
            var existingItems = await writeContext.PPE
                .Where(p => p.SKUMarker == skuId)
                .ToListAsync();
            if (existingItems.Count == 0) return;

            var target = existingItems.FirstOrDefault(p => p.PartitionKey == item.PartitionKey)
                         ?? existingItems[0];

            target.RowKey = item.RowKey;
            target.SKUMarker = item.SKUMarker;
            target.UnloadedDate = item.UnloadedDate;
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            var target = await writeContext.PPE.FindAsync(partitionKey);
            if (target is null) return;

            if (projected.HasValue) target.Projected = projected.Value;
            if (unloadedDate.HasValue) target.UnloadedDate = unloadedDate.Value;
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            var items = await writeContext.PPE
                .Where(p => p.SKUMarker == skuId)
                .ToListAsync();
            if (items.Count == 0) return false;
            writeContext.PPE.RemoveRange(items);
            return true;
        }

        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var target = await writeContext.PPE.FindAsync(partitionKey);
            if (target is null) return false;
            writeContext.PPE.Remove(target);
            return true;
        }
    }
}
