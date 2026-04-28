using Xunit;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;
using Switchyard.LogisticsAPI.Data.Repositories;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Tests.Repositories
{
    public class BillOfLadingRepositoryTests : IDisposable
    {
        private readonly SqliteConnection _connection;
        private readonly LogisticsContext _context;
        private readonly LogisticsReadContext _readContext;
        private readonly BillOfLadingRepository _repository;

        public BillOfLadingRepositoryTests()
        {
            _connection = new SqliteConnection("DataSource=:memory:");
            _connection.Open();
            var writeOptions = new DbContextOptionsBuilder<LogisticsContext>()
                .UseSqlite(_connection)
                .Options;
            var readOptions = new DbContextOptionsBuilder<LogisticsReadContext>()
                .UseSqlite(_connection)
                .Options;
            _context = new LogisticsContext(writeOptions);
            _context.Database.EnsureCreated();
            _readContext = new LogisticsReadContext(readOptions);
            _repository = new BillOfLadingRepository(_context, _readContext);
        }

        public void Dispose()
        {
            _context.Dispose();
            _readContext.Dispose();
            _connection.Dispose();
        }

        private static BillOfLading MakeBOL(string transactionId) => new()
        {
            PartitionKey = $"{transactionId}-testguid",
            TransactionId = transactionId,
            Status = "Submitted",
            CustomerFirstName = "Jane",
            CustomerLastName = "Doe",
            City = "Springfield",
            State = "IL",
            LineEntries = []
        };

        [Fact]
        public async Task GetAllAsync_ReturnsAllBOLs()
        {
            _context.BillsOfLading.AddRange(MakeBOL("txn001"), MakeBOL("txn002"));
            await _context.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }

        [Fact]
        public async Task GetByTransactionIdAsync_ReturnsBOL_WhenFound()
        {
            _context.BillsOfLading.Add(MakeBOL("txn001"));
            await _context.SaveChangesAsync();

            var result = await _repository.GetByTransactionIdAsync("txn001");

            Assert.NotNull(result);
            Assert.Equal("txn001", result.TransactionId);
        }

        [Fact]
        public async Task GetByTransactionIdAsync_ReturnsNull_WhenNotFound()
        {
            var result = await _repository.GetByTransactionIdAsync("txn999");

            Assert.Null(result);
        }

        [Fact]
        public async Task AddAsync_AddsBOL()
        {
            var bol = MakeBOL("txn001");

            await _repository.AddAsync(bol);
            await _context.SaveChangesAsync();

            Assert.Single(_context.BillsOfLading);
        }

        [Fact]
        public async Task UpdateAsync_UpdatesScalarFields_WhenFound()
        {
            _context.BillsOfLading.Add(MakeBOL("txn001"));
            await _context.SaveChangesAsync();

            var updated = MakeBOL("txn001");
            updated.Status = "In Transit";
            updated.CustomerFirstName = "John";

            await _repository.UpdateAsync(updated);
            await _context.SaveChangesAsync();

            var result = await _repository.GetByTransactionIdAsync("txn001");
            Assert.Equal("In Transit", result!.Status);
            Assert.Equal("John", result.CustomerFirstName);
        }

        [Fact]
        public async Task UpdateAsync_DoesNothing_WhenNotFound()
        {
            var updated = MakeBOL("txn999");
            updated.Status = "In Transit";

            await _repository.UpdateAsync(updated);
            await _context.SaveChangesAsync();

            Assert.Empty(_context.BillsOfLading);
        }

        [Fact]
        public async Task DeleteByTransactionIdAsync_ReturnsTrue_AndRemoves_WhenFound()
        {
            _context.BillsOfLading.Add(MakeBOL("txn001"));
            await _context.SaveChangesAsync();

            var result = await _repository.DeleteByTransactionIdAsync("txn001");
            await _context.SaveChangesAsync();

            Assert.True(result);
            Assert.Empty(_context.BillsOfLading);
        }

        [Fact]
        public async Task DeleteByTransactionIdAsync_ReturnsFalse_WhenNotFound()
        {
            var result = await _repository.DeleteByTransactionIdAsync("txn999");

            Assert.False(result);
        }
    }
}
