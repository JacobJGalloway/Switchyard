using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Services.Interfaces;

namespace Switchyard.LogisticsAPI.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class BillOfLadingController(IBillOfLadingService billOfLadingService) : ControllerBase
    {
        private readonly IBillOfLadingService _billOfLadingService = billOfLadingService;

        /// <summary>Returns all bills of lading.</summary>
        [HttpGet]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<IEnumerable<BillOfLading>>> GetAllAsync()
            => Ok(await _billOfLadingService.GetAllAsync());

        /// <summary>Returns a single BOL with its line entries, or 404 if not found.</summary>
        [HttpGet("{transactionId}")]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<BillOfLading>> GetByTransactionIdAsync(string transactionId)
        {
            var bol = await _billOfLadingService.GetByTransactionIdAsync(transactionId);
            if (bol is null) return NotFound();
            return Ok(bol);
        }

        /// <summary>Returns the line entries for a given BOL transaction.</summary>
        [HttpGet("{transactionId}/line-entry")]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<List<LineEntry>>> GetLineEntriesAsync(string transactionId)
            => Ok(await _billOfLadingService.GetLineEntriesByTransactionIdAsync(transactionId));

        /// <summary>Creates a new BOL, generates a transaction ID, persists line entries, and writes a document to Downloads.</summary>
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

        /// <summary>Marks all line entries for the given location on a BOL as processed.</summary>
        [HttpPost("{transactionId}/process/{locationId}")]
        [Authorize(Policy = "ModifyBOL")]
        public async Task<IActionResult> ProcessLocationStop(string transactionId, string locationId)
        {
            await _billOfLadingService.ProcessLocationStop(transactionId, locationId);
            return NoContent();
        }

        /// <summary>Atomically moves unprocessed line entries from one location to another. Returns 409 if the old location is already processed.</summary>
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
