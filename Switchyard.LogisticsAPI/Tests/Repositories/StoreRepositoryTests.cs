using Xunit;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;
using Switchyard.LogisticsAPI.Data.Repositories;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Tests.Repositories
{
    public class StoreRepositoryTests : IDisposable
    {
        private readonly SqliteConnection _connection;
        private readonly LogisticsReadContext _readContext;
        private readonly StoreRepository _repository;

        public StoreRepositoryTests()
        {
            _connection = new SqliteConnection("DataSource=:memory:");
            _connection.Open();
            var readOptions = new DbContextOptionsBuilder<LogisticsReadContext>()
                .UseSqlite(_connection)
                .Options;
            _readContext = new LogisticsReadContext(readOptions);
            _readContext.Database.EnsureCreated();
            _repository = new StoreRepository(_readContext);
        }

        public void Dispose()
        {
            _readContext.Dispose();
            _connection.Dispose();
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
