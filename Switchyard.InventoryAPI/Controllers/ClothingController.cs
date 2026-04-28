using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Services.Interfaces;

namespace Switchyard.InventoryAPI.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class ClothingController(IClothingService clothingService) : ControllerBase
    {
        private readonly IClothingService _clothingService = clothingService;

        /// <summary>Returns all clothing inventory items.</summary>
        [HttpGet]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<IEnumerable<Clothing>>> GetAll()
        {
            return Ok(await _clothingService.GetAllAsync());
        }

        /// <summary>Returns clothing items matching the given SKU.</summary>
        [HttpGet("{skuId}")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<Clothing>>> GetBySKUId(string skuId)
        {
            var items = await _clothingService.GetBySKUIdAsync(skuId);
            if (items.Count == 0) return Ok(new List<Clothing>());
            return Ok(items);
        }

        /// <summary>Returns clothing items at the given location.</summary>
        [HttpGet("location/{locationId}")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<Clothing>>> GetByLocationAsync(string locationId)
        {
            return Ok(await _clothingService.GetByLocationAsync(locationId));
        }

        /// <summary>Returns clothing items filtered by location and SKU.</summary>
        [HttpGet("filter")]
        [Authorize(Policy = "ReadInventory")]
        public async Task<ActionResult<List<Clothing>>> GetByLocationAndSKU([FromQuery] string locationId, [FromQuery] string skuId)
        {
            return Ok(await _clothingService.GetByLocationAndSKUAsync(locationId, skuId));
        }

        /// <summary>Stages a new clothing item for the current warehouse.</summary>
        [HttpPost]
        public async Task<ActionResult<Clothing>> Create(Clothing item)
        {
            var created = await _clothingService.AddAsync(item);
            return CreatedAtAction(nameof(Create), new { id = created.PartitionKey }, created);
        }

        /// <summary>Replaces all fields on clothing items matching the given SKU.</summary>
        [HttpPut("{skuId}")]
        public async Task<IActionResult> UpdateBySKUId(string skuId, Clothing item)
        {
            if (skuId != item.SKUMarker) return BadRequest();
            await _clothingService.UpdateBySKUIdAsync(skuId, item);
            return NoContent();
        }

        /// <summary>Partially updates the Projected flag or UnloadedDate on a clothing item.</summary>
        [HttpPatch("item/{partitionKey}")]
        public async Task<IActionResult> Patch(string partitionKey, [FromBody] InventoryPatchRequest request)
        {
            await _clothingService.PatchAsync(partitionKey, request.Projected, request.UnloadedDate);
            return NoContent();
        }

        /// <summary>Deletes a specific clothing item by partition key.</summary>
        [HttpDelete("item/{partitionKey}")]
        public async Task<IActionResult> DeleteByPartitionKey(string partitionKey)
        {
            if (!await _clothingService.DeleteByPartitionKeyAsync(partitionKey)) return NotFound();
            return NoContent();
        }
    }
}
