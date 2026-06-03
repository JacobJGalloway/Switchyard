using Switchyard.Domain.Interfaces;

namespace Switchyard.Domain;

public class Warehouse : IWarehouse
{
    public string PartitionKey { get; set; } = string.Empty;
    public required string WarehouseId { get; set; }
    public required string City { get; set; }
    public required string State { get; set; }
}
