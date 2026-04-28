namespace Switchyard.LogisticsAPI.Models.Interfaces
{
    public interface IStore
    {
        string PartitionKey { get; set; }
        string StoreId { get; set; }
        string BaseWarehouseId { get; set; }
        string City { get; set; }
        string State { get; set; }
    }
}