using Xunit;
using Moq;
using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Models;
using WarehouseLogistics_Claude.Services;

namespace WarehouseLogistics_Claude.Tests.Services
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
    }
}
