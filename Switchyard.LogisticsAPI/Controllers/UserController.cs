using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Services.Interfaces;

namespace Switchyard.LogisticsAPI.Controllers
{
    // TODO: When the data layer is extracted into its own service,
    //       move this controller to that project.
    [Authorize]
    [ApiController]
    [Route("api/[controller]")]
    public class UserController(IUserManagementService userManagementService) : ControllerBase
    {
        private readonly IUserManagementService _userManagementService = userManagementService;

        /// <summary>Returns all Auth0 users with their assigned roles and location metadata.</summary>
        [HttpGet]
        [Authorize(Policy = "ManageUsers")]
        public async Task<ActionResult<IEnumerable<AppUser>>> GetAllAsync()
            => Ok(await _userManagementService.GetAllUsersAsync());

        /// <summary>Creates a new username/password user in Auth0 and assigns a role.</summary>
        [HttpPost]
        [Authorize(Policy = "ManageUsers")]
        public async Task<IActionResult> CreateAsync([FromBody] CreateUserRequest request)
        {
            try
            {
                var user = await _userManagementService.CreateUserAsync(request);
                return CreatedAtAction(nameof(CreateAsync), new { userId = user.UserId }, user);
            }
            catch (ArgumentException ex)
            {
                return BadRequest(ex.Message);
            }
        }

        /// <summary>Blocks a user in Auth0 (soft deactivation; does not delete the account).</summary>
        [HttpPatch("{userId}/deactivate")]
        [Authorize(Policy = "ManageUsers")]
        public async Task<IActionResult> DeactivateAsync(string userId)
        {
            await _userManagementService.DeactivateUserAsync(userId);
            return NoContent();
        }
    }
}
