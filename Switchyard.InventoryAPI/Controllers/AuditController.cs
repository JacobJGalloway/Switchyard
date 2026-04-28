using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Data;

namespace Switchyard.InventoryAPI.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class AuditController(InventoryContext writeCtx, InventoryReadContext readCtx) : ControllerBase
    {
        /// <summary>Returns write vs read DB row counts for Clothing, PPE, and Tool tables, with an InSync flag per table.</summary>
        [HttpGet]
        public async Task<IActionResult> GetAsync()
        {
            var results = new[]
            {
                new { Table = "Clothing", Write = await writeCtx.Clothing.CountAsync(), Read = await readCtx.Clothing.CountAsync() },
                new { Table = "PPE",      Write = await writeCtx.PPE.CountAsync(),      Read = await readCtx.PPE.CountAsync() },
                new { Table = "Tool",     Write = await writeCtx.Tool.CountAsync(),      Read = await readCtx.Tool.CountAsync() },
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
