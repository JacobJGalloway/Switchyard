using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Data.Interfaces;

namespace WarehouseLogistics_Claude.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class StoreController(IUnitOfWork unitOfWork) : ControllerBase
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Store>>> GetAllAsync()
        {
            return Ok(await _unitOfWork.Stores.GetAllAsync());
        }
    }
}
