using Microsoft.AspNetCore.Mvc;
using Moq;
using Xunit;
using WarehouseInventory_Claude.Controllers;
using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;

namespace WarehouseInventory_Claude.Tests.Controllers;

public class PPEControllerTests
{
    private readonly Mock<IPPERepository> _mockRepo;
    private readonly PPEController _controller;

    public PPEControllerTests()
    {
        _mockRepo = new Mock<IPPERepository>();
        _controller = new PPEController(_mockRepo.Object);
    }

    [Fact]
    public async Task GetAll_ReturnsOkWithAllItems()
    {
        var items = new List<PPE>
        {
            new() { PartitionKey = "123-SKU001-a1b2c3d4e5f6478a9b0cdef123456789", SKUMarker = "SKU001", Type = "Gloves", Size = "L" },
            new() { PartitionKey = "123-SKU002-b2c3d4e5f6a7489b0c1def234567890a", SKUMarker = "SKU002", Type = "Helmet", Size = "M" }
        };
        _mockRepo.Setup(r => r.GetAllAsync()).ReturnsAsync(items);

        var result = await _controller.GetAll();

        var ok = Assert.IsType<OkObjectResult>(result.Result);
        var returned = Assert.IsAssignableFrom<IEnumerable<PPE>>(ok.Value);
        Assert.Equal(2, returned.Count());
    }

    [Fact]
    public async Task GetBySKUId_ReturnsItem_WhenFound()
    {
        var item = new PPE { PartitionKey = "123-SKU001-a1b2c3d4e5f6478a9b0cdef123456789", SKUMarker = "SKU001", Type = "Gloves" };
        _mockRepo.Setup(r => r.GetBySKUIdAsync("SKU001")).ReturnsAsync(item);

        var result = await _controller.GetBySKUId("SKU001");

        var returned = Assert.IsType<PPE>(result.Value);
        Assert.Equal("SKU001", returned.SKUMarker);
    }

    [Fact]
    public async Task GetBySKUId_ReturnsNotFound_WhenMissing()
    {
        _mockRepo.Setup(r => r.GetBySKUIdAsync("MISSING")).ReturnsAsync((PPE?)null);

        var result = await _controller.GetBySKUId("MISSING");

        Assert.IsType<NotFoundResult>(result.Result);
    }

    [Fact]
    public async Task Create_ReturnsCreatedAtAction_WithCreatedItem()
    {
        var item = new PPE { PartitionKey = "123-SKU001-a1b2c3d4e5f6478a9b0cdef123456789", SKUMarker = "SKU001", Type = "Gloves" };
        _mockRepo.Setup(r => r.AddAsync(item)).ReturnsAsync(item);

        var result = await _controller.Create(item);

        var created = Assert.IsType<CreatedAtActionResult>(result.Result);
        Assert.Equal(nameof(PPEController.Create), created.ActionName);
        Assert.Equal(item, created.Value);
    }

    [Fact]
    public async Task UpdateBySKUId_ReturnsNoContent_WhenSKUMatches()
    {
        var item = new PPE { PartitionKey = "123-SKU001-a1b2c3d4e5f6478a9b0cdef123456789", SKUMarker = "SKU001" };
        _mockRepo.Setup(r => r.UpdateBySKUIdAsync("SKU001", item)).Returns(Task.CompletedTask);

        var result = await _controller.UpdateBySKUId("SKU001", item);

        Assert.IsType<NoContentResult>(result);
    }

    [Fact]
    public async Task UpdateBySKUId_ReturnsBadRequest_WhenSKUMismatch()
    {
        var item = new PPE { PartitionKey = "123-SKU999-c3d4e5f6a7b8490c1d2ef345678901b", SKUMarker = "SKU999" };

        var result = await _controller.UpdateBySKUId("SKU001", item);

        Assert.IsType<BadRequestResult>(result);
    }

    [Fact]
    public async Task DeleteBySKUId_ReturnsNoContent_WhenDeleted()
    {
        _mockRepo.Setup(r => r.DeleteBySKUIdAsync("SKU001")).ReturnsAsync(true);

        var result = await _controller.DeleteBySKUId("SKU001");

        Assert.IsType<NoContentResult>(result);
    }

    [Fact]
    public async Task DeleteBySKUId_ReturnsNotFound_WhenMissing()
    {
        _mockRepo.Setup(r => r.DeleteBySKUIdAsync("MISSING")).ReturnsAsync(false);

        var result = await _controller.DeleteBySKUId("MISSING");

        Assert.IsType<NotFoundResult>(result);
    }
}
