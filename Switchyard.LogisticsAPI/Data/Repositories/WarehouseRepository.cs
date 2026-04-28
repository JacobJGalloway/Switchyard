using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Repositories
{
    public class WarehouseRepository(LogisticsReadContext readContext) : IWarehouseRepository
    {
        public async Task<IEnumerable<Warehouse>> GetAllAsync()
        {
            return await readContext.Warehouses.AsNoTracking().ToListAsync();
        }
    }
}
