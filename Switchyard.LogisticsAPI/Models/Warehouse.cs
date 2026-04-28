using Switchyard.LogisticsAPI.Models.Interfaces;

namespace Switchyard.LogisticsAPI.Models
{
    public class Warehouse : IWarehouse
    {
        // PartitionKey is "{WarehouseId}-{GUID}" to ensure uniqueness and 
        // allow for multiple entries for the same warehouse if needed
        // while still enabling efficient querying by WarehouseId.
        // Also allows for future expansion into document databases like 
        // Cosmos DB where partition keys are important for performance.
        public string PartitionKey { get; set; } = string.Empty;
        // WarehouseId is "WH" and a three-digit number which is the WarehouseId
        public required string WarehouseId { get; set; }
        public required string City { get; set; }
        public required string State { get; set; }
    }
}