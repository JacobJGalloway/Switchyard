using Xunit;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Models.Interfaces;

namespace Switchyard.LogisticsAPI.Tests.Models
{
    public class StoreTests
    {
        [Fact]
        public void Store_ImplementsIStore()
        {
            Assert.IsAssignableFrom<IStore>(new Store { StoreId = "ST0001", BaseWarehouseId = "WH001", City = "Springfield", State = "IL" });
        }

        [Fact]
        public void DefaultPartitionKey_IsEmpty()
        {
            var store = new Store { StoreId = "ST0001", BaseWarehouseId = "WH001", City = "Springfield", State = "IL" };
            Assert.Equal(string.Empty, store.PartitionKey);
        }

        [Fact]
        public void Properties_SetAndGetCorrectly()
        {
            var store = new Store
            {
                PartitionKey = "ST0001-WH001-guid",
                StoreId = "ST0001",
                BaseWarehouseId = "WH001",
                City = "Springfield",
                State = "IL"
            };

            Assert.Equal("ST0001-WH001-guid", store.PartitionKey);
            Assert.Equal("ST0001", store.StoreId);
            Assert.Equal("WH001", store.BaseWarehouseId);
            Assert.Equal("Springfield", store.City);
            Assert.Equal("IL", store.State);
        }
    }
}
