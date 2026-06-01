using Switchyard.Domain.Interfaces;

namespace Switchyard.Domain;

public class Store : IStore
{
    public string PartitionKey { get; set; } = string.Empty;
    public required string StoreId { get; set; }
    public required string BaseWarehouseId { get; set; }
    public required string City { get; set; }
    public required string State { get; set; }
}
