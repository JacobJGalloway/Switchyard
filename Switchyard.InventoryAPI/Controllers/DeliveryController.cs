using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.Domain;
using Switchyard.InventoryAPI.Data.Interfaces;

namespace Switchyard.InventoryAPI.Controllers;

[Authorize]
[ApiController]
[Route("api/deliveries")]
public class DeliveryController(IUnitOfWork unitOfWork) : ControllerBase
{
    private readonly IUnitOfWork _unitOfWork = unitOfWork;

    /// <summary>Confirms receipt of a delivery: sets Projected=false and LocationId on matching in-transit inventory items.</summary>
    [HttpPost("receive")]
    public async Task<IActionResult> ReceiveDelivery([FromBody] ReceiveDeliveryRequest request)
    {
        if (request.Items.Count == 0) return BadRequest("No items specified.");
        await _unitOfWork.ReceiveDeliveryAsync(request.LocationId, request.Items);
        return NoContent();
    }
}
