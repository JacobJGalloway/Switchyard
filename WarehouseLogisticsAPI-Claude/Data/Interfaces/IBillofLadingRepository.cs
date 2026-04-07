using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Interfaces
{
    public interface IBillOfLadingRepository
    {
        Task<IEnumerable<BillOfLading>> GetAllAsync();
        Task<BillOfLading?> GetByTransactionIdAsync(string transactionId);
        Task<BillOfLading> AddAsync(BillOfLading billOfLading);
        Task UpdateAsync(BillOfLading billOfLading);
        Task<bool> DeleteByTransactionIdAsync(string transactionId);
    }
}
