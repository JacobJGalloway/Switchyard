using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Services.Interfaces
{
    public interface IBillOfLadingService
    {
        Task<string> CreateAsync(BillOfLading billOfLading);
        Task ProcessLocationStop(string transactionId, string locationId);
        Task ReplaceLocationStopAsync(string transactionId, string oldLocationId, string newLocationId);
        Task<IEnumerable<BillOfLading>> GetAllAsync();
        Task<BillOfLading?> GetByTransactionIdAsync(string transactionId);
        Task<List<LineEntry>> GetLineEntriesByTransactionIdAsync(string transactionId);
    }
}
