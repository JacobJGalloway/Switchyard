using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data.Interfaces;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data.Repositories
{
    public class ToolRepository(InventoryContext writeContext, InventoryReadContext readContext) : IToolRepository
    {
        public async Task<IEnumerable<Tool>> GetAllAsync()
        {
            return await readContext.Tool.AsNoTracking().ToListAsync();
        }

        public async Task<List<Tool>> GetBySKUIdAsync(string skuId)
        {
            return await readContext.Tool
                .Where(t => t.SKUMarker == skuId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<List<Tool>> GetByLocationAsync(string locationId)
        {
            return await readContext.Tool
                .Where(t => t.LocationId == locationId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<List<Tool>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await readContext.Tool
                .Where(t => t.LocationId == locationId && t.SKUMarker == skuId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<Tool> AddAsync(Tool item)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            item.LocationId = warehouseId;
            item.PartitionKey = $"{warehouseId}-{item.SKUMarker}-{Guid.NewGuid():N}";
            writeContext.Tool.Add(item);
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Tool item)
        {
            var existingItems = await writeContext.Tool
                .Where(t => t.SKUMarker == skuId)
                .ToListAsync();
            if (existingItems.Count == 0) return;

            var target = existingItems.FirstOrDefault(t => t.PartitionKey == item.PartitionKey)
                         ?? existingItems[0];

            target.RowKey = item.RowKey;
            target.SKUMarker = item.SKUMarker;
            target.UnloadedDate = item.UnloadedDate;
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            var target = await writeContext.Tool.FindAsync(partitionKey);
            if (target is null) return;

            if (projected.HasValue) target.Projected = projected.Value;
            if (unloadedDate.HasValue) target.UnloadedDate = unloadedDate.Value;
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            var items = await writeContext.Tool
                .Where(t => t.SKUMarker == skuId)
                .ToListAsync();
            if (items.Count == 0) return false;
            writeContext.Tool.RemoveRange(items);
            return true;
        }

        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var target = await writeContext.Tool.FindAsync(partitionKey);
            if (target is null) return false;
            writeContext.Tool.Remove(target);
            return true;
        }
    }
}
