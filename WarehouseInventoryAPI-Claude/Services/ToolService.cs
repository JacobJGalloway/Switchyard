using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Models;
using WarehouseInventory_Claude.Services.Interfaces;

namespace WarehouseInventory_Claude.Services
{
    public class ToolService(IUnitOfWork unitOfWork) : IToolService
    {
        private readonly IUnitOfWork _unitOfWork = unitOfWork;

        public async Task<IEnumerable<Tool>> GetAllAsync()
        {
            return await _unitOfWork.Tools.GetAllAsync();
        }

        public async Task<List<Tool>> GetBySKUIdAsync(string skuId)
        {
            return await _unitOfWork.GetToolBySKUIdAsync(skuId);
        }

        public async Task<List<Tool>> GetByLocationAsync(string locationId)
        {
            return await _unitOfWork.Tools.GetByLocationAsync(locationId);
        }

        public async Task<List<Tool>> GetByLocationAndSKUAsync(string locationId, string skuId)
        {
            return await _unitOfWork.Tools.GetByLocationAndSKUAsync(locationId, skuId);
        }

        public async Task<Tool> AddAsync(Tool item)
        {
            var created = await _unitOfWork.Tools.AddAsync(item);
            await _unitOfWork.SaveChangesAsync();
            return created;
        }

        public async Task UpdateBySKUIdAsync(string skuId, Tool item)
        {
            await _unitOfWork.Tools.UpdateBySKUIdAsync(skuId, item);
            await _unitOfWork.SaveChangesAsync();
        }

        public async Task PatchAsync(string partitionKey, bool? projected, DateTime? unloadedDate)
        {
            await _unitOfWork.Tools.PatchAsync(partitionKey, projected, unloadedDate);
            await _unitOfWork.SaveChangesAsync();
        }

        public async Task<bool> DeleteByPartitionKeyAsync(string partitionKey)
        {
            var result = await _unitOfWork.Tools.DeleteByPartitionKeyAsync(partitionKey);
            if (result) await _unitOfWork.SaveChangesAsync();
            return result;
        }
    }
}
