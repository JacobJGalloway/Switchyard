namespace Switchyard.LogisticsAPI.Models
{
    public class ReplaceStopRequest
    {
        public required string OldLocationId { get; set; }
        public required string NewLocationId { get; set; }
    }
}
