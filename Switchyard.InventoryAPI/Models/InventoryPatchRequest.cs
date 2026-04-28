namespace Switchyard.InventoryAPI.Models
{
    public class InventoryPatchRequest
    {
        public bool? Projected { get; set; }
        public DateTime? UnloadedDate { get; set; }
    }
}
