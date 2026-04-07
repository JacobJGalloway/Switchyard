using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Services
{
    public class PPEService(IUnitOfWork unitOfWork) : IPPEService
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        public async Task<IEnumerable<PPE>> GetAllAsync()
        {
            return await _unitOfWork.PPE.GetAllAsync();
        }

        public async Task<List<PPE>> GetBySKUIdAsync(string skuId)
        {
            return await _unitOfWork.GetPPEBySKUIdAsync(skuId);
        }

        public async Task<List<PPE>> GetByLocationAsync(string locationId)
        {
            return await _unitOfWork.PPE.GetByLocationAsync(locationId);
        }

        public async Task<List<PPE>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await _unitOfWork.PPE.GetByLocationAndSKUAsync(locationId, skuId);
        }

        public async Task<PPE> AddAsync(PPE item)
        {
            var created = await _unitOfWork.PPE.AddAsync(item);
            await _unitOfWork.SaveChangesAsync();
            return created;
        }

        public async Task UpdateBySKUIdAsync(string skuId, PPE item)
        {
            await _unitOfWork.PPE.UpdateBySKUIdAsync(skuId, item);
            await _unitOfWork.SaveChangesAsync();
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            await _unitOfWork.PPE.PatchAsync(partitionKey, projected, unloadedDate);
            await _unitOfWork.SaveChangesAsync();
        }

        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var result = await _unitOfWork.PPE.DeleteByPartitionKeyAsync(partitionKey);
            if (result) await _unitOfWork.SaveChangesAsync();
            return result;
        }
    }
}
