using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;

namespace Switchyard.LogisticsAPI.Controllers;

[ApiController]
[Route("api/analytics")]
public class AnalyticsController(LogisticsContext db) : ControllerBase
{
    [HttpGet("sku-movement")]
    public async Task<IActionResult> SkuMovement(
        [FromQuery] int days = 14,
        [FromQuery] string? warehouses = null,
        [FromQuery] string? skus = null)
    {
        var cutoff = DateTime.UtcNow.Date.AddDays(-days);

        // Fetch base data first — Split() calls stay out of the EF Core expression tree
        // to avoid a .NET 10 / EF Core funcletizer bug with ReadOnlySpan<string> overloads.
        var baseRows = await (
            from le in db.LineEntries
            join bol in db.BillsOfLading on le.TransactionId equals bol.TransactionId
            where bol.CommittedDate != null && bol.CommittedDate >= cutoff
            select new { le, bol }
        ).ToListAsync();

        var whFilter  = warehouses?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);
        var skuFilter = skus?.Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries);

        var filtered = baseRows.AsEnumerable();
        if (whFilter  is { Length: > 0 }) filtered = filtered.Where(x => whFilter.Contains(x.le.LocationId));
        if (skuFilter is { Length: > 0 }) filtered = filtered.Where(x => skuFilter.Contains(x.le.SKUMarker));

        var rows = filtered
            .GroupBy(x => new { x.bol.CommittedDate!.Value.Date, x.le.SKUMarker })
            .Select(g => new
            {
                date        = g.Key.Date,
                skuMarker   = g.Key.SKUMarker,
                quantityMoved = g.Sum(x => Math.Abs(x.le.Quantity))
            })
            .OrderBy(r => r.date)
            .ThenBy(r => r.skuMarker)
            .ToList();

        return Ok(rows);
    }
}
