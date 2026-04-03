using Xunit;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Tests.Models;

public class PPETests
{
    [Fact]
    public void PPE_ImplementsIItem()
    {
        var ppe = new PPE();

        Assert.IsAssignableFrom<IItem>(ppe);
    }

    [Fact]
    public void PPE_DefaultProperties_AreEmptyStrings()
    {
        var ppe = new PPE();

        Assert.Equal(string.Empty, ppe.PartitionKey);
        Assert.Equal(string.Empty, ppe.RowKey);
        Assert.Equal(string.Empty, ppe.SKUMarker);
        Assert.Equal(string.Empty, ppe.Category);
        Assert.Equal(string.Empty, ppe.Type);
        Assert.Equal(string.Empty, ppe.Size);
    }

    [Fact]
    public void PPE_UnloadedDate_DefaultsToUtcNow()
    {
        var before = DateTime.UtcNow;
        var ppe = new PPE();
        var after = DateTime.UtcNow;

        Assert.InRange(ppe.UnloadedDate, before, after);
    }

    [Fact]
    public void PPE_Properties_CanBeSet()
    {
        var date = new DateTime(2026, 1, 15, 0, 0, 0, DateTimeKind.Utc);
        var ppe = new PPE
        {
            PartitionKey = "pk-001",
            RowKey = "rk-001",
            SKUMarker = "SKU001",
            UnloadedDate = date,
            Category = "Head Protection",
            Type = "Helmet",
            Size = "L"
        };

        Assert.Equal("pk-001", ppe.PartitionKey);
        Assert.Equal("rk-001", ppe.RowKey);
        Assert.Equal("SKU001", ppe.SKUMarker);
        Assert.Equal(date, ppe.UnloadedDate);
        Assert.Equal("Head Protection", ppe.Category);
        Assert.Equal("Helmet", ppe.Type);
        Assert.Equal("L", ppe.Size);
    }
}
