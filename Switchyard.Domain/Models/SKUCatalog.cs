namespace Switchyard.Domain;

public class SKUCatalog
{
    public string SKUMarker { get; set; } = string.Empty;
    public string Category { get; set; } = string.Empty;
    public string Type { get; set; } = string.Empty;
    public string Color { get; set; } = string.Empty;
    public string Size { get; set; } = string.Empty;
    public decimal UnitPrice { get; set; }
}
