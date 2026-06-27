namespace Switchyard.Domain.Interfaces;

public interface IItem
{
    string PartitionKey { get; set; }
    string RowKey { get; set; }
    string LocationId { get; set; }
    string SKUMarker { get; set; }
    DateTime UnloadedDate { get; set; }
    bool Projected { get; set; }
    decimal UnitPrice { get; set; }
    string PriceCurrency { get; set; }
    DateTime? PriceEffectiveDate { get; set; }
}
