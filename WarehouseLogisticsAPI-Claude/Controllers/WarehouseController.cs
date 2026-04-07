using Microsoft.AspNetCore.Mvc;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Data.Interfaces;

namespace WarehouseLogistics_Claude.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class WarehouseController(IUnitOfWork unitOfWork) : ControllerBase
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Warehouse>>> GetAllAsync()
        {
            return Ok(await _unitOfWork.Warehouses.GetAllAsync());
        }
    }
}
