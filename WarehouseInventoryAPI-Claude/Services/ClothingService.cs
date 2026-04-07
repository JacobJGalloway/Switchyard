using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Services
{
    public class ClothingService(IUnitOfWork unitOfWork) : IClothingService
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        public async Task<IEnumerable<Clothing>> GetAllAsync()
        {
            return await _unitOfWork.Clothing.GetAllAsync();
        }

        public async Task<List<Clothing>> GetBySKUIdAsync(string skuId)
        {
            return await _unitOfWork.GetClothingBySKUIdAsync(skuId);
        }

        public async Task<List<Clothing>> GetByLocationAsync(string locationId)
        {
            return await _unitOfWork.Clothing.GetByLocationAsync(locationId);
        }

        public async Task<List<Clothing>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await _unitOfWork.Clothing.GetByLocationAndSKUAsync(locationId, skuId);
        }

        public async Task<Clothing> AddAsync(Clothing item)
        {
            var created = await _unitOfWork.Clothing.AddAsync(item);
            await _unitOfWork.SaveChangesAsync();
            return created;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Clothing item)
        {
            await _unitOfWork.Clothing.UpdateBySKUIdAsync(skuId, item);
            await _unitOfWork.SaveChangesAsync();
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            await _unitOfWork.Clothing.PatchAsync(partitionKey, projected, unloadedDate);
            await _unitOfWork.SaveChangesAsync();
        }

        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var result = await _unitOfWork.Clothing.DeleteByPartitionKeyAsync(partitionKey);
            if (result) await _unitOfWork.SaveChangesAsync();
            return result;
        }
    }
}
