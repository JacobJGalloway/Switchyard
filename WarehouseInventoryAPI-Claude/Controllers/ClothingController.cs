using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class ClothingController(IClothingService clothingService) : ControllerBase
    {
        private readonly IClothingService _clothingService = clothingService;

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Clothing>>> GetAll()
        {
            return Ok(await _clothingService.GetAllAsync());
        }

        [HttpGet("{skuId}")]
        public async Task<ActionResult<List<Clothing>>> GetBySKUId(string skuId)
        {
            var items = await _clothingService.GetBySKUIdAsync(skuId);
            if (items.Count == 0) return Ok(new List<Clothing>());
            return Ok(items);
        }

        [HttpGet("location/{locationId}")]
        public async Task<ActionResult<List<Clothing>>> GetByLocationAsync(string locationId)
        {
            return Ok(await _clothingService.GetByLocationAsync(locationId));
        }

        [HttpGet("filter")]
        public async Task<ActionResult<List<Clothing>>> GetByLocationAndSKU([FromQuery] string locationId, [FromQuery] string skuId)
        {
            return Ok(await _clothingService.GetByLocationAndSKUAsync(locationId, skuId));
        }

        [HttpPost]
        public async Task<ActionResult<Clothing>> Create(Clothing item)
        {
            var created = await _clothingService.AddAsync(item);
            return CreatedAtAction(nameof(Create), new { id = created.PartitionKey }, created);
        }

        [HttpPut("{skuId}")]
        public async Task<IActionResult> UpdateBySKUId(string skuId, Clothing item)
        {
            if (skuId != item.SKUMarker) return BadRequest();
            await _clothingService.UpdateBySKUIdAsync(skuId, item);
            return NoContent();
        }

        [HttpPatch("item/{partitionKey}")]
        public async Task<IActionResult> Patch(string partitionKey, [FromBody] InventoryPatchRequest request)
        {
            await _clothingService.PatchAsync(partitionKey, request.Projected, request.UnloadedDate);
            return NoContent();
        }

        [HttpDelete("item/{partitionKey}")]
        public async Task<IActionResult> DeleteByPartitionKey(string partitionKey)
        {
            if (!await _clothingService.DeleteByPartitionKeyAsync(partitionKey)) return NotFound();
            return NoContent();
        }
    }
}
