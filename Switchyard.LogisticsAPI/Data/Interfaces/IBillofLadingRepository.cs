using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Interfaces
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
