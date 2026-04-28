using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Controllers;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Tests.Controllers
{
    public class WarehouseControllerTests
    {
        private readonly Mock<IUnitOfWork> _mockUoW;
        private readonly Mock<IWarehouseRepository> _mockWarehouseRepo;
        private readonly WarehouseController _controller;

        public WarehouseControllerTests()
        {
            _mockUoW = new Mock<IUnitOfWork>();
            _mockWarehouseRepo = new Mock<IWarehouseRepository>();
            _mockUoW.Setup(u => u.Warehouses).Returns(_mockWarehouseRepo.Object);
            _controller = new WarehouseController(_mockUoW.Object);
        }

        [Fact]
        public async Task GetAllAsync_ReturnsOk_WithWarehouses()
        {
            var warehouses = new List<Warehouse>
            {
                new() { WarehouseId = "WH001", City = "Chicago", State = "IL" },
                new() { WarehouseId = "WH002", City = "Indianapolis", State = "IN" },
            };
            _mockWarehouseRepo.Setup(r => r.GetAllAsync()).ReturnsAsync(warehouses);

            var result = await _controller.GetAllAsync();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(warehouses, ok.Value);
        }
    }
}
