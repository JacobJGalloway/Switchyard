using Switchyard.Domain;

namespace Switchyard.InventoryAPI.Data.Interfaces
{
    public interface IUnitOfWork : IDisposable
    {

        IClothingRepository Clothing { get; }
        IPPERepository PPE { get; }
        IToolRepository Tools { get; }
        Task<List<Clothing>> GetClothingBySKUIdAsync(string skuId);
        Task<List<PPE>> GetPPEBySKUIdAsync(string skuId);
        Task<List<Tool>> GetToolBySKUIdAsync(string skuId);
        Task ReceiveDeliveryAsync(string locationId, List<DeliveryLineItem> items);
        Task<int> SaveChangesAsync();
    }
}
