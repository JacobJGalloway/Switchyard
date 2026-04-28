using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data;

namespace Switchyard.LogisticsAPI.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class AuditController(LogisticsContext writeCtx, LogisticsReadContext readCtx) : ControllerBase
    {
        /// <summary>Returns write vs read DB row counts per table, with an InSync flag for each.</summary>
        [HttpGet]
        public async Task<IActionResult> GetAsync()
        {
            var results = new[]
            {
                new { Table = "BillsOfLading", Write = await writeCtx.BillsOfLading.CountAsync(), Read = await readCtx.BillsOfLading.CountAsync() },
                new { Table = "LineEntries",   Write = await writeCtx.LineEntries.CountAsync(),   Read = await readCtx.LineEntries.CountAsync() },
                new { Table = "Warehouses",    Write = await writeCtx.Warehouses.CountAsync(),    Read = await readCtx.Warehouses.CountAsync() },
                new { Table = "Stores",        Write = await writeCtx.Stores.CountAsync(),        Read = await readCtx.Stores.CountAsync() },
            };

            return Ok(results.Select(r => new
            {
                r.Table,
                r.Write,
                r.Read,
                InSync = r.Write == r.Read,
            }));
        }
    }
}
