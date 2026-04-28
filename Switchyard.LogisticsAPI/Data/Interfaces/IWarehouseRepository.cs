using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Interfaces
{
    public interface IWarehouseRepository
    {
        Task<IEnumerable<Warehouse>> GetAllAsync();
    }
}
