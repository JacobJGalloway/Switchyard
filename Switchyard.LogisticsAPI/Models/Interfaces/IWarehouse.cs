namespace Switchyard.LogisticsAPI.Models.Interfaces
{
    public interface IWarehouse
    {
        string PartitionKey { get; set; }
        string WarehouseId { get; set; }
        string City { get; set; }
        string State { get; set; }
    }
}