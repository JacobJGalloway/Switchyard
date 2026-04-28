using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Repositories
{
    public class StoreRepository(LogisticsReadContext readContext) : IStoreRepository
    {
        public async Task<IEnumerable<Store>> GetAllAsync()
        {
            return await readContext.Stores.AsNoTracking().ToListAsync();
        }
    }
}
