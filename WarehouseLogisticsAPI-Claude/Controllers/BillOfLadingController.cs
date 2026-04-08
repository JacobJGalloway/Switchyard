using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Services.Interfaces;

namespace WarehouseLogistics_Claude.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class BillOfLadingController(IBillOfLadingService billOfLadingService) : ControllerBase
    {
        private readonly IBillOfLadingService _billOfLadingService = billOfLadingService;

        [HttpGet]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<IEnumerable<BillOfLading>>> GetAllAsync()
            => Ok(await _billOfLadingService.GetAllAsync());

        [HttpGet("{transactionId}")]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<BillOfLading>> GetByTransactionIdAsync(string transactionId)
        {
            var bol = await _billOfLadingService.GetByTransactionIdAsync(transactionId);
            if (bol is null) return NotFound();
            return Ok(bol);
        }

        [HttpGet("{transactionId}/line-entry")]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<List<LineEntry>>> GetLineEntriesAsync(string transactionId)
            => Ok(await _billOfLadingService.GetLineEntriesByTransactionIdAsync(transactionId));

        [HttpPost]
        [Authorize(Policy = "CreateBOL")]
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
        [Authorize(Policy = "ModifyBOL")]
        public async Task<IActionResult> ProcessLocationStop(string transactionId, string locationId)
        {
            await _billOfLadingService.ProcessLocationStop(transactionId, locationId);
            return NoContent();
        }

        [HttpPost("{transactionId}/replace-stop")]
        [Authorize(Policy = "ModifyBOL")]
        public async Task<IActionResult> ReplaceLocationStopAsync(string transactionId, [FromBody] ReplaceStopRequest request)
        {
            try
            {
                await _billOfLadingService.ReplaceLocationStopAsync(transactionId, request.OldLocationId, request.NewLocationId);
                return NoContent();
            }
            catch (ArgumentException ex)
            {
                return BadRequest(ex.Message);
            }
            catch (InvalidOperationException ex)
            {
                return Conflict(ex.Message);
            }
        }
    }
}
