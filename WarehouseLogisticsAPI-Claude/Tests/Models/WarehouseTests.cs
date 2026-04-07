using Xunit;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Models.Interfaces;

namespace WarehouseLogistics_Claude.Tests.Models
{
    public class WarehouseTests
    {
        [Fact]
        public void Warehouse_ImplementsIWarehouse()
        {
            Assert.IsAssignableFrom<IWarehouse>(new Warehouse { WarehouseId = "WH001", City = "Springfield", State = "IL" });
        }

        [Fact]
        public void DefaultPartitionKey_IsEmpty()
        {
            var warehouse = new Warehouse { WarehouseId = "WH001", City = "Springfield", State = "IL" };
            Assert.Equal(string.Empty, warehouse.PartitionKey);
        }

        [Fact]
        public void Properties_SetAndGetCorrectly()
        {
            var warehouse = new Warehouse
            {
                PartitionKey = "WH001-guid",
                WarehouseId = "WH001",
                City = "Springfield",
                State = "IL"
            };

            Assert.Equal("WH001-guid", warehouse.PartitionKey);
            Assert.Equal("WH001", warehouse.WarehouseId);
            Assert.Equal("Springfield", warehouse.City);
            Assert.Equal("IL", warehouse.State);
        }
    }
}
