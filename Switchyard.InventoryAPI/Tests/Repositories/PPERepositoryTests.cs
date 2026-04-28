using Xunit;
using Moq;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data;
using Switchyard.InventoryAPI.Data.Repositories;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Tests.Repositories
{
    public class PPERepositoryTests : IDisposable
    {
        private readonly SqliteConnection _connection;
        private readonly InventoryContext _context;
        private readonly InventoryReadContext _readContext;
        private readonly PPERepository _repository;

        public PPERepositoryTests()
        {
            _connection = new SqliteConnection("DataSource=:memory:");
            _connection.Open();
            var writeOptions = new DbContextOptionsBuilder<InventoryContext>()
                .UseSqlite(_connection)
                .Options;
            var readOptions = new DbContextOptionsBuilder<InventoryReadContext>()
                .UseSqlite(_connection)
                .Options;
            _context = new InventoryContext(writeOptions);
            _context.Database.EnsureCreated();
            _readContext = new InventoryReadContext(readOptions);
            _repository = new PPERepository(_context, _readContext);
        }

        public void Dispose()
        {
            _context.Dispose();
            _readContext.Dispose();
            _connection.Dispose();
        }

        [Fact]
        public async Task GetAllAsync_ReturnsAllItems()
        {
            _context.PPE.AddRange(
                new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow },
                new PPE { PartitionKey = "WH001-SPPE002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "SPPE002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsMatchingItems_WhenFound()
        {
            _context.PPE.AddRange(
                new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow },
                new PPE { PartitionKey = "WH001-SPPE001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow },
                new PPE { PartitionKey = "WH001-SPPE002-c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8", SKUMarker = "SPPE002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.GetBySKUIdAsync("SPPE001");

            Assert.Equal(2, result.Count);
            Assert.All(result, p => Assert.Equal("SPPE001", p.SKUMarker));
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsEmptyList_WhenNotFound()
        {
            var result = await _repository.GetBySKUIdAsync("SPPE999");

            Assert.Empty(result);
        }

        [Fact]
        public async Task AddAsync_StagesItem_PersistedAfterSave()
        {
            var item = new PPE { SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };

            await _repository.AddAsync(item);
            await _context.SaveChangesAsync();

            Assert.Equal(1, await _context.PPE.CountAsync());
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_UpdatesByPartitionKey_WhenMatch()
        {
            var original = new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            var other = new PPE { PartitionKey = "WH001-SPPE001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            _context.PPE.AddRange(original, other);
            await _context.SaveChangesAsync();

            var updated = new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow.AddDays(1) };
            await _repository.UpdateBySKUIdAsync("SPPE001", updated);
            await _context.SaveChangesAsync();

            var result = await _context.PPE.FindAsync("WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(updated.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_FallsBackToFirst_WhenNoPartitionKeyMatch()
        {
            var item = new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            _context.PPE.Add(item);
            await _context.SaveChangesAsync();

            var updated = new PPE { PartitionKey = "WH001-SPPE001-ffffffffffffffffffffffffffffffff", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow.AddDays(1) };
            await _repository.UpdateBySKUIdAsync("SPPE001", updated);
            await _context.SaveChangesAsync();

            var result = await _context.PPE.FindAsync("WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(updated.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_DoesNothing_WhenSkuNotFound()
        {
            var item = new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            _context.PPE.Add(item);
            await _context.SaveChangesAsync();

            await _repository.UpdateBySKUIdAsync("SPPE999", new PPE { PartitionKey = "WH001-SPPE999-ffffffffffffffffffffffffffffffff", SKUMarker = "SPPE999", UnloadedDate = DateTime.UtcNow.AddDays(1) });
            await _context.SaveChangesAsync();

            var result = await _context.PPE.FindAsync("WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(item.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task DeleteBySKUIdAsync_ReturnsFalse_WhenNotFound()
        {
            var result = await _repository.DeleteBySKUIdAsync("SPPE999");

            Assert.False(result);
        }

        [Fact]
        public async Task DeleteBySKUIdAsync_ReturnsTrue_AndRemovesAllMatchingItems()
        {
            _context.PPE.AddRange(
                new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow },
                new PPE { PartitionKey = "WH001-SPPE001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow },
                new PPE { PartitionKey = "WH001-SPPE002-c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8", SKUMarker = "SPPE002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.DeleteBySKUIdAsync("SPPE001");
            await _context.SaveChangesAsync();

            Assert.True(result);
            Assert.Equal(1, await _context.PPE.CountAsync());
        }
    }
}
