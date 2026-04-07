using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using WarehouseInventory_Claude.Controllers;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Tests.Controllers
{
    public class PPEControllerTests
    {
        private readonly Mock<IPPEService> _mockService;
        private readonly PPEController _controller;

        public PPEControllerTests()
        {
            _mockService = new Mock<IPPEService>();
            _controller = new PPEController(_mockService.Object);
        }

        [Fact]
        public async Task GetAll_ReturnsOkWithItems()
        {
            var items = new List<PPE>
            {
                new() { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow },
                new() { PartitionKey = "WH001-SPPE002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "SPPE002", UnloadedDate = DateTime.UtcNow }
            };
            _mockService.Setup(s => s.GetAllAsync()).ReturnsAsync(items);

            var result = await _controller.GetAll();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(items, ok.Value);
        }

        [Fact]
        public async Task GetBySKUId_ReturnsOkWithItems_WhenFound()
        {
            var items = new List<PPE>
            {
                new() { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow }
            };
            _mockService.Setup(s => s.GetBySKUIdAsync("SPPE001")).ReturnsAsync(items);

            var result = await _controller.GetBySKUId("SPPE001");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(items, ok.Value);
        }

        [Fact]
        public async Task GetBySKUId_ReturnsEmptyList_WhenNotFound()
        {
            _mockService.Setup(s => s.GetBySKUIdAsync("SPPE999")).ReturnsAsync(new List<PPE>());

            var result = await _controller.GetBySKUId("SPPE999");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            var list = Assert.IsType<List<PPE>>(ok.Value);
            Assert.Empty(list);
        }

        [Fact]
        public async Task Create_ReturnsCreatedAtAction()
        {
            var item = new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            _mockService.Setup(s => s.AddAsync(item)).ReturnsAsync(item);

            var result = await _controller.Create(item);

            Assert.IsType<CreatedAtActionResult>(result.Result);
        }

        [Fact]
        public async Task UpdateBySKUId_ReturnsBadRequest_WhenSkuMismatch()
        {
            var item = new PPE { PartitionKey = "WH001-SPPE002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "SPPE002", UnloadedDate = DateTime.UtcNow };

            var result = await _controller.UpdateBySKUId("SPPE001", item);

            Assert.IsType<BadRequestResult>(result);
            _mockService.Verify(s => s.UpdateBySKUIdAsync(It.IsAny<string>(), It.IsAny<PPE>()), Times.Never);
        }

        [Fact]
        public async Task UpdateBySKUId_ReturnsNoContent_WhenValid()
        {
            var item = new PPE { PartitionKey = "WH001-SPPE001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "SPPE001", UnloadedDate = DateTime.UtcNow };
            _mockService.Setup(s => s.UpdateBySKUIdAsync("SPPE001", item)).Returns(Task.CompletedTask);

            var result = await _controller.UpdateBySKUId("SPPE001", item);

            Assert.IsType<NoContentResult>(result);
        }
    }
}
