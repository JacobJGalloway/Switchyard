using Xunit;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Data;
using WarehouseLogistics_Claude.Data.Repositories;
using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Tests.Repositories
{
    public class LineEntryRepositoryTests : IDisposable
    {
        private readonly SqliteConnection _connection;
        private readonly LogisticsContext _context;
        private readonly LineEntryRepository _repository;

        public LineEntryRepositoryTests()
        {
            _connection = new SqliteConnection("DataSource=:memory:");
            _connection.Open();
            var options = new DbContextOptionsBuilder<LogisticsContext>()
                .UseSqlite(_connection)
                .Options;
            _context = new LogisticsContext(options);
            _context.Database.EnsureCreated();
            _repository = new LineEntryRepository(_context);
        }

        public void Dispose()
        {
            _context.Dispose();
            _connection.Dispose();
        }

        private static LineEntry MakeEntry(string transactionId, string locationId = "WH001") => new()
        {
            PartitionKey = $"{transactionId}-{Guid.NewGuid():N}",
            TransactionId = transactionId,
            LocationId = locationId,
            SKUMarker = "CLTH001",
            Quantity = 5
        };

        [Fact]
        public async Task GetAllAsync_ReturnsAllEntries()
        {
            _context.LineEntries.AddRange(MakeEntry("txn001"), MakeEntry("txn002"));
            await _context.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }

        [Fact]
        public async Task GetLineEntriesByTransactionIdAsync_ReturnsMatchingEntries()
        {
            _context.LineEntries.AddRange(MakeEntry("txn001"), MakeEntry("txn001"), MakeEntry("txn002"));
            await _context.SaveChangesAsync();

            var result = await _repository.GetLineEntriesByTransactionIdAsync("txn001");

            Assert.Equal(2, result.Count);
            Assert.All(result, e => Assert.Equal("txn001", e.TransactionId));
        }

        [Fact]
        public async Task AddAsync_SetsPartitionKey_AndAddsEntry()
        {
            var entry = new LineEntry
            {
                TransactionId = "txn001",
                LocationId = "WH001",
                SKUMarker = "CLTH001",
                Quantity = 3
            };

            await _repository.AddAsync(entry);
            await _context.SaveChangesAsync();

            Assert.Single(_context.LineEntries);
            Assert.StartsWith("txn001-", entry.PartitionKey);
        }

        [Fact]
        public async Task UpdateLineEntryAsync_UpdatesEntry()
        {
            var entry = MakeEntry("txn001");
            _context.LineEntries.Add(entry);
            await _context.SaveChangesAsync();

            entry.IsProcessed = true;
            entry.ProcessedDate = DateTime.UtcNow;
            await _repository.UpdateLineEntryAsync(entry);
            await _context.SaveChangesAsync();

            var result = await _context.LineEntries.FindAsync(entry.PartitionKey);
            Assert.True(result!.IsProcessed);
            Assert.NotNull(result.ProcessedDate);
        }

        [Fact]
        public async Task DeleteByTransactionIdAsync_ReturnsTrue_AndRemovesEntries_WhenFound()
        {
            _context.LineEntries.AddRange(MakeEntry("txn001"), MakeEntry("txn001"), MakeEntry("txn002"));
            await _context.SaveChangesAsync();

            var result = await _repository.DeleteByTransactionIdAsync("txn001");
            await _context.SaveChangesAsync();

            Assert.True(result);
            Assert.Single(_context.LineEntries);
        }

        [Fact]
        public async Task DeleteByTransactionIdAsync_ReturnsFalse_WhenNotFound()
        {
            var result = await _repository.DeleteByTransactionIdAsync("txn999");

            Assert.False(result);
        }
    }
}
