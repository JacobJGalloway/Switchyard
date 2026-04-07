using Microsoft.AspNetCore.Mvc;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class PPEController(IPPEService ppeService) : ControllerBase
    {
        private readonly IPPEService _ppeService = ppeService;

        [HttpGet]
        public async Task<ActionResult<IEnumerable<PPE>>> GetAll()
        {
            return Ok(await _ppeService.GetAllAsync());
        }

        [HttpGet("{skuId}")]
        public async Task<ActionResult<List<PPE>>> GetBySKUId(string skuId)
        {
            var items = await _ppeService.GetBySKUIdAsync(skuId);
            if (items.Count == 0) return Ok(new List<PPE>());
            return Ok(items);
        }

        [HttpGet("location/{locationId}")]
        public async Task<ActionResult<List<PPE>>> GetByLocationAsync(string locationId)
        {
            return Ok(await _ppeService.GetByLocationAsync(locationId));
        }

        [HttpGet("filter")]
        public async Task<ActionResult<List<PPE>>> GetByLocationAndSKU([FromQuery] string locationId, [FromQuery] string skuId)
        {
            return Ok(await _ppeService.GetByLocationAndSKUAsync(locationId, skuId));
        }

        [HttpPost]
        public async Task<ActionResult<PPE>> Create(PPE item)
        {
            var created = await _ppeService.AddAsync(item);
            return CreatedAtAction(nameof(Create), new { id = created.PartitionKey }, created);
        }

        [HttpPut("{skuId}")]
        public async Task<IActionResult> UpdateBySKUId(string skuId, PPE item)
        {
            if (skuId != item.SKUMarker) return BadRequest();
            await _ppeService.UpdateBySKUIdAsync(skuId, item);
            return NoContent();
        }

        [HttpPatch("item/{partitionKey}")]
        public async Task<IActionResult> Patch(string partitionKey, [FromBody] InventoryPatchRequest request)
        {
            await _ppeService.PatchAsync(partitionKey, request.Projected, request.UnloadedDate);
            return NoContent();
        }

        [HttpDelete("item/{partitionKey}")]
        public async Task<IActionResult> DeleteByPartitionKey(string partitionKey)
        {
            if (!await _ppeService.DeleteByPartitionKeyAsync(partitionKey)) return NotFound();
            return NoContent();
        }
    }
}
