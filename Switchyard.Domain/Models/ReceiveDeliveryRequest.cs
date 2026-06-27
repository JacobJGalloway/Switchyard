namespace Switchyard.Domain;

public class ReceiveDeliveryRequest
{
    public required string LocationId { get; set; }
    public required List<DeliveryLineItem> Items { get; set; }
}

public class DeliveryLineItem
{
    public required string SKUMarker { get; set; }
    public required int Quantity { get; set; }
}
