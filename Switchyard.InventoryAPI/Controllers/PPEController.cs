using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Services.Interfaces;

namespace Switchyard.InventoryAPI.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class PPEController(IPPEService ppeService) : ControllerBase
    {
        private readonly IPPEService _ppeService = ppeService;

        /// <summary>Returns all PPE inventory items.</summary>
        [HttpGet]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<IEnumerable<PPE>>> GetAll()
        {
            return Ok(await _ppeService.GetAllAsync());
        }

        /// <summary>Returns PPE items matching the given SKU.</summary>
        [HttpGet("{skuId}")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<PPE>>> GetBySKUId(string skuId)
        {
            var items = await _ppeService.GetBySKUIdAsync(skuId);
            if (items.Count == 0) return Ok(new List<PPE>());
            return Ok(items);
        }

        /// <summary>Returns PPE items at the given location.</summary>
        [HttpGet("location/{locationId}")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<PPE>>> GetByLocationAsync(string locationId)
        {
            return Ok(await _ppeService.GetByLocationAsync(locationId));
        }

        /// <summary>Returns PPE items filtered by location and SKU.</summary>
        [HttpGet("filter")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<PPE>>> GetByLocationAndSKU([FromQuery] string locationId, [FromQuery] string skuId)
        {
            return Ok(await _ppeService.GetByLocationAndSKUAsync(locationId, skuId));
        }

        /// <summary>Stages a new PPE item for the current warehouse.</summary>
        [HttpPost]
        public async Task<ActionResult<PPE>> Create(PPE item)
        {
            var created = await _ppeService.AddAsync(item);
            return CreatedAtAction(nameof(Create), new { id = created.PartitionKey }, created);
        }

        /// <summary>Replaces all fields on PPE items matching the given SKU.</summary>
        [HttpPut("{skuId}")]
        public async Task<IActionResult> UpdateBySKUId(string skuId, PPE item)
        {
            if (skuId != item.SKUMarker) return BadRequest();
            await _ppeService.UpdateBySKUIdAsync(skuId, item);
            return NoContent();
        }

        /// <summary>Partially updates the Projected flag or UnloadedDate on a PPE item.</summary>
        [HttpPatch("item/{partitionKey}")]
        public async Task<IActionResult> Patch(string partitionKey, [FromBody] InventoryPatchRequest request)
        {
            await _ppeService.PatchAsync(partitionKey, request.Projected, request.UnloadedDate);
            return NoContent();
        }

        /// <summary>Deletes a specific PPE item by partition key.</summary>
        [HttpDelete("item/{partitionKey}")]
        public async Task<IActionResult> DeleteByPartitionKey(string partitionKey)
        {
            if (!await _ppeService.DeleteByPartitionKeyAsync(partitionKey)) return NotFound();
            return NoContent();
        }
    }
}
