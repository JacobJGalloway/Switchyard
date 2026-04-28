using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Repositories
{
    public class LineEntryRepository(LogisticsContext writeContext, LogisticsReadContext readContext) : ILineEntryRepository
    {
        public async Task<IEnumerable<LineEntry>> GetAllAsync()
        {
            return await readContext.LineEntries.AsNoTracking().ToListAsync();
        }

        public async Task<List<LineEntry>> GetLineEntriesByTransactionIdAsync(string transactionId)
        {
            return await readContext.LineEntries
                .Where(le => le.TransactionId == transactionId)
                .AsNoTracking()
                .ToListAsync();
        }

        public async Task<LineEntry> AddAsync(LineEntry lineEntry)
        {
            lineEntry.PartitionKey = $"{lineEntry.TransactionId}-{Guid.NewGuid():N}";
            writeContext.LineEntries.Add(lineEntry);
            return lineEntry;
        }

        public async Task UpdateLineEntryAsync(LineEntry lineEntry)
        {
            writeContext.LineEntries.Update(lineEntry);
        }

        public async Task<bool> DeleteByTransactionIdAsync(string transactionId)
        {
            var entries = await writeContext.LineEntries
                .Where(le => le.TransactionId == transactionId)
                .ToListAsync();
            if (entries.Count == 0) return false;

            writeContext.LineEntries.RemoveRange(entries);
            return true;
        }

        public async Task<bool> DeleteByLocationAsync(string transactionId, string locationId)
        {
            var entries = await writeContext.LineEntries
                .Where(le => le.TransactionId == transactionId && le.LocationId == locationId)
                .ToListAsync();
            if (entries.Count == 0) return false;

            writeContext.LineEntries.RemoveRange(entries);
            return true;
        }
    }
}
