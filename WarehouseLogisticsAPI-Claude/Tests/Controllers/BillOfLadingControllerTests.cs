using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using WarehouseLogistics_Claude.Controllers;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Services.Interfaces;

namespace WarehouseLogistics_Claude.Tests.Controllers
{
    public class BillOfLadingControllerTests
    {
        private readonly Mock<IBillOfLadingService> _mockService;
        private readonly BillOfLadingController _controller;

        public BillOfLadingControllerTests()
        {
            _mockService = new Mock<IBillOfLadingService>();
            _controller = new BillOfLadingController(_mockService.Object);
        }

        private static BillOfLading MakeValidBOL() => new()
        {
            Status = "Pending",
            CustomerFirstName = "Jane",
            CustomerLastName = "Doe",
            City = "Springfield",
            State = "IL",
            LineEntries =
            [
                new LineEntry { LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 10 }
            ]
        };

        [Fact]
        public async Task Create_ReturnsCreatedAtAction_WhenValid()
        {
            var bol = MakeValidBOL();
            _mockService.Setup(s => s.CreateAsync(bol)).ReturnsAsync("txn001");

            var result = await _controller.Create(bol);

            Assert.IsType<CreatedAtActionResult>(result);
        }

        [Fact]
        public async Task Create_ReturnsBadRequest_WhenArgumentExceptionThrown()
        {
            var bol = MakeValidBOL();
            _mockService.Setup(s => s.CreateAsync(bol)).ThrowsAsync(new ArgumentException("A Pending Bill of Lading is required."));

            var result = await _controller.Create(bol);

            var bad = Assert.IsType<BadRequestObjectResult>(result);
            Assert.Equal("A Pending Bill of Lading is required.", bad.Value);
        }

        [Fact]
        public async Task ProcessLocationStop_ReturnsNoContent()
        {
            _mockService.Setup(s => s.ProcessLocationStop("txn001", "WH001")).Returns(Task.CompletedTask);

            var result = await _controller.ProcessLocationStop("txn001", "WH001");

            Assert.IsType<NoContentResult>(result);
            _mockService.Verify(s => s.ProcessLocationStop("txn001", "WH001"), Times.Once);
        }
    }
}
