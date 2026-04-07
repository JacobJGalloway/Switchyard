using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Data.Interfaces
{
    public interface IPPERepository
    {
        Task<IEnumerable<PPE>> GetAllAsync();
        Task<List<PPE>> GetBySKUIdAsync(string skuId);
        Task<List<PPE>> GetByLocationAsync(string locationId);
        Task<List<PPE>> GetByLocationAndSKUAsync(string locationId, string skuId);
        Task<PPE> AddAsync(PPE item);
        Task UpdateBySKUIdAsync(string skuId, PPE item);
        Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate);
        Task<bool> DeleteBySKUIdAsync(string skuId);
        Task<bool> DeleteByPartitionKeyAsync(string partitionKey);
    }
}
