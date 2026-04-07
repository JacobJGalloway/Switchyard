using Xunit;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Models.Interfaces;

namespace WarehouseLogistics_Claude.Tests.Models
{
    public class BillOfLadingTests
    {
        [Fact]
        public void BillOfLading_ImplementsIBillOfLading()
        {
            Assert.IsAssignableFrom<IBillOfLading>(new BillOfLading { LineEntries = [] });
        }

        [Fact]
        public void DefaultPartitionKey_IsEmpty()
        {
            var bol = new BillOfLading { LineEntries = [] };
            Assert.Equal(string.Empty, bol.PartitionKey);
        }

        [Fact]
        public void DefaultStatus_IsPending()
        {
            var bol = new BillOfLading { LineEntries = [] };
            Assert.Equal("Pending", bol.Status);
        }

        [Fact]
        public void Properties_SetAndGetCorrectly()
        {
            var entries = new List<LineEntry>
            {
                new() { LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 10 }
            };
            var bol = new BillOfLading
            {
                PartitionKey = "abc123-guid",
                TransactionId = "abc123",
                Status = "Submitted",
                CustomerFirstName = "Jane",
                CustomerLastName = "Doe",
                City = "Springfield",
                State = "IL",
                LineEntries = entries
            };

            Assert.Equal("abc123-guid", bol.PartitionKey);
            Assert.Equal("abc123", bol.TransactionId);
            Assert.Equal("Submitted", bol.Status);
            Assert.Equal("Jane", bol.CustomerFirstName);
            Assert.Equal("Doe", bol.CustomerLastName);
            Assert.Equal("Springfield", bol.City);
            Assert.Equal("IL", bol.State);
            Assert.Single(bol.LineEntries);
        }
    }
}
