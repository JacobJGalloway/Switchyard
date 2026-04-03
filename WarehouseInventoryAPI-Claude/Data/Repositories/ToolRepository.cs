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
            return await _context.Tools.ToListAsync();
        }

        public async Task<Tool?> GetBySKUIdAsync(string skuId)
        {
            return await _context.Tools.FindAsync(skuId);
        }

        public async Task<Tool> AddAsync(Tool item)
        {
            _context.Tools.Add(item);
            await _context.SaveChangesAsync();
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Tool item)
        {
            var existingItem = await GetBySKUIdAsync(skuId);
            if (existingItem is null) return;

            _context.Entry(existingItem).CurrentValues.SetValues(item);
            await _context.SaveChangesAsync();
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            List<Tool> items = await _context.Tools.Where(t => t.PartitionKey == skuId).ToListAsync();
            if (items.Count == 0) return false;
            _context.Tools.RemoveRange(items);
            await _context.SaveChangesAsync();
            return true;
        }
    }
}

