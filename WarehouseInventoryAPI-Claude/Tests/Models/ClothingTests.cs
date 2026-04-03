using Xunit;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Tests.Models;

public class ClothingTests
{
    [Fact]
    public void Clothing_ImplementsIItem()
    {
        var clothing = new Clothing();

        Assert.IsAssignableFrom<IItem>(clothing);
    }

    [Fact]
    public void Clothing_DefaultProperties_AreEmptyStrings()
    {
        var clothing = new Clothing();

        Assert.Equal(string.Empty, clothing.PartitionKey);
        Assert.Equal(string.Empty, clothing.RowKey);
        Assert.Equal(string.Empty, clothing.SKUMarker);
        Assert.Equal(string.Empty, clothing.Category);
        Assert.Equal(string.Empty, clothing.Type);
        Assert.Equal(string.Empty, clothing.Color);
        Assert.Equal(string.Empty, clothing.Description);
    }

    [Fact]
    public void Clothing_UnloadedDate_DefaultsToUtcNow()
    {
        var before = DateTime.UtcNow;
        var clothing = new Clothing();
        var after = DateTime.UtcNow;

        Assert.InRange(clothing.UnloadedDate, before, after);
    }

    [Fact]
    public void Clothing_Properties_CanBeSet()
    {
        var date = new DateTime(2026, 1, 15, 0, 0, 0, DateTimeKind.Utc);
        var clothing = new Clothing
        {
            PartitionKey = "pk-001",
            RowKey = "rk-001",
            SKUMarker = "SKU001",
            UnloadedDate = date,
            Category = "Workwear",
            Type = "Shirt",
            Color = "Blue",
            Description = "Heavy-duty work shirt"
        };

        Assert.Equal("pk-001", clothing.PartitionKey);
        Assert.Equal("rk-001", clothing.RowKey);
        Assert.Equal("SKU001", clothing.SKUMarker);
        Assert.Equal(date, clothing.UnloadedDate);
        Assert.Equal("Workwear", clothing.Category);
        Assert.Equal("Shirt", clothing.Type);
        Assert.Equal("Blue", clothing.Color);
        Assert.Equal("Heavy-duty work shirt", clothing.Description);
    }
}
