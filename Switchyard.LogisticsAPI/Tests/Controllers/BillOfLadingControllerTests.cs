using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Controllers;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Services.Interfaces;

namespace Switchyard.LogisticsAPI.Tests.Controllers
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

        [Fact]
        public async Task GetAllAsync_ReturnsOk_WithBOLs()
        {
            var bols = new List<BillOfLading> { MakeValidBOL(), MakeValidBOL() };
            _mockService.Setup(s => s.GetAllAsync()).ReturnsAsync(bols);

            var result = await _controller.GetAllAsync();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(bols, ok.Value);
        }

        [Fact]
        public async Task GetByTransactionIdAsync_ReturnsOk_WhenFound()
        {
            var bol = MakeValidBOL();
            _mockService.Setup(s => s.GetByTransactionIdAsync("txn001")).ReturnsAsync(bol);

            var result = await _controller.GetByTransactionIdAsync("txn001");

            Assert.IsType<OkObjectResult>(result.Result);
        }

        [Fact]
        public async Task GetByTransactionIdAsync_ReturnsNotFound_WhenMissing()
        {
            _mockService.Setup(s => s.GetByTransactionIdAsync("txn999")).ReturnsAsync((BillOfLading?)null);

            var result = await _controller.GetByTransactionIdAsync("txn999");

            Assert.IsType<NotFoundResult>(result.Result);
        }

        [Fact]
        public async Task GetLineEntriesAsync_ReturnsOk_WithEntries()
        {
            var entries = new List<LineEntry>
            {
                new() { TransactionId = "txn001", LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 5 }
            };
            _mockService.Setup(s => s.GetLineEntriesByTransactionIdAsync("txn001")).ReturnsAsync(entries);

            var result = await _controller.GetLineEntriesAsync("txn001");

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(entries, ok.Value);
        }

        [Fact]
        public async Task ReplaceLocationStopAsync_ReturnsNoContent_WhenValid()
        {
            var request = new ReplaceStopRequest { OldLocationId = "WH001", NewLocationId = "ST0001" };
            _mockService.Setup(s => s.ReplaceLocationStopAsync("txn001", "WH001", "ST0001")).Returns(Task.CompletedTask);

            var result = await _controller.ReplaceLocationStopAsync("txn001", request);

            Assert.IsType<NoContentResult>(result);
        }

        [Fact]
        public async Task ReplaceLocationStopAsync_ReturnsBadRequest_WhenLocationNotFound()
        {
            var request = new ReplaceStopRequest { OldLocationId = "ST9999", NewLocationId = "ST0001" };
            _mockService.Setup(s => s.ReplaceLocationStopAsync("txn001", "ST9999", "ST0001"))
                .ThrowsAsync(new ArgumentException("No line entries found for location ST9999 on transaction txn001."));

            var result = await _controller.ReplaceLocationStopAsync("txn001", request);

            var bad = Assert.IsType<BadRequestObjectResult>(result);
            Assert.Contains("ST9999", bad.Value!.ToString());
        }

        [Fact]
        public async Task ReplaceLocationStopAsync_ReturnsConflict_WhenAlreadyProcessed()
        {
            var request = new ReplaceStopRequest { OldLocationId = "WH001", NewLocationId = "ST0001" };
            _mockService.Setup(s => s.ReplaceLocationStopAsync("txn001", "WH001", "ST0001"))
                .ThrowsAsync(new InvalidOperationException("Location WH001 has already been processed and cannot be replaced."));

            var result = await _controller.ReplaceLocationStopAsync("txn001", request);

            var conflict = Assert.IsType<ConflictObjectResult>(result);
            Assert.Contains("WH001", conflict.Value!.ToString());
        }
    }
}
