using Switchyard.Domain.Interfaces;

namespace Switchyard.Domain;

public class BillOfLading : IBillOfLading
{
    public string PartitionKey { get; set; } = string.Empty;
    public string TransactionId { get; set; } = string.Empty;
    public string Status { get; set; } = "Pending";
    public string CustomerFirstName { get; set; } = string.Empty;
    public string CustomerLastName { get; set; } = string.Empty;
    public string City { get; set; } = string.Empty;
    public string State { get; set; } = string.Empty;
    public required List<LineEntry> LineEntries { get; set; } = [];
}
