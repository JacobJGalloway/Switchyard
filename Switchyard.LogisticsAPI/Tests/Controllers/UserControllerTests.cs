using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Controllers;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Services.Interfaces;

namespace Switchyard.LogisticsAPI.Tests.Controllers
{
    public class UserControllerTests
    {
        private readonly Mock<IUserManagementService> _mockService;
        private readonly UserController _controller;

        public UserControllerTests()
        {
            _mockService = new Mock<IUserManagementService>();
            _controller = new UserController(_mockService.Object);
        }

        private static AppUser MakeUser(string id = "auth0|001") => new()
        {
            UserId = id,
            Email = "test@example.com",
            Name = "Test User",
            Role = "Warehouse",
            WarehouseId = "WH001",
        };

        [Fact]
        public async Task GetAllAsync_ReturnsOk_WithUsers()
        {
            var users = new List<AppUser> { MakeUser("auth0|001"), MakeUser("auth0|002") };
            _mockService.Setup(s => s.GetAllUsersAsync()).ReturnsAsync(users);

            var result = await _controller.GetAllAsync();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(users, ok.Value);
        }

        [Fact]
        public async Task CreateAsync_ReturnsCreatedAtAction_WhenValid()
        {
            var request = new CreateUserRequest
            {
                Email = "new@example.com", Password = "Pass1!", FirstName = "Jane",
                LastName = "Doe", Role = "Warehouse", WarehouseId = "WH001"
            };
            var created = MakeUser("auth0|new");
            _mockService.Setup(s => s.CreateUserAsync(request)).ReturnsAsync(created);

            var result = await _controller.CreateAsync(request);

            Assert.IsType<CreatedAtActionResult>(result);
        }

        [Fact]
        public async Task CreateAsync_ReturnsBadRequest_WhenInvalidRole()
        {
            var request = new CreateUserRequest
            {
                Email = "new@example.com", Password = "Pass1!", FirstName = "Jane",
                LastName = "Doe", Role = "NonExistentRole"
            };
            _mockService.Setup(s => s.CreateUserAsync(request))
                .ThrowsAsync(new ArgumentException("Role 'NonExistentRole' does not exist."));

            var result = await _controller.CreateAsync(request);

            var bad = Assert.IsType<BadRequestObjectResult>(result);
            Assert.Equal("Role 'NonExistentRole' does not exist.", bad.Value);
        }

        [Fact]
        public async Task DeactivateAsync_ReturnsNoContent()
        {
            _mockService.Setup(s => s.DeactivateUserAsync("auth0|001")).Returns(Task.CompletedTask);

            var result = await _controller.DeactivateAsync("auth0|001");

            Assert.IsType<NoContentResult>(result);
            _mockService.Verify(s => s.DeactivateUserAsync("auth0|001"), Times.Once);
        }
    }
}
