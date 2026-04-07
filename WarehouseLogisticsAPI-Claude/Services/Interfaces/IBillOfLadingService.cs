using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Services.Interfaces
{
    public interface IBillOfLadingService
    {
        Task<string> CreateAsync(BillOfLading billOfLading);
        Task ProcessLocationStop(string transactionId, string locationId);
    }
}
