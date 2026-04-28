using Xunit;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;
using Switchyard.LogisticsAPI.Data.Repositories;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Tests.Repositories
{
    public class WarehouseRepositoryTests : IDisposable
    {
        private readonly SqliteConnection _connection;
        private readonly LogisticsReadContext _readContext;
        private readonly WarehouseRepository _repository;

        public WarehouseRepositoryTests()
        {
            _connection = new SqliteConnection("DataSource=:memory:");
            _connection.Open();
            var readOptions = new DbContextOptionsBuilder<LogisticsReadContext>()
                .UseSqlite(_connection)
                .Options;
            _readContext = new LogisticsReadContext(readOptions);
            _readContext.Database.EnsureCreated();
            _repository = new WarehouseRepository(_readContext);
        }

        public void Dispose()
        {
            _readContext.Dispose();
            _connection.Dispose();
        }

        [Fact]
        public async Task GetAllAsync_ReturnsAllWarehouses()
        {
            _readContext.Warehouses.AddRange(
                new Warehouse { WarehouseId = "WH001", City = "Chicago", State = "IL" },
                new Warehouse { WarehouseId = "WH002", City = "Indianapolis", State = "IN" }
            );
            await _readContext.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }
    }
}
