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

        public async Task<PPE?> GetBySKUIdAsync(string skuId)
        {
            return await _context.PPE.FindAsync(skuId);
        }

        public async Task<PPE> AddAsync(PPE item)
        {
            _context.PPE.Add(item);
            await _context.SaveChangesAsync();
            return item;
        }

        public async Task UpdateBySKUIdAsync(string skuId, PPE item)
        {
            var existingItem = await GetBySKUIdAsync(skuId);
            if (existingItem is null) return;

            _context.Entry(existingItem).CurrentValues.SetValues(item);
            await _context.SaveChangesAsync();
        }

        public async Task<bool> DeleteBySKUIdAsync(string skuId)
        {
            List<PPE> items = await _context.PPE.Where(p => p.PartitionKey == skuId).ToListAsync();
            if (items.Count == 0) return false;
            _context.PPE.RemoveRange(items);
            await _context.SaveChangesAsync();
            return true;
        }
    }
}

