using Xunit;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Models.Interfaces;

namespace Switchyard.InventoryAPI.Tests.Models
{
    public class ClothingTests
    {
        [Fact]
        public void Clothing_ImplementsIItem()
        {
            Assert.IsAssignableFrom<IItem>(new Clothing { SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow });
        }

        [Fact]
        public void Clothing_DefaultPartitionKey_IsEmpty()
        {
            var item = new Clothing { SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            Assert.Equal(string.Empty, item.PartitionKey);
        }

        [Fact]
        public void Clothing_DefaultRowKey_IsEmpty()
        {
            var item = new Clothing { SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            Assert.Equal(string.Empty, item.RowKey);
        }

        [Fact]
        public void Clothing_Properties_SetAndGetCorrectly()
        {
            var date = new DateTime(2025, 1, 10, 8, 0, 0, DateTimeKind.Utc);
            var item = new Clothing
            {
                PartitionKey = "WH001-CLTH001-abc123",
                RowKey = "CLTH001",
                SKUMarker = "CLTH001",
                UnloadedDate = date
            };

            Assert.Equal("WH001-CLTH001-abc123", item.PartitionKey);
            Assert.Equal("CLTH001", item.RowKey);
            Assert.Equal("CLTH001", item.SKUMarker);
            Assert.Equal(date, item.UnloadedDate);
        }
    }
}
