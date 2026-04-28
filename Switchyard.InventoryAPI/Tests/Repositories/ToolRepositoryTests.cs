using Xunit;
using Moq;
using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data;
using Switchyard.InventoryAPI.Data.Repositories;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Tests.Repositories
{
    public class ToolRepositoryTests : IDisposable
    {
        private readonly SqliteConnection _connection;
        private readonly InventoryContext _context;
        private readonly InventoryReadContext _readContext;
        private readonly ToolRepository _repository;

        public ToolRepositoryTests()
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
            _repository = new ToolRepository(_context, _readContext);
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
            _context.Tool.AddRange(
                new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new Tool { PartitionKey = "WH001-PWTL002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsMatchingItems_WhenFound()
        {
            _context.Tool.AddRange(
                new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new Tool { PartitionKey = "WH001-PWTL001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new Tool { PartitionKey = "WH001-PWTL002-c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8", SKUMarker = "PWTL002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.GetBySKUIdAsync("PWTL001");

            Assert.Equal(2, result.Count);
            Assert.All(result, t => Assert.Equal("PWTL001", t.SKUMarker));
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsEmptyList_WhenNotFound()
        {
            var result = await _repository.GetBySKUIdAsync("PWTL999");

            Assert.Empty(result);
        }

        [Fact]
        public async Task AddAsync_StagesItem_PersistedAfterSave()
        {
            var item = new Tool { SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };

            await _repository.AddAsync(item);
            await _context.SaveChangesAsync();

            Assert.Equal(1, await _context.Tool.CountAsync());
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_UpdatesByPartitionKey_WhenMatch()
        {
            var original = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            var other = new Tool { PartitionKey = "WH001-PWTL001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _context.Tool.AddRange(original, other);
            await _context.SaveChangesAsync();

            var updated = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow.AddDays(1) };
            await _repository.UpdateBySKUIdAsync("PWTL001", updated);
            await _context.SaveChangesAsync();

            var result = await _context.Tool.FindAsync("WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(updated.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_FallsBackToFirst_WhenNoPartitionKeyMatch()
        {
            var item = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _context.Tool.Add(item);
            await _context.SaveChangesAsync();

            var updated = new Tool { PartitionKey = "WH001-PWTL001-ffffffffffffffffffffffffffffffff", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow.AddDays(1) };
            await _repository.UpdateBySKUIdAsync("PWTL001", updated);
            await _context.SaveChangesAsync();

            var result = await _context.Tool.FindAsync("WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(updated.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_DoesNothing_WhenSkuNotFound()
        {
            var item = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _context.Tool.Add(item);
            await _context.SaveChangesAsync();

            await _repository.UpdateBySKUIdAsync("PWTL999", new Tool { PartitionKey = "WH001-PWTL999-ffffffffffffffffffffffffffffffff", SKUMarker = "PWTL999", UnloadedDate = DateTime.UtcNow.AddDays(1) });
            await _context.SaveChangesAsync();

            var result = await _context.Tool.FindAsync("WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(item.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task DeleteBySKUIdAsync_ReturnsFalse_WhenNotFound()
        {
            var result = await _repository.DeleteBySKUIdAsync("PWTL999");

            Assert.False(result);
        }

        [Fact]
        public async Task DeleteBySKUIdAsync_ReturnsTrue_AndRemovesAllMatchingItems()
        {
            _context.Tool.AddRange(
                new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new Tool { PartitionKey = "WH001-PWTL001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new Tool { PartitionKey = "WH001-PWTL002-c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8", SKUMarker = "PWTL002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.DeleteBySKUIdAsync("PWTL001");
            await _context.SaveChangesAsync();

            Assert.True(result);
            Assert.Equal(1, await _context.Tool.CountAsync());
        }
    }
}
