using Xunit;
using Moq;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data;
using Switchyard.InventoryAPI.Data.Repositories;
using Switchyard.Domain;

namespace Switchyard.InventoryAPI.Tests.Repositories
{
    public class ClothingRepositoryTests : IDisposable
    {
        private readonly InventoryContext _context;
        private readonly InventoryReadContext _readContext;
        private readonly ClothingRepository _repository;

        public ClothingRepositoryTests()
        {
            var dbName = Guid.NewGuid().ToString();
            var writeOptions = new DbContextOptionsBuilder<InventoryContext>()
                .UseInMemoryDatabase(dbName)
                .Options;
            var readOptions = new DbContextOptionsBuilder<InventoryReadContext>()
                .UseInMemoryDatabase(dbName)
                .Options;
            _context = new InventoryContext(writeOptions);
            _context.Database.EnsureCreated();
            _readContext = new InventoryReadContext(readOptions);
            _repository = new ClothingRepository(_context, _readContext);
        }

        public void Dispose()
        {
            _context.Dispose();
            _readContext.Dispose();
        }

        [Fact]
        public async Task GetAllAsync_ReturnsAllItems()
        {
            _context.Clothing.AddRange(
                new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow },
                new Clothing { PartitionKey = "WH001-CLTH002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "CLTH002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.GetAllAsync();

            Assert.Equal(2, result.Count());
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsMatchingItems_WhenFound()
        {
            _context.Clothing.AddRange(
                new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow },
                new Clothing { PartitionKey = "WH001-CLTH001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow },
                new Clothing { PartitionKey = "WH001-CLTH002-c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8", SKUMarker = "CLTH002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.GetBySKUIdAsync("CLTH001");

            Assert.Equal(2, result.Count);
            Assert.All(result, c => Assert.Equal("CLTH001", c.SKUMarker));
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsEmptyList_WhenNotFound()
        {
            var result = await _repository.GetBySKUIdAsync("CLTH999");

            Assert.Empty(result);
        }

        [Fact]
        public async Task AddAsync_StagesItem_PersistedAfterSave()
        {
            var item = new Clothing { SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };

            await _repository.AddAsync(item);
            await _context.SaveChangesAsync();

            Assert.Equal(1, await _context.Clothing.CountAsync());
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_UpdatesByPartitionKey_WhenMatch()
        {
            var original = new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            var other = new Clothing { PartitionKey = "WH001-CLTH001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            _context.Clothing.AddRange(original, other);
            await _context.SaveChangesAsync();

            var updated = new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow.AddDays(1) };
            await _repository.UpdateBySKUIdAsync("CLTH001", updated);
            await _context.SaveChangesAsync();

            var result = await _context.Clothing.FindAsync("WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(updated.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_FallsBackToFirst_WhenNoPartitionKeyMatch()
        {
            var item = new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            _context.Clothing.Add(item);
            await _context.SaveChangesAsync();

            var updated = new Clothing { PartitionKey = "WH001-CLTH001-ffffffffffffffffffffffffffffffff", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow.AddDays(1) };
            await _repository.UpdateBySKUIdAsync("CLTH001", updated);
            await _context.SaveChangesAsync();

            var result = await _context.Clothing.FindAsync("WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(updated.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_DoesNothing_WhenSkuNotFound()
        {
            var item = new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            _context.Clothing.Add(item);
            await _context.SaveChangesAsync();

            await _repository.UpdateBySKUIdAsync("CLTH999", new Clothing { PartitionKey = "WH001-CLTH999-ffffffffffffffffffffffffffffffff", SKUMarker = "CLTH999", UnloadedDate = DateTime.UtcNow.AddDays(1) });
            await _context.SaveChangesAsync();

            var result = await _context.Clothing.FindAsync("WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6");
            Assert.Equal(item.UnloadedDate, result!.UnloadedDate);
        }

        [Fact]
        public async Task DeleteBySKUIdAsync_ReturnsFalse_WhenNotFound()
        {
            var result = await _repository.DeleteBySKUIdAsync("CLTH999");

            Assert.False(result);
        }

        [Fact]
        public async Task DeleteBySKUIdAsync_ReturnsTrue_AndRemovesAllMatchingItems()
        {
            _context.Clothing.AddRange(
                new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow },
                new Clothing { PartitionKey = "WH001-CLTH001-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow },
                new Clothing { PartitionKey = "WH001-CLTH002-c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8", SKUMarker = "CLTH002", UnloadedDate = DateTime.UtcNow }
            );
            await _context.SaveChangesAsync();

            var result = await _repository.DeleteBySKUIdAsync("CLTH001");
            await _context.SaveChangesAsync();

            Assert.True(result);
            Assert.Equal(1, await _context.Clothing.CountAsync());
        }
    }
}
