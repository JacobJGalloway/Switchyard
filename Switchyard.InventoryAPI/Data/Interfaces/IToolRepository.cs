using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data.Interfaces
{
    public interface IToolRepository
    {
        Task<IEnumerable<Tool>> GetAllAsync();
        Task<List<Tool>> GetBySKUIdAsync(string skuId);
        Task<List<Tool>> GetByLocationAsync(string locationId);
        Task<List<Tool>> GetByLocationAndSKUAsync(string locationId, string skuId);
        Task<Tool> AddAsync(Tool item);
        Task UpdateBySKUIdAsync(string skuId, Tool item);
        Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate);
        Task<bool> DeleteBySKUIdAsync(string skuId);
        Task<bool> DeleteByPartitionKeyAsync(string partitionKey);
    }
}
