using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data.Interfaces;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data.Repositories
{
    public class ClothingRepository(InventoryContext writeContext, InventoryReadContext readContext) : IClothingRepository
    {
        public async Task<IEnumerable<Clothing>> GetAllAsync()
        {
            return await readContext.Clothing.AsNoTracking().ToListAsync();
        }

        public async Task<List<Clothing>> GetBySKUIdAsync(string skuId)
        {
            return await readContext.Clothing
                .Where(c => c.SKUMarker == skuId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<List<Clothing>> GetByLocationAsync(string locationId)
        {
            return await readContext.Clothing
                .Where(c => c.LocationId == locationId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<List<Clothing>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await readContext.Clothing
                .Where(c => c.LocationId == locationId && c.SKUMarker == skuId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<Clothing> AddAsync(Clothing item)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            item.LocationId = warehouseId;
            item.PartitionKey = $"{warehouseId}-{item.SKUMarker}-{Guid.NewGuid():N}";
            writeContext.Clothing.Add(item);
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Clothing item)
        {
            var existingItems = await writeContext.Clothing
                .Where(c => c.SKUMarker == skuId)
                .ToListAsync();
            if (existingItems.Count == 0) return;

            var target = existingItems.FirstOrDefault(c => c.PartitionKey == item.PartitionKey)
                         ?? existingItems[0];

            target.RowKey = item.RowKey;
            target.SKUMarker = item.SKUMarker;
            target.UnloadedDate = item.UnloadedDate;
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            var target = await writeContext.Clothing.FindAsync(partitionKey);
            if (target is null) return;

            if (projected.HasValue) target.Projected = projected.Value;
            if (unloadedDate.HasValue) target.UnloadedDate = unloadedDate.Value;
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            var items = await writeContext.Clothing
                .Where(c => c.SKUMarker == skuId)
                .ToListAsync();
            if (items.Count == 0) return false;
            writeContext.Clothing.RemoveRange(items);
            return true;
        }

        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var target = await writeContext.Clothing.FindAsync(partitionKey);
            if (target is null) return false;
            writeContext.Clothing.Remove(target);
            return true;
        }
    }
}
