using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Repositories
{
    public class WarehouseRepository(LogisticsContext context) : IWarehouseRepository
    {
        private readonly LogisticsContext _context = context;

        public async Task<IEnumerable<Warehouse>> GetAllAsync()
        {
            return await _context.Warehouses.ToListAsync();
        }
    }
}
