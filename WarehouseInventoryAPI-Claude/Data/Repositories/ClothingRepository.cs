using Microsoft.EntityFrameworkCore;
using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Data.Repositories
{
    public class ClothingRepository(InventoryContext context) : IClothingRepository
    {
        private readonly InventoryContext _context = context;

        public async Task<IEnumerable<Clothing>> GetAllAsync()
        {
            return await _context.Clothing.ToListAsync();
        }

        public async Task<List<Clothing>> GetBySKUIdAsync(string skuId)
        {
            var items = await _context.GetClothingBySKUIdsync(skuId);
            
            if (items.Count == 0) return new List<Clothing>();

            return items;
        }

        public async Task<List<Clothing>> GetByLocationAsync(string locationId)
        {
            return await _context.Clothing
                .Where(c => c.LocationId == locationId)
                .ToListAsync();
        }

        public async Task<List<Clothing>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await _context.Clothing
                .Where(c => c.LocationId == locationId && c.SKUMarker == skuId)
                .ToListAsync();
        }

        public async Task<Clothing> AddAsync(Clothing item)
        {
            var warehouseId = Environment.GetEnvironmentVariable("LOCATION_ID") ?? string.Empty;
            item.LocationId = warehouseId;
            item.PartitionKey = $"{warehouseId}-{item.SKUMarker}-{Guid.NewGuid():N}";
            _context.Clothing.Add(item);
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Clothing item)
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
            var target = await _context.Clothing.FindAsync(partitionKey);
            if (target is null) return;

            if (projected.HasValue) target.Projected = projected.Value;
            if (unloadedDate.HasValue) target.UnloadedDate = unloadedDate.Value;
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            List<Clothing> items = await _context.Clothing.Where(c => c.SKUMarker == skuId).ToListAsync();
            if (items.Count == 0) return false;
            _context.Clothing.RemoveRange(items);
            return true;
        }
        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var target = await _context.Clothing.FindAsync(partitionKey);
            if (target is null) return false;

            _context.Clothing.Remove(target);
            return true;
        }
    }
}