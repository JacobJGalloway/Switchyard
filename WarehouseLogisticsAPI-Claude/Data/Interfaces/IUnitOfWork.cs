namespace WarehouseLogistics_Claude.Data.Interfaces
{
    public interface IUnitOfWork : IDisposable
    {
        IBillOfLadingRepository BillsOfLading { get; }
        ILineEntryRepository LineEntries { get; }
        IWarehouseRepository Warehouses { get; }
        IStoreRepository Stores { get; }
        Task<int> SaveChangesAsync();
    }
}
