using Switchyard.Domain;

namespace Switchyard.LogisticsAPI.Data.Interfaces
{
    public interface IStoreRepository
    {
        Task<IEnumerable<Store>> GetAllAsync();
    }
}
