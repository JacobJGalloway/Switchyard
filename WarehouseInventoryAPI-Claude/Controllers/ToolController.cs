using Microsoft.AspNetCore.Mvc;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class ToolController(IToolService toolService) : ControllerBase
    {
        private readonly IToolService _toolService = toolService;

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Tool>>> GetAll()
        {
            return Ok(await _toolService.GetAllAsync());
        }

        [HttpGet("location/{locationId}")]
        public async Task<ActionResult<List<Tool>>> GetByLocationAsync(string locationId)
        {
            return Ok(await _toolService.GetByLocationAsync(locationId));
        }

        [HttpGet("filter")]
        public async Task<ActionResult<List<Tool>>> GetByLocationAndSKU([FromQuery] string locationId, [FromQuery] string skuId)
        {
            return Ok(await _toolService.GetByLocationAndSKUAsync(locationId, skuId));
        }

        [HttpGet("{skuId}")]
        public async Task<ActionResult<List<Tool>>> GetBySKUId(string skuId)
        {
            var items = await _toolService.GetBySKUIdAsync(skuId);
            if (items.Count == 0) return Ok(new List<Tool>());
            return Ok(items);
        }

        [HttpPost]
        public async Task<ActionResult<Tool>> Create(Tool item)
        {
            var created = await _toolService.AddAsync(item);
            return CreatedAtAction(nameof(Create), new { id = created.PartitionKey }, created);
        }

        [HttpPut("{skuId}")]
        public async Task<IActionResult> UpdateBySKUId(string skuId, Tool item)
        {
            if (skuId != item.SKUMarker) return BadRequest();
            await _toolService.UpdateBySKUIdAsync(skuId, item);
            return NoContent();
        }

        [HttpPatch("item/{partitionKey}")]
        public async Task<IActionResult> Patch(string partitionKey, [FromBody] InventoryPatchRequest request)
        {
            await _toolService.PatchAsync(partitionKey, request.Projected, request.UnloadedDate);
            return NoContent();
        }

        [HttpDelete("item/{partitionKey}")]
        public async Task<IActionResult> DeleteByPartitionKey(string partitionKey)
        {
            if (!await _toolService.DeleteByPartitionKeyAsync(partitionKey)) return NotFound();
            return NoContent();
        }
    }
}
