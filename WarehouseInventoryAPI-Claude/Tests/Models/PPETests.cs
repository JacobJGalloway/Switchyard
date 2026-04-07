using Xunit;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Tests.Models
{
    public class PPETests
    {
        [Fact]
        public void PPE_ImplementsIItem()
        {
            Assert.IsAssignableFrom<IItem>(new PPE { SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow });
        }

        [Fact]
        public void PPE_DefaultPartitionKey_IsEmpty()
        {
            var item = new PPE { SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            Assert.Equal(string.Empty, item.PartitionKey);
        }

        [Fact]
        public void PPE_DefaultRowKey_IsEmpty()
        {
            var item = new PPE { SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            Assert.Equal(string.Empty, item.RowKey);
        }

        [Fact]
        public void PPE_Properties_SetAndGetCorrectly()
        {
            var date = new DateTime(2025, 1, 10, 8, 0, 0, DateTimeKind.Utc);
            var item = new PPE
            {
                PartitionKey = "WH001-SPPE001-abc123",
                RowKey = "SPPE001",
                SKUMarker = "SPPE001",
                UnloadedDate = date
            };

            Assert.Equal("WH001-SPPE001-abc123", item.PartitionKey);
            Assert.Equal("SPPE001", item.RowKey);
            Assert.Equal("SPPE001", item.SKUMarker);
            Assert.Equal(date, item.UnloadedDate);
        }
    }
}
