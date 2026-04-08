using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Repositories
{
    public class LineEntryRepository(LogisticsContext context) : ILineEntryRepository
    {
        private readonly LogisticsContext _context = context;

        public async Task<IEnumerable<LineEntry>> GetAllAsync()
        {
            return await _context.LineEntries.ToListAsync();
        }

        public async Task<List<LineEntry>> GetLineEntriesByTransactionIdAsync(string transactionId)
        {
            return await _context.LineEntries
                .Where(le => le.TransactionId == transactionId)
                .ToListAsync();
        }

        public async Task<LineEntry> AddAsync(LineEntry lineEntry)
        {
            lineEntry.PartitionKey = $"{lineEntry.TransactionId}-{Guid.NewGuid():N}";
            _context.LineEntries.Add(lineEntry);
            return lineEntry;
        }

        public async Task UpdateLineEntryAsync(LineEntry lineEntry)
        {
            _context.LineEntries.Update(lineEntry);
        }

        public async Task<bool> DeleteByTransactionIdAsync(string transactionId)
        {
            var entries = await _context.LineEntries
                .Where(le => le.TransactionId == transactionId)
                .ToListAsync();
            if (entries.Count == 0) return false;

            _context.LineEntries.RemoveRange(entries);
            return true;
        }

        public async Task<bool> DeleteByLocationAsync(string transactionId, string locationId)
        {
            var entries = await _context.LineEntries
                .Where(le => le.TransactionId == transactionId && le.LocationId == locationId)
                .ToListAsync();
            if (entries.Count == 0) return false;

            _context.LineEntries.RemoveRange(entries);
            return true;
        }
    }
}
