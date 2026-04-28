using Xunit;
using Moq;
using Microsoft.AspNetCore.Mvc;
using Switchyard.LogisticsAPI.Controllers;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Tests.Controllers
{
    public class StoreControllerTests
    {
        private readonly Mock<IUnitOfWork> _mockUoW;
        private readonly Mock<IStoreRepository> _mockStoreRepo;
        private readonly StoreController _controller;

        public StoreControllerTests()
        {
            _mockUoW = new Mock<IUnitOfWork>();
            _mockStoreRepo = new Mock<IStoreRepository>();
            _mockUoW.Setup(u => u.Stores).Returns(_mockStoreRepo.Object);
            _controller = new StoreController(_mockUoW.Object);
        }

        [Fact]
        public async Task GetAllAsync_ReturnsOk_WithStores()
        {
            var stores = new List<Store>
            {
                new() { PartitionKey = "ST0001-pk", StoreId = "ST0001", BaseWarehouseId = "WH001", City = "Chicago", State = "IL" },
                new() { PartitionKey = "ST0002-pk", StoreId = "ST0002", BaseWarehouseId = "WH001", City = "Naperville", State = "IL" },
            };
            _mockStoreRepo.Setup(r => r.GetAllAsync()).ReturnsAsync(stores);

            var result = await _controller.GetAllAsync();

            var ok = Assert.IsType<OkObjectResult>(result.Result);
            Assert.Equal(stores, ok.Value);
        }
    }
}
