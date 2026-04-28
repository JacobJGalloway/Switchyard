using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using Switchyard.InventoryAPI.Controllers;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Services.Interfaces;

namespace Switchyard.InventoryAPI.Tests.Controllers
{
    public class ToolControllerTests
    {
        private readonly Mock<IToolService> _mockService;
        private readonly ToolController _controller;

        public ToolControllerTests()
        {
            _mockService = new Mock<IToolService>();
            _controller = new ToolController(_mockService.Object);
        }

        [Fact]
        public async Task GetAll_ReturnsOkWithItems()
        {
            var items = new List<Tool>
            {
                new() { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new() { PartitionKey = "WH001-PWTL002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL002", UnloadedDate = DateTime.UtcNow }
            };
            _mockService.Setup(s => s.GetAllAsync()).ReturnsAsync(items);

            var result = await _controller.GetAll();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(items, ok.Value);
        }

        [Fact]
        public async Task GetBySKUId_ReturnsOkWithItems_WhenFound()
        {
            var items = new List<Tool>
            {
                new() { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow }
            };
            _mockService.Setup(s => s.GetBySKUIdAsync("PWTL001")).ReturnsAsync(items);

            var result = await _controller.GetBySKUId("PWTL001");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(items, ok.Value);
        }

        [Fact]
        public async Task GetBySKUId_ReturnsEmptyList_WhenNotFound()
        {
            _mockService.Setup(s => s.GetBySKUIdAsync("PWTL999")).ReturnsAsync(new List<Tool>());

            var result = await _controller.GetBySKUId("PWTL999");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            var list = Assert.IsType<List<Tool>>(ok.Value);
            Assert.Empty(list);
        }

        [Fact]
        public async Task Create_ReturnsCreatedAtAction()
        {
            var item = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _mockService.Setup(s => s.AddAsync(item)).ReturnsAsync(item);

            var result = await _controller.Create(item);

            Assert.IsType<CreatedAtActionResult>(result.Result);
        }

        [Fact]
        public async Task UpdateBySKUId_ReturnsBadRequest_WhenSkuMismatch()
        {
            var item = new Tool { PartitionKey = "WH001-PWTL002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL002", UnloadedDate = DateTime.UtcNow };

            var result = await _controller.UpdateBySKUId("PWTL001", item);

            Assert.IsType<BadRequestResult>(result);
            _mockService.Verify(s => s.UpdateBySKUIdAsync(It.IsAny<string>(), It.IsAny<Tool>()), Times.Never);
        }

        [Fact]
        public async Task UpdateBySKUId_ReturnsNoContent_WhenValid()
        {
            var item = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _mockService.Setup(s => s.UpdateBySKUIdAsync("PWTL001", item)).Returns(Task.CompletedTask);

            var result = await _controller.UpdateBySKUId("PWTL001", item);

            Assert.IsType<NoContentResult>(result);
        }
    }
}
