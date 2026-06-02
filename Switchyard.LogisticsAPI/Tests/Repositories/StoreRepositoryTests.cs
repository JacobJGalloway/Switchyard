using Xunit;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;
using Switchyard.LogisticsAPI.Data.Repositories;
using Switchyard.Domain;

namespace Switchyard.LogisticsAPI.Tests.Repositories
{
    public class StoreRepositoryTests : IDisposable
    {
        private readonly LogisticsReadContext _readContext;
        private readonly StoreRepository _repository;

        public StoreRepositoryTests()
        {
            var dbName = Guid.NewGuid().ToString();
            var readOptions = new DbContextOptionsBuilder<LogisticsReadContext>()
                .UseInMemoryDatabase(dbName)
                .Options;
            _readContext = new LogisticsReadContext(readOptions);
            _readContext.Database.EnsureCreated();
            _repository = new StoreRepository(_readContext);
        }

        public void Dispose()
        {
            _readContext.Dispose();
        }

        [Fact]
        public async Task GetAllAsync_ReturnsAllStores()
        {
            _readContext.Stores.AddRange(
                new Store { PartitionKey = "ST0001-pk", StoreId = "ST0001", BaseWarehouseId = "WH001", City = "Chicago", State = "IL" },
                new Store { PartitionKey = "ST0002-pk", StoreId = "ST0002", BaseWarehouseId = "WH001", City = "Naperville", State = "IL" }
            );
            await _readContext.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }
    }
}
