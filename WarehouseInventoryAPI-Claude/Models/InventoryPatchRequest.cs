namespace WarehouseInventory_Claude.Models
{
    public class InventoryPatchRequest
    {
        public bool? Projected { get; set; }
        public DateTime? UnloadedDate { get; set; }
    }
}
