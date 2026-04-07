using Xunit;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Tests.Models
{
    public class ToolTests
    {
        [Fact]
        public void Tool_ImplementsIItem()
        {
            Assert.IsAssignableFrom<IItem>(new Tool { SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow });
        }

        [Fact]
        public void Tool_DefaultPartitionKey_IsEmpty()
        {
            var item = new Tool { SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            Assert.Equal(string.Empty, item.PartitionKey);
        }

        [Fact]
        public void Tool_DefaultRowKey_IsEmpty()
        {
            var item = new Tool { SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            Assert.Equal(string.Empty, item.RowKey);
        }

        [Fact]
        public void Tool_Properties_SetAndGetCorrectly()
        {
            var date = new DateTime(2025, 1, 10, 8, 0, 0, DateTimeKind.Utc);
            var item = new Tool
            {
                PartitionKey = "WH001-PWTL001-abc123",
                RowKey = "PWTL001",
                SKUMarker = "PWTL001",
                UnloadedDate = date
            };

            Assert.Equal("WH001-PWTL001-abc123", item.PartitionKey);
            Assert.Equal("PWTL001", item.RowKey);
            Assert.Equal("PWTL001", item.SKUMarker);
            Assert.Equal(date, item.UnloadedDate);
        }
    }
}
