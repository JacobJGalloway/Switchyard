using Xunit;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Models.Interfaces;

namespace Switchyard.LogisticsAPI.Tests.Models
{
    public class LineEntryTests
    {
        [Fact]
        public void LineEntry_ImplementsILineEntry()
        {
            Assert.IsAssignableFrom<ILineEntry>(new LineEntry { LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 1 });
        }

        [Fact]
        public void DefaultPartitionKey_IsEmpty()
        {
            var entry = new LineEntry { LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 1 };
            Assert.Equal(string.Empty, entry.PartitionKey);
        }

        [Fact]
        public void DefaultIsProcessed_IsFalse()
        {
            var entry = new LineEntry { LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 1 };
            Assert.False(entry.IsProcessed);
        }

        [Fact]
        public void Properties_SetAndGetCorrectly()
        {
            var processedDate = DateTime.UtcNow;
            var entry = new LineEntry
            {
                PartitionKey = "abc123-guid",
                TransactionId = "abc123",
                LocationId = "WH001",
                SKUMarker = "CLTH001",
                Quantity = -5,
                IsProcessed = true,
                ProcessedDate = processedDate
            };

            Assert.Equal("abc123-guid", entry.PartitionKey);
            Assert.Equal("abc123", entry.TransactionId);
            Assert.Equal("WH001", entry.LocationId);
            Assert.Equal("CLTH001", entry.SKUMarker);
            Assert.Equal(-5, entry.Quantity);
            Assert.True(entry.IsProcessed);
            Assert.Equal(processedDate, entry.ProcessedDate);
        }
    }
}
