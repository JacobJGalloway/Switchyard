using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Data.Interfaces;

namespace Switchyard.LogisticsAPI.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class WarehouseController(IUnitOfWork unitOfWork) : ControllerBase
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        /// <summary>Returns all warehouses.</summary>
        [HttpGet]
        [Authorize(Policy = "ReadBOL")]
        public async Task<ActionResult<IEnumerable<Warehouse>>> GetAllAsync()
        {
            return Ok(await _unitOfWork.Warehouses.GetAllAsync());
        }
    }
}
