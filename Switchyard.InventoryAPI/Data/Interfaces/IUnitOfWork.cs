namespace Switchyard.InventoryAPI.Data.Interfaces
{
    public interface IUnitOfWork : IDisposable
    {
        IClothingRepository Clothing { get; }
        IPPERepository PPE { get; }
        IToolRepository Tools { get; }
        Task<List<Models.Clothing>> GetClothingBySKUIdAsync(string skuId);
        Task<List<Models.PPE>> GetPPEBySKUIdAsync(string skuId);
        Task<List<Models.Tool>> GetToolBySKUIdAsync(string skuId);
        Task<int> SaveChangesAsync();
    }
}
