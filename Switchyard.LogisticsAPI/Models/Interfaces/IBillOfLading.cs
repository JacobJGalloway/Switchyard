namespace Switchyard.LogisticsAPI.Models.Interfaces
{
    public interface IBillOfLading
    {
        string PartitionKey { get; set; }
        string TransactionId { get; set; }
        string CustomerFirstName { get; set; }
        string CustomerLastName { get; set; }
        string City { get; set; }
        string State { get; set; }
        List<LineEntry> LineEntries { get; set; }
    }
}
