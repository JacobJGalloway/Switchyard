using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data.Repositories
{
    public class BillOfLadingRepository(LogisticsContext context) : IBillOfLadingRepository
    {
        private readonly LogisticsContext _context = context;

        public async Task<IEnumerable<BillOfLading>> GetAllAsync()
        {
            return await _context.BillsOfLading.ToListAsync();
        }

        public async Task<BillOfLading?> GetByTransactionIdAsync(string transactionId)
        {
            return await _context.BillsOfLading
                .FirstOrDefaultAsync(b => b.TransactionId == transactionId);
        }

        public async Task<BillOfLading> AddAsync(BillOfLading billOfLading)
        {
            _context.BillsOfLading.Add(billOfLading);
            return billOfLading;
        }

        public async Task UpdateAsync(BillOfLading billOfLading)
        {
            var target = await _context.BillsOfLading
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
            var target = await _context.BillsOfLading
                .FirstOrDefaultAsync(b => b.TransactionId == transactionId);
            if (target is null) return false;

            _context.BillsOfLading.Remove(target);
            return true;
        }
    }
}
