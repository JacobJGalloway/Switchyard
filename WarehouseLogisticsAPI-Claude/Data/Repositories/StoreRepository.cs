using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Repositories
{
    public class StoreRepository(LogisticsContext context) : IStoreRepository
    {
        private readonly LogisticsContext _context = context;

        public async Task<IEnumerable<Store>> GetAllAsync()
        {
            return await _context.Stores.ToListAsync();
        }
    }
}
