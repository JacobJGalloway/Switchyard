namespace Switchyard.LogisticsAPI.Models.Interfaces
{
    public interface ILineEntry
    {
        string PartitionKey { get; set; }
        string TransactionId { get; set; }
        string LocationId { get; set; }
        string SKUMarker { get; set; }
        int Quantity { get; set; }
        bool IsProcessed { get; set; }
    }
}
