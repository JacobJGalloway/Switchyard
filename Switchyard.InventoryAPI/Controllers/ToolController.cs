using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Services.Interfaces;

namespace Switchyard.InventoryAPI.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class ToolController(IToolService toolService) : ControllerBase
    {
        private readonly IToolService _toolService = toolService;

        /// <summary>Returns all tool inventory items.</summary>
        [HttpGet]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<IEnumerable<Tool>>> GetAll()
        {
            return Ok(await _toolService.GetAllAsync());
        }

        /// <summary>Returns tool items at the given location.</summary>
        [HttpGet("location/{locationId}")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<Tool>>> GetByLocationAsync(string locationId)
        {
            return Ok(await _toolService.GetByLocationAsync(locationId));
        }

        /// <summary>Returns tool items filtered by location and SKU.</summary>
        [HttpGet("filter")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<Tool>>> GetByLocationAndSKU([FromQuery] string locationId, [FromQuery] string skuId)
        {
            return Ok(await _toolService.GetByLocationAndSKUAsync(locationId, skuId));
        }

        /// <summary>Returns tool items matching the given SKU.</summary>
        [HttpGet("{skuId}")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<Tool>>> GetBySKUId(string skuId)
        {
            var items = await _toolService.GetBySKUIdAsync(skuId);
            if (items.Count == 0) return Ok(new List<Tool>());
            return Ok(items);
        }

        /// <summary>Stages a new tool item for the current warehouse.</summary>
        [HttpPost]
        public async Task<ActionResult<Tool>> Create(Tool item)
        {
            var created = await _toolService.AddAsync(item);
            return CreatedAtAction(nameof(Create), new { id = created.PartitionKey }, created);
        }

        /// <summary>Replaces all fields on tool items matching the given SKU.</summary>
        [HttpPut("{skuId}")]
        public async Task<IActionResult> UpdateBySKUId(string skuId, Tool item)
        {
            if (skuId != item.SKUMarker) return BadRequest();
            await _toolService.UpdateBySKUIdAsync(skuId, item);
            return NoContent();
        }

        /// <summary>Partially updates the Projected flag or UnloadedDate on a tool item.</summary>
        [HttpPatch("item/{partitionKey}")]
        public async Task<IActionResult> Patch(string partitionKey, [FromBody] InventoryPatchRequest request)
        {
            await _toolService.PatchAsync(partitionKey, request.Projected, request.UnloadedDate);
            return NoContent();
        }

        /// <summary>Deletes a specific tool item by partition key.</summary>
        [HttpDelete("item/{partitionKey}")]
        public async Task<IActionResult> DeleteByPartitionKey(string partitionKey)
        {
            if (!await _toolService.DeleteByPartitionKeyAsync(partitionKey)) return NotFound();
            return NoContent();
        }
    }
}
