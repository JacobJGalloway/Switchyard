using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using Switchyard.InventoryAPI.Controllers;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Services.Interfaces;

namespace Switchyard.InventoryAPI.Tests.Controllers
{
    public class ClothingControllerTests
    {
        private readonly Mock<IClothingService> _mockService;
        private readonly ClothingController _controller;

        public ClothingControllerTests()
        {
            _mockService = new Mock<IClothingService>();
            _controller = new ClothingController(_mockService.Object);
        }

        [Fact]
        public async Task GetAll_ReturnsOkWithItems()
        {
            var items = new List<Clothing>
            {
                new() { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow },
                new() { PartitionKey = "WH001-CLTH002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "CLTH002", UnloadedDate = DateTime.UtcNow }
            };
            _mockService.Setup(s => s.GetAllAsync()).ReturnsAsync(items);

            var result = await _controller.GetAll();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(items, ok.Value);
        }

        [Fact]
        public async Task GetBySKUId_ReturnsOkWithItems_WhenFound()
        {
            var items = new List<Clothing>
            {
                new() { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow }
            };
            _mockService.Setup(s => s.GetBySKUIdAsync("CLTH001")).ReturnsAsync(items);

            var result = await _controller.GetBySKUId("CLTH001");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(items, ok.Value);
        }

        [Fact]
        public async Task GetBySKUId_ReturnsEmptyList_WhenNotFound()
        {
            _mockService.Setup(s => s.GetBySKUIdAsync("CLTH999")).ReturnsAsync(new List<Clothing>());

            var result = await _controller.GetBySKUId("CLTH999");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            var list = Assert.IsType<List<Clothing>>(ok.Value);
            Assert.Empty(list);
        }

        [Fact]
        public async Task Create_ReturnsCreatedAtAction()
        {
            var item = new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            _mockService.Setup(s => s.AddAsync(item)).ReturnsAsync(item);

            var result = await _controller.Create(item);

            Assert.IsType<CreatedAtActionResult>(result.Result);
        }

        [Fact]
        public async Task UpdateBySKUId_ReturnsBadRequest_WhenSkuMismatch()
        {
            var item = new Clothing { PartitionKey = "WH001-CLTH002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "CLTH002", UnloadedDate = DateTime.UtcNow };

            var result = await _controller.UpdateBySKUId("CLTH001", item);

            Assert.IsType<BadRequestResult>(result);
            _mockService.Verify(s => s.UpdateBySKUIdAsync(It.IsAny<string>(), It.IsAny<Clothing>()), Times.Never);
        }

        [Fact]
        public async Task UpdateBySKUId_ReturnsNoContent_WhenValid()
        {
            var item = new Clothing { PartitionKey = "WH001-CLTH001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "CLTH001", UnloadedDate = DateTime.UtcNow };
            _mockService.Setup(s => s.UpdateBySKUIdAsync("CLTH001", item)).Returns(Task.CompletedTask);

            var result = await _controller.UpdateBySKUId("CLTH001", item);

            Assert.IsType<NoContentResult>(result);
        }
    }
}
