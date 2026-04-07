using Microsoft.AspNetCore.Mvc;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Services.Interfaces;

namespace WarehouseLogistics_Claude.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class BillOfLadingController(IBillOfLadingService billOfLadingService) : ControllerBase
    {
        private readonly IBillOfLadingService _billOfLadingService = billOfLadingService;

        [HttpPost]
        public async Task<IActionResult> Create([FromBody] BillOfLading billOfLading)
        {
            try
            {
                var transactionId = await _billOfLadingService.CreateAsync(billOfLading);
                return CreatedAtAction(nameof(Create), new { transactionId }, transactionId);
            }
            catch (ArgumentException ex)
            {
                return BadRequest(ex.Message);
            }
        }

        [HttpPost("{transactionId}/process/{locationId}")]
        public async Task<IActionResult> ProcessLocationStop(string transactionId, string locationId)
        {
            await _billOfLadingService.ProcessLocationStop(transactionId, locationId);
            return NoContent();
        }
    }
}
