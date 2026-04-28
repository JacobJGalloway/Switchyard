using Xunit;
using Moq;
using Switchyard.InventoryAPI.Data.Interfaces;
using Switchyard.InventoryAPI.Models;
using Switchyard.InventoryAPI.Services;

namespace Switchyard.InventoryAPI.Tests.Services
{
    public class ToolServiceTests
    {
        private readonly Mock<IUnitOfWork> _mockUoW;
        private readonly Mock<IToolRepository> _mockRepo;
        private readonly ToolService _service;

        public ToolServiceTests()
        {
            _mockUoW = new Mock<IUnitOfWork>();
            _mockRepo = new Mock<IToolRepository>();
            _mockUoW.Setup(u => u.Tools).Returns(_mockRepo.Object);
            _service = new ToolService(_mockUoW.Object);
        }

        [Fact]
        public async Task GetAllAsync_ReturnsAllItems()
        {
            var items = new List<Tool>
            {
                new() { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow },
                new() { PartitionKey = "WH001-PWTL002-b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7", SKUMarker = "PWTL002", UnloadedDate = DateTime.UtcNow }
            };
            _mockRepo.Setup(r => r.GetAllAsync()).ReturnsAsync(items);

            var result = await _service.GetAllAsync();

            Assert.Equal(items, result);
        }

        [Fact]
        public async Task GetBySKUIdAsync_ReturnsItems()
        {
            var items = new List<Tool>
            {
                new() { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow }
            };
            _mockUoW.Setup(u => u.GetToolBySKUIdAsync("PWTL001")).ReturnsAsync(items);

            var result = await _service.GetBySKUIdAsync("PWTL001");

            Assert.Equal(items, result);
        }

        [Fact]
        public async Task GetByLocationAndSKUAsync_ReturnsItems()
        {
            var items = new List<Tool>
            {
                new() { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", LocationId = "WH001", UnloadedDate = DateTime.UtcNow }
            };
            _mockRepo.Setup(r => r.GetByLocationAndSKUAsync("WH001", "PWTL001")).ReturnsAsync(items);

            var result = await _service.GetByLocationAndSKUAsync("WH001", "PWTL001");

            Assert.Equal(items, result);
        }

        [Fact]
        public async Task AddAsync_ReturnsCreatedItem_AndSavesChanges()
        {
            var item = new Tool { SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _mockRepo.Setup(r => r.AddAsync(item)).ReturnsAsync(item);

            var result = await _service.AddAsync(item);

            Assert.Equal(item, result);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Once);
        }

        [Fact]
        public async Task UpdateBySKUIdAsync_CallsRepoAndSavesChanges()
        {
            var item = new Tool { PartitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6", SKUMarker = "PWTL001", UnloadedDate = DateTime.UtcNow };
            _mockRepo.Setup(r => r.UpdateBySKUIdAsync("PWTL001", item)).Returns(Task.CompletedTask);

            await _service.UpdateBySKUIdAsync("PWTL001", item);

            _mockRepo.Verify(r => r.UpdateBySKUIdAsync("PWTL001", item), Times.Once);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Once);
        }

        [Fact]
        public async Task PatchAsync_CallsRepoAndSavesChanges()
        {
            var partitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6";
            _mockRepo.Setup(r => r.PatchAsync(partitionKey, false, It.IsAny<DateTime?>())).Returns(Task.CompletedTask);

            await _service.PatchAsync(partitionKey, false, null);

            _mockRepo.Verify(r => r.PatchAsync(partitionKey, false, null), Times.Once);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Once);
        }

        [Fact]
        public async Task DeleteByPartitionKeyAsync_ReturnsTrue_AndSavesChanges_WhenFound()
        {
            var partitionKey = "WH001-PWTL001-a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6";
            _mockRepo.Setup(r => r.DeleteByPartitionKeyAsync(partitionKey)).ReturnsAsync(true);

            var result = await _service.DeleteByPartitionKeyAsync(partitionKey);

            Assert.True(result);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Once);
        }

        [Fact]
        public async Task DeleteByPartitionKeyAsync_ReturnsFalse_AndDoesNotSaveChanges_WhenNotFound()
        {
            var partitionKey = "WH001-PWTL999-ffffffffffffffffffffffffffffffff";
            _mockRepo.Setup(r => r.DeleteByPartitionKeyAsync(partitionKey)).ReturnsAsync(false);

            var result = await _service.DeleteByPartitionKeyAsync(partitionKey);

            Assert.False(result);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Never);
        }
    }
}
