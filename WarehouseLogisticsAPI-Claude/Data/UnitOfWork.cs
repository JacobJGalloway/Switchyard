using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Data.Repositories;

namespace WarehouseLogistics_Claude.Data
{
    public class UnitOfWork(LogisticsContext context) : IUnitOfWork
    {
        private readonly LogisticsContext _context = context;

        private IBillOfLadingRepository? _billsOfLading;
        private ILineEntryRepository? _lineEntries;
        private IWarehouseRepository? _warehouses;
        private IStoreRepository? _stores;

        public IBillOfLadingRepository BillsOfLading => _billsOfLading ??= new BillOfLadingRepository(_context);
        public ILineEntryRepository LineEntries => _lineEntries ??= new LineEntryRepository(_context);
        public IWarehouseRepository Warehouses => _warehouses ??= new WarehouseRepository(_context);
        public IStoreRepository Stores => _stores ??= new StoreRepository(_context);

        public async Task<int> SaveChangesAsync() => await _context.SaveChangesAsync();

        public void Dispose() => _context.Dispose();
    }
}
