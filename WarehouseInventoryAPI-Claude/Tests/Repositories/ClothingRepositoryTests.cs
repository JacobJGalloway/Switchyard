using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Xunit;
using WarehouseInventory_Claude.Data;
using WarehouseInventory_Claude.Data.Repositories;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Tests.Repositories;

public class ClothingRepositoryTests : IDisposable
{
    private readonly SqliteConnection _connection;
    private readonly InventoryContext _context;
    private readonly ClothingRepository _repository;

    public ClothingRepositoryTests()
    {
        _connection = new SqliteConnection("DataSource=:memory:");
        _connection.Open();

        var options = new DbContextOptionsBuilder<InventoryContext>()
            .UseSqlite(_connection)
            .Options;

        _context = new InventoryContext(options);
        _context.Database.EnsureCreated();
        _repository = new ClothingRepository(_context);
    }

    public void Dispose()
    {
        _context.Dispose();
        _connection.Dispose();
    }

    [Fact]
    public async Task GetAllAsync_ReturnsAllItems()
    {
        _context.Clothing.AddRange(
            new Clothing { PartitionKey = "1", SKUMarker = "SKU001", Type = "Shirt" },
            new Clothing { PartitionKey = "2", SKUMarker = "SKU002", Type = "Pants" }
        );
        await _context.SaveChangesAsync();

        var result = await _repository.GetAllAsync();

        Assert.Equal(2, result.Count());
    }

    [Fact]
    public async Task GetBySKUIdAsync_ReturnsItem_WhenFound()
    {
        _context.Clothing.Add(new Clothing { PartitionKey = "1", SKUMarker = "SKU001", Type = "Shirt" });
        await _context.SaveChangesAsync();

        var result = await _repository.GetBySKUIdAsync("SKU001");

        Assert.NotNull(result);
        Assert.Equal("SKU001", result.SKUMarker);
    }

    [Fact]
    public async Task GetBySKUIdAsync_ThrowsKeyNotFoundException_WhenMissing()
    {
        await Assert.ThrowsAsync<KeyNotFoundException>(() =>
            _repository.GetBySKUIdAsync("MISSING"));
    }

    [Fact]
    public async Task AddAsync_PersistsAndReturnsItem()
    {
        var item = new Clothing { PartitionKey = "1", SKUMarker = "SKU001", Type = "Shirt", Color = "Blue" };

        var result = await _repository.AddAsync(item);

        Assert.Equal("SKU001", result.SKUMarker);
        Assert.Equal(1, _context.Clothing.Count());
    }

    [Fact]
    public async Task UpdateBySKUIdAsync_UpdatesExistingItem()
    {
        _context.Clothing.Add(new Clothing { PartitionKey = "1", SKUMarker = "SKU001", Type = "Shirt", Color = "Blue" });
        await _context.SaveChangesAsync();

        var updated = new Clothing { PartitionKey = "1", SKUMarker = "SKU001", Type = "Shirt", Color = "Red" };
        await _repository.UpdateBySKUIdAsync("SKU001", updated);

        var result = await _context.Clothing.FindAsync("1");
        Assert.Equal("Red", result!.Color);
    }

    [Fact]
    public async Task UpdateBySKUIdAsync_ThrowsKeyNotFoundException_WhenSKUNotFound()
    {
        // FindClothingBySKUAsync throws KeyNotFoundException instead of returning null
        await Assert.ThrowsAsync<KeyNotFoundException>(() =>
            _repository.UpdateBySKUIdAsync("MISSING", new Clothing { SKUMarker = "MISSING" }));
    }

    [Fact]
    public async Task DeleteBySKUIdAsync_ReturnsTrueAndRemovesItems_WhenFound()
    {
        _context.Clothing.AddRange(
            new Clothing { PartitionKey = "1", SKUMarker = "SKU001", Type = "Shirt" },
            new Clothing { PartitionKey = "2", SKUMarker = "SKU001", Type = "Shirt" }
        );
        await _context.SaveChangesAsync();

        var result = await _repository.DeleteBySKUIdAsync("SKU001");

        Assert.True(result);
        Assert.Equal(0, _context.Clothing.Count());
    }

    [Fact]
    public async Task DeleteBySKUIdAsync_ReturnsFalse_WhenNotFound()
    {
        var result = await _repository.DeleteBySKUIdAsync("MISSING");

        Assert.False(result);
    }
}
