using Microsoft.EntityFrameworkCore;
using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Data.Repositories
{
    public class ToolRepository(InventoryContext context) : IToolRepository
    {
        private readonly InventoryContext _context = context;

        public async Task<IEnumerable<Tool>> GetAllAsync()
        {
            return await _context.Tool.ToListAsync();
        }

        public async Task<List<Tool>> GetBySKUIdAsync(string skuId)
        {
            return await _context.GetToolBySKUIdsync(skuId);
        }

        public async Task<List<Tool>> GetByLocationAsync(string locationId)
        {
            return await _context.Tool
                .Where(t => t.LocationId == locationId)
                .ToListAsync();
        }

        public async Task<List<Tool>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await _context.Tool
                .Where(t => t.LocationId == locationId && t.SKUMarker == skuId)
                .ToListAsync();
        }

        public async Task<Tool> AddAsync(Tool item)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            item.LocationId = warehouseId;
            item.PartitionKey = $"{warehouseId}-{item.SKUMarker}-{Guid.NewGuid():N}";
            _context.Tool.Add(item);
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Tool item)
        {
            var existingItems = await GetBySKUIdAsync(skuId);
            if (existingItems.Count == 0) return;

            var target = existingItems.FirstOrDefault(t => t.PartitionKey == item.PartitionKey)
                         ?? existingItems[0];

            target.RowKey = item.RowKey;
            target.SKUMarker = item.SKUMarker;
            target.UnloadedDate = item.UnloadedDate;
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            var target = await _context.Tool.FindAsync(partitionKey);
            if (target is null) return;

            if (projected.HasValue) target.Projected = projected.Value;
            if (unloadedDate.HasValue) target.UnloadedDate = unloadedDate.Value;
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            List<Tool> items = await _context.Tool.Where(t => t.SKUMarker == skuId).ToListAsync();
            if (items.Count == 0) return false;
            _context.Tool.RemoveRange(items);
            return true;
        }
        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var target = await _context.Tool.FindAsync(partitionKey);
            if (target is null) return false;

            _context.Tool.Remove(target);
            return true;
        }
    }
}