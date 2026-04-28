using Xunit;
using Moq;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Services;

namespace Switchyard.LogisticsAPI.Tests.Services
{
    public class BillOfLadingServiceTests
    {
        private readonly Mock<IUnitOfWork> _mockUoW;
        private readonly Mock<IBillOfLadingRepository> _mockBOLRepo;
        private readonly Mock<ILineEntryRepository> _mockLineEntryRepo;
        private readonly BillOfLadingService _service;

        public BillOfLadingServiceTests()
        {
            _mockUoW = new Mock<IUnitOfWork>();
            _mockBOLRepo = new Mock<IBillOfLadingRepository>();
            _mockLineEntryRepo = new Mock<ILineEntryRepository>();
            _mockUoW.Setup(u => u.BillsOfLading).Returns(_mockBOLRepo.Object);
            _mockUoW.Setup(u => u.LineEntries).Returns(_mockLineEntryRepo.Object);
            _service = new BillOfLadingService(_mockUoW.Object);
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
        public async Task CreateAsync_ThrowsArgumentException_WhenStatusIsNotPending()
        {
            var bol = MakeValidBOL();
            bol.Status = "Submitted";

            await Assert.ThrowsAsync<ArgumentException>(() => _service.CreateAsync(bol));
        }

        [Fact]
        public async Task CreateAsync_ThrowsArgumentException_WhenCustomerInfoMissing()
        {
            var bol = MakeValidBOL();
            bol.CustomerFirstName = string.Empty;

            await Assert.ThrowsAsync<ArgumentException>(() => _service.CreateAsync(bol));
        }

        [Fact]
        public async Task CreateAsync_ThrowsArgumentException_WhenLineEntriesEmpty()
        {
            var bol = MakeValidBOL();
            bol.LineEntries = [];

            await Assert.ThrowsAsync<ArgumentException>(() => _service.CreateAsync(bol));
        }

        [Fact]
        public async Task CreateAsync_ThrowsArgumentException_WhenLocationIdNotConfigured()
        {
            Environment.SetEnvironmentVariable("LOCATION_ID", null);
            var bol = MakeValidBOL();

            await Assert.ThrowsAsync<ArgumentException>(() => _service.CreateAsync(bol));
        }

        [Fact]
        public async Task CreateAsync_ReturnsTransactionId_AndSavesChanges_WhenValid()
        {
            Environment.SetEnvironmentVariable("LOCATION_ID", "WH001");

            var bol = MakeValidBOL();
            _mockBOLRepo.Setup(r => r.AddAsync(It.IsAny<BillOfLading>())).ReturnsAsync(bol);
            _mockLineEntryRepo.Setup(r => r.AddAsync(It.IsAny<LineEntry>())).ReturnsAsync((LineEntry e) => e);
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync(It.IsAny<string>()))
                .ReturnsAsync([]);
            _mockUoW.Setup(u => u.SaveChangesAsync()).ReturnsAsync(1);

            var transactionId = await _service.CreateAsync(bol);

            Assert.False(string.IsNullOrWhiteSpace(transactionId));
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.AtLeastOnce);

            Environment.SetEnvironmentVariable("LOCATION_ID", null);
        }

        [Fact]
        public async Task ProcessLocationStop_MarksMatchingEntries_AndSavesChanges()
        {
            var entries = new List<LineEntry>
            {
                new() { PartitionKey = "txn001-guid1", TransactionId = "txn001", LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 5 },
                new() { PartitionKey = "txn001-guid2", TransactionId = "txn001", LocationId = "ST0001", SKUMarker = "CLTH001", Quantity = -5 }
            };
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001")).ReturnsAsync(entries);
            _mockLineEntryRepo.Setup(r => r.UpdateLineEntryAsync(It.IsAny<LineEntry>())).Returns(Task.CompletedTask);
            _mockUoW.Setup(u => u.SaveChangesAsync()).ReturnsAsync(1);

            await _service.ProcessLocationStop("txn001", "WH001");

            _mockLineEntryRepo.Verify(r => r.UpdateLineEntryAsync(It.Is<LineEntry>(e => e.LocationId == "WH001")), Times.Once);
            _mockLineEntryRepo.Verify(r => r.UpdateLineEntryAsync(It.Is<LineEntry>(e => e.LocationId == "ST0001")), Times.Never);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Once);
        }

        [Fact]
        public async Task ProcessLocationStop_DoesNotSaveChanges_WhenNoMatchingEntries()
        {
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001"))
                .ReturnsAsync([]);

            await _service.ProcessLocationStop("txn001", "WH001");

            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Never);
        }

        [Fact]
        public async Task GetAllAsync_ReturnsBOLs()
        {
            var bols = new List<BillOfLading> { MakeValidBOL() };
            _mockBOLRepo.Setup(r => r.GetAllAsync()).ReturnsAsync(bols);

            var result = await _service.GetAllAsync();

            Assert.Single(result);
        }

        [Fact]
        public async Task GetByTransactionIdAsync_ReturnsNull_WhenNotFound()
        {
            _mockBOLRepo.Setup(r => r.GetByTransactionIdAsync("txn999")).ReturnsAsync((BillOfLading?)null);

            var result = await _service.GetByTransactionIdAsync("txn999");

            Assert.Null(result);
        }

        [Fact]
        public async Task GetByTransactionIdAsync_AttachesLineEntries_WhenFound()
        {
            var bol = MakeValidBOL();
            bol.TransactionId = "txn001";
            var entries = new List<LineEntry>
            {
                new() { TransactionId = "txn001", LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 5 }
            };
            _mockBOLRepo.Setup(r => r.GetByTransactionIdAsync("txn001")).ReturnsAsync(bol);
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001")).ReturnsAsync(entries);

            var result = await _service.GetByTransactionIdAsync("txn001");

            Assert.NotNull(result);
            Assert.Single(result!.LineEntries);
        }

        [Fact]
        public async Task GetLineEntriesByTransactionIdAsync_ReturnsEntries()
        {
            var entries = new List<LineEntry>
            {
                new() { TransactionId = "txn001", LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 5 }
            };
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001")).ReturnsAsync(entries);

            var result = await _service.GetLineEntriesByTransactionIdAsync("txn001");

            Assert.Single(result);
        }

        [Fact]
        public async Task ReplaceLocationStopAsync_ThrowsArgumentException_WhenNoEntriesForLocation()
        {
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001"))
                .ReturnsAsync([]);

            await Assert.ThrowsAsync<ArgumentException>(() =>
                _service.ReplaceLocationStopAsync("txn001", "ST9999", "ST0001"));
        }

        [Fact]
        public async Task ReplaceLocationStopAsync_ThrowsInvalidOperationException_WhenAlreadyProcessed()
        {
            var entries = new List<LineEntry>
            {
                new() { TransactionId = "txn001", LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 5, IsProcessed = true }
            };
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001")).ReturnsAsync(entries);

            await Assert.ThrowsAsync<InvalidOperationException>(() =>
                _service.ReplaceLocationStopAsync("txn001", "WH001", "ST0001"));
        }

        [Fact]
        public async Task ReplaceLocationStopAsync_AddsNewEntries_DeletesOld_AndSaves()
        {
            var entries = new List<LineEntry>
            {
                new() { TransactionId = "txn001", LocationId = "WH001", SKUMarker = "CLTH001", Quantity = 5, IsProcessed = false },
                new() { TransactionId = "txn001", LocationId = "WH001", SKUMarker = "SPPE001", Quantity = 2, IsProcessed = false },
            };
            _mockLineEntryRepo.Setup(r => r.GetLineEntriesByTransactionIdAsync("txn001")).ReturnsAsync(entries);
            _mockLineEntryRepo.Setup(r => r.AddAsync(It.IsAny<LineEntry>())).ReturnsAsync((LineEntry e) => e);
            _mockLineEntryRepo.Setup(r => r.DeleteByLocationAsync("txn001", "WH001")).ReturnsAsync(true);
            _mockUoW.Setup(u => u.SaveChangesAsync()).ReturnsAsync(1);

            await _service.ReplaceLocationStopAsync("txn001", "WH001", "ST0001");

            _mockLineEntryRepo.Verify(r => r.AddAsync(It.Is<LineEntry>(e => e.LocationId == "ST0001")), Times.Exactly(2));
            _mockLineEntryRepo.Verify(r => r.DeleteByLocationAsync("txn001", "WH001"), Times.Once);
            _mockUoW.Verify(u => u.SaveChangesAsync(), Times.Once);
        }
    }
}
