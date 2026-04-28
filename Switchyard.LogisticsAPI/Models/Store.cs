using Switchyard.LogisticsAPI.Models.Interfaces;

namespace Switchyard.LogisticsAPI.Models
{
    public class Store : Models.Interfaces.IStore
    {
        // PartitionKey is "{StoreId}-{BaseWarehouseId}-{GUID}". Allows multiple
        // records per store and future Cosmos DB compatibility.
        public string PartitionKey { get; set; } = string.Empty;
        // StoreId is "ST" followed by a four-digit number, e.g. "ST0001"
        public required string StoreId { get; set; }
        // BaseWarehouseId is the supplying warehouse for this store (e.g. "WH001")
        public required string BaseWarehouseId { get; set; }
        public required string City { get; set; }
        public required string State { get; set; }
    }
}