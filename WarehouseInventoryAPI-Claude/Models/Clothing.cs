using WarehouseInventory_Claude.Models.Interfaces;

namespace WarehouseInventory_Claude.Models
{
    public class Clothing : IItem
    {
        public string PartitionKey { get; set; } = string.Empty;
        public string RowKey { get; set; } = string.Empty;
        public string LocationId { get; set; } = string.Empty;
        public required string SKUMarker { get; set; }
        public required DateTime UnloadedDate { get; set; }
        public bool Projected { get; set; } = true;
    }

}