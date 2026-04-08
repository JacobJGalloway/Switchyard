namespace WarehouseLogistics_Claude.Models
{
    public class ReplaceStopRequest
    {
        public required string OldLocationId { get; set; }
        public required string NewLocationId { get; set; }
    }
}
