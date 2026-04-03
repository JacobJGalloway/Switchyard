using Microsoft.Data.Sqlite;
using Microsoft.EntityFrameworkCore;
using Xunit;
using WarehouseInventory_Claude.Data;
using WarehouseInventory_Claude.Data.Repositories;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Tests.Repositories;

public class ToolRepositoryTests : IDisposable
{
    private readonly SqliteConnection _connection;
    private readonly InventoryContext _context;
    private readonly ToolRepository _repository;

    public ToolRepositoryTests()
    {
        _connection = new SqliteConnection("DataSource=:memory:");
        _connection.Open();

        var options = new DbContextOptionsBuilder<InventoryContext>()
            .UseSqlite(_connection)
            .Options;

        _context = new InventoryContext(options);
        _context.Database.EnsureCreated();
        _repository = new ToolRepository(_context);
    }

    public void Dispose()
    {
        _context.Dispose();
        _connection.Dispose();
    }

    [Fact]
    public async Task GetAllAsync_ReturnsAllItems()
    {
        _context.Tools.AddRange(
            new Tool { PartitionKey = "1", SKUMarker = "SKU001", Type = "Wrench", Size = 12.5 },
            new Tool { PartitionKey = "2", SKUMarker = "SKU002", Type = "Hammer", Size = 16.0 }
        );
        await _context.SaveChangesAsync();

        var result = await _repository.GetAllAsync();

        Assert.Equal(2, result.Count());
    }

    [Fact]
    public async Task GetBySKUIdAsync_ReturnsItem_WhenFound()
    {
        _context.Tools.Add(new Tool { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Wrench" });
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
        var item = new Tool { PartitionKey = "1", SKUMarker = "SKU001", Type = "Wrench", Size = 12.5 };

        var result = await _repository.AddAsync(item);

        Assert.Equal("SKU001", result.SKUMarker);
        Assert.Equal(1, _context.Tools.Count());
    }

    [Fact]
    public async Task UpdateBySKUIdAsync_UpdatesExistingItem()
    {
        _context.Tools.Add(new Tool { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Wrench", Size = 10.0 });
        await _context.SaveChangesAsync();

        var updated = new Tool { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Wrench", Size = 14.0 };
        await _repository.UpdateBySKUIdAsync("SKU001", updated);

        var result = await _context.Tools.FindAsync("SKU001");
        Assert.Equal(14.0, result!.Size);
    }

    [Fact]
    public async Task UpdateBySKUIdAsync_DoesNothing_WhenSKUNotFound()
    {
        var exception = await Record.ExceptionAsync(() =>
            _repository.UpdateBySKUIdAsync("MISSING", new Tool { SKUMarker = "MISSING" }));

        Assert.Null(exception);
    }

    [Fact]
    public async Task DeleteBySKUIdAsync_ReturnsTrueAndRemovesItem_WhenFound()
    {
        _context.Tools.Add(new Tool { PartitionKey = "SKU001", SKUMarker = "SKU001", Type = "Wrench" });
        await _context.SaveChangesAsync();

        var result = await _repository.DeleteBySKUIdAsync("SKU001");

        Assert.True(result);
        Assert.Equal(0, _context.Tools.Count());
    }

    [Fact]
    public async Task DeleteBySKUIdAsync_ReturnsFalse_WhenNotFound()
    {
        var result = await _repository.DeleteBySKUIdAsync("MISSING");

        Assert.False(result);
    }
}
