using Switchyard.Domain.Interfaces;

namespace Switchyard.Domain;

public class LineEntry : ILineEntry
{
    public string PartitionKey { get; set; } = string.Empty;
    public string TransactionId { get; set; } = string.Empty;
    public required string LocationId { get; set; }
    public required string SKUMarker { get; set; }
    public required int Quantity { get; set; }
    public bool IsProcessed { get; set; } = false;
    public DateTime? ProcessedDate { get; set; }
}
