using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Models
{
    public class Tool : IItem
    {
        public string PartitionKey { get; set; } = string.Empty;
        public string RowKey { get; set; } = string.Empty;
        public string SKUMarker { get; set; } = string.Empty;
        public DateTime UnloadedDate { get; set; } = DateTime.UtcNow;
        public string Type { get; set; } = string.Empty;
        public double Size { get; set; }
        public string Description { get; set; } = string.Empty;
    }
}