using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Interfaces
{
    public interface IStoreRepository
    {
        Task<IEnumerable<Store>> GetAllAsync();
    }
}
