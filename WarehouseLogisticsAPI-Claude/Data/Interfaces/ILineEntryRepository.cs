using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Interfaces
{
    public interface ILineEntryRepository
    {
        Task<IEnumerable<LineEntry>> GetAllAsync();
        Task<List<LineEntry>> GetLineEntriesByTransactionIdAsync(string transactionId);
        Task<LineEntry> AddAsync(LineEntry lineEntry);
        Task UpdateLineEntryAsync(LineEntry lineEntry);
        Task<bool> DeleteByTransactionIdAsync(string transactionId);
        Task<bool> DeleteByLocationAsync(string transactionId, string locationId);
    }
}
