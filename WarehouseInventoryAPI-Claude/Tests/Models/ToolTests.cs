using Xunit;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Tests.Models;

public class ToolTests
{
    [Fact]
    public void Tool_ImplementsIItem()
    {
        var tool = new Tool();

        Assert.IsAssignableFrom<IItem>(tool);
    }

    [Fact]
    public void Tool_DefaultStringProperties_AreEmptyStrings()
    {
        var tool = new Tool();

        Assert.Equal(string.Empty, tool.PartitionKey);
        Assert.Equal(string.Empty, tool.RowKey);
        Assert.Equal(string.Empty, tool.SKUMarker);
        Assert.Equal(string.Empty, tool.Type);
        Assert.Equal(string.Empty, tool.Description);
    }

    [Fact]
    public void Tool_DefaultSize_IsZero()
    {
        var tool = new Tool();

        Assert.Equal(0.0, tool.Size);
    }

    [Fact]
    public void Tool_UnloadedDate_DefaultsToUtcNow()
    {
        var before = DateTime.UtcNow;
        var tool = new Tool();
        var after = DateTime.UtcNow;

        Assert.InRange(tool.UnloadedDate, before, after);
    }

    [Fact]
    public void Tool_Properties_CanBeSet()
    {
        var date = new DateTime(2026, 1, 15, 0, 0, 0, DateTimeKind.Utc);
        var tool = new Tool
        {
            PartitionKey = "pk-001",
            RowKey = "rk-001",
            SKUMarker = "SKU001",
            UnloadedDate = date,
            Type = "Wrench",
            Size = 12.5,
            Description = "Adjustable crescent wrench"
        };

        Assert.Equal("pk-001", tool.PartitionKey);
        Assert.Equal("rk-001", tool.RowKey);
        Assert.Equal("SKU001", tool.SKUMarker);
        Assert.Equal(date, tool.UnloadedDate);
        Assert.Equal("Wrench", tool.Type);
        Assert.Equal(12.5, tool.Size);
        Assert.Equal("Adjustable crescent wrench", tool.Description);
    }
}
