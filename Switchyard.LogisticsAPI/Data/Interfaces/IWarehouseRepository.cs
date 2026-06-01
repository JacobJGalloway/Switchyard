using Switchyard.Domain;

namespace Switchyard.LogisticsAPI.Data.Interfaces
{
    public interface IWarehouseRepository
    {
        Task<IEnumerable<Warehouse>> GetAllAsync();
    }
}
