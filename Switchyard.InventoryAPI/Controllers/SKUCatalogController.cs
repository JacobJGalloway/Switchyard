using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data;

namespace Switchyard.InventoryAPI.Controllers;

[ApiController]
[Route("api/sku-catalog")]
public class SKUCatalogController(InventoryContext db) : ControllerBase
{
    [HttpGet]
    public async Task<IActionResult> Get([FromQuery] string? skus)
    {
        var query = db.SKUCatalog.AsQueryable();

        if (!string.IsNullOrWhiteSpace(skus))
        {
            var skuList = skus.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
            query = query.Where(s => skuList.Contains(s.SKUMarker));
        }

        var result = await query
            .OrderBy(s => s.Category)
            .ThenBy(s => s.SKUMarker)
            .Select(s => new { skuMarker = s.SKUMarker, category = s.Category, type = s.Type, color = s.Color, size = s.Size, unitPrice = s.UnitPrice })
            .ToListAsync();

        return Ok(result);
    }
}
