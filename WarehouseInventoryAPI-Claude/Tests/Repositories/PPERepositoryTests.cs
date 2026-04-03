using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Xunit;
using WarehouseInventory_Claude.Data;
using WarehouseInventory_Claude.Data.Repositories;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Tests.Repositories;

public class PPERepositoryTests : IDisposable
{
    private readonly SqliteConnection _connection;
    private readonly InventoryContext _context;
    private readonly PPERepository _repository;

    public PPERepositoryTests()
    {
        _connection = new SqliteConnection("DataSource=:memory:");
        _connection.Open();

        var options = new DbContextOptionsBuilder<InventoryContext>()
            .UseSqlite(_connection)
            .Options;

        _context = new InventoryContext(options);
        _context.Database.EnsureCreated();
        _repository = new PPERepository(_context);
    }

    public void Dispose()
    {
        _context.Dispose();
        _connection.Dispose();
    }

    [Fact]
    public async Task GetAllAsync_ReturnsAllItems()
    {
        _context.PPE.AddRange(
            new PPE { PartitionKey = "1", SKUMarker = "SKU001", Type = "Gloves", Size = "L" },
            new PPE { PartitionKey = "2", SKUMarker = "SKU002", Type = "Helmet", Size = "M" }
        );
        await _context.SaveChangesAsync();

        var result = await _repository.GetAllAsync();

        Assert.Equal(2, result.Count());
    }

    [Fact]
    public async Task GetBySKUIdAsync_ReturnsItem_WhenFound()
    {
        _context.PPE.Add(new PPE { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Gloves" });
        await _context.SaveChangesAsync();

        var result = await _repository.GetBySKUIdAsync("SKU001");

        Assert.NotNull(result);
        Assert.Equal("SKU001", result.SKUMarker);
    }

    [Fact]
    public async Task GetBySKUIdAsync_ReturnsNull_WhenMissing()
    {
        var result = await _repository.GetBySKUIdAsync("MISSING");

        Assert.Null(result);
    }

    [Fact]
    public async Task AddAsync_PersistsAndReturnsItem()
    {
        var item = new PPE { PartitionKey = "1", SKUMarker = "SKU001", Type = "Gloves", Size = "L" };

        var result = await _repository.AddAsync(item);

        Assert.Equal("SKU001", result.SKUMarker);
        Assert.Equal(1, _context.PPE.Count());
    }

    [Fact]
    public async Task UpdateBySKUIdAsync_UpdatesExistingItem()
    {
        _context.PPE.Add(new PPE { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Gloves", Size = "M" });
        await _context.SaveChangesAsync();

        var updated = new PPE { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Gloves", Size = "L" };
        await _repository.UpdateBySKUIdAsync("SKU001", updated);

        var result = await _context.PPE.FindAsync("SKU001");
        Assert.Equal("L", result!.Size);
    }

    [Fact]
    public async Task UpdateBySKUIdAsync_DoesNothing_WhenSKUNotFound()
    {
        var exception = await Record.ExceptionAsync(() =>
            _repository.UpdateBySKUIdAsync("MISSING", new PPE { SKUMarker = "MISSING" }));

        Assert.Null(exception);
    }

    [Fact]
    public async Task DeleteBySKUIdAsync_ReturnsTrueAndRemovesItem_WhenFound()
    {
        _context.PPE.Add(new PPE { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Gloves" });
        await _context.SaveChangesAsync();

        var result = await _repository.DeleteBySKUIdAsync("SKU001");

        Assert.True(result);
        Assert.Equal(0, _context.PPE.Count());
    }

    [Fact]
    public async Task DeleteBySKUIdAsync_ReturnsFalse_WhenNotFound()
    {
        var result = await _repository.DeleteBySKUIdAsync("MISSING");

        Assert.False(result);
    }
}
