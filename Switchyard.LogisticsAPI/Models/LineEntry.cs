using System.ComponentModel.DataAnnotations;
using Switchyard.LogisticsAPI.Models.Interfaces;

namespace Switchyard.LogisticsAPI.Models
{
    public class LineEntry : ILineEntry
    {
        // Partition key is "{TransactionID}-{GUID}" to allow for 
        // efficient querying of line entries by transaction while 
        // ensuring uniqueness for storage in Cosmos DB or similar NoSQL databases.
        public string PartitionKey { get; set; } = string.Empty;
        public string TransactionId { get; set; } = string.Empty;
        public required string LocationId { get; set; }
        public required string SKUMarker { get; set; }
        // Positive quantity = loaded from this location (OUT); negative quantity = delivered to this location (IN)
        public required int Quantity { get; set; }
        public bool IsProcessed { get; set; } = false;
        public DateTime? ProcessedDate { get; set; }
    }
}
