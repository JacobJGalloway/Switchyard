using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Interfaces
{
    public interface IStoreRepository
    {
        Task<IEnumerable<Store>> GetAllAsync();
    }
}
