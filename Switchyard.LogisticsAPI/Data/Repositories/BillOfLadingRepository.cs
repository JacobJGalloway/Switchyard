using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Data.Interfaces;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Repositories
{
    public class BillOfLadingRepository(LogisticsContext writeContext, LogisticsReadContext readContext) : IBillOfLadingRepository
    {
        public async Task<IEnumerable<BillOfLading>> GetAllAsync()
        {
            return await readContext.BillsOfLading.AsNoTracking().ToListAsync();
        }

        public async Task<BillOfLading?> GetByTransactionIdAsync(string transactionId)
        {
            return await readContext.BillsOfLading
                .AsNoTracking()
                .FirstOrDefaultAsync(b => b.TransactionId == transactionId);
        }

        public async Task<BillOfLading> AddAsync(BillOfLading billOfLading)
        {
            writeContext.BillsOfLading.Add(billOfLading);
            return billOfLading;
        }

        public async Task UpdateAsync(BillOfLading billOfLading)
        {
            var target = await writeContext.BillsOfLading
                .FirstOrDefaultAsync(b => b.TransactionId == billOfLading.TransactionId);
            if (target is null) return;

            target.Status = billOfLading.Status;
            target.CustomerFirstName = billOfLading.CustomerFirstName;
            target.CustomerLastName = billOfLading.CustomerLastName;
            target.City = billOfLading.City;
            target.State = billOfLading.State;
        }

        public async Task<bool> DeleteByTransactionIdAsync(string transactionId)
        {
            var target = await writeContext.BillsOfLading
                .FirstOrDefaultAsync(b => b.TransactionId == transactionId);
            if (target is null) return false;

            writeContext.BillsOfLading.Remove(target);
            return true;
        }
    }
}
