using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Interfaces
{
    public interface IWarehouseRepository
    {
        Task<IEnumerable<Warehouse>> GetAllAsync();
    }
}
