using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;

namespace Switchyard.LogisticsAPI.Controllers;

[ApiController]
[Route("api/[controller]")]
public class AnalyticsController(LogisticsContext db) : ControllerBase
{
    [HttpGet("sku-movement")]
    public async Task<IActionResult> SkuMovement(
        [FromQuery] int days = 30,
        [FromQuery] string? warehouses = null,
        [FromQuery] string? skus = null)
    {
        var cutoff = DateTime.UtcNow.Date.AddDays(-days);

        var warehouseList = warehouses?
            .Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
        var skuList = skus?
            .Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

        var query =
            from le in db.LineEntries
            join bol in db.BillsOfLading on le.TransactionId equals bol.TransactionId
            where bol.CommittedDate != null && bol.CommittedDate >= cutoff
            select new { le, bol };

        if (warehouseList is { Length: > 0 })
            query = query.Where(x => warehouseList.Contains(x.le.LocationId));

        if (skuList is { Length: > 0 })
            query = query.Where(x => skuList.Contains(x.le.SKUMarker));

        var rows = await query
            .GroupBy(x => new
            {
                Date = x.bol.CommittedDate!.Value.Date,
                x.le.SKUMarker
            })
            .Select(g => new
            {
                date = g.Key.Date,
                skuMarker = g.Key.SKUMarker,
                quantityMoved = g.Sum(x => Math.Abs(x.le.Quantity))
            })
            .OrderBy(r => r.date)
            .ThenBy(r => r.skuMarker)
            .ToListAsync();

        return Ok(rows);
    }
}
