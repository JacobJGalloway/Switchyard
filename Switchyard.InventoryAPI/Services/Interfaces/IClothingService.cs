using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Services.Interfaces
{
    public interface IClothingService
    {
        Task<IEnumerable<Clothing>> GetAllAsync();
        Task<List<Clothing>> GetBySKUIdAsync(string skuId);
        Task<List<Clothing>> GetByLocationAsync(string locationId);
        Task<List<Clothing>> GetByLocationAndSKUAsync(string locationId, string skuId);
        Task<Clothing> AddAsync(Clothing item);
        Task UpdateBySKUIdAsync(string skuId, Clothing item);
        Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate);
        Task<bool> DeleteByPartitionKeyAsync(string partitionKey);
    }
}
