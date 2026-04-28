using System.ComponentModel.DataAnnotations;
using Switchyard.LogisticsAPI.Models.Interfaces;

namespace Switchyard.LogisticsAPI.Models
{
    public class BillOfLading : IBillOfLading
    {
        // Partition key is "{TransactionID}-{GUID}" to allow for 
        // efficient querying of line entries by transaction while ensuring 
        // uniqueness for storage in Cosmos DB or similar NoSQL databases.
        public string PartitionKey { get; set; } = string.Empty;
        public string TransactionId { get; set; } = string.Empty;
        public string Status { get; set; } = "Pending"; // Possible values: Pending, Submitted, Created, In Transit, Delivered
        public string CustomerFirstName { get; set; } = string.Empty;
        public string CustomerLastName { get; set; } = string.Empty;
        public string City { get; set; } = string.Empty;
        public string State { get; set; } = string.Empty;
        public required List<LineEntry> LineEntries { get; set; } = [];
    }
}
