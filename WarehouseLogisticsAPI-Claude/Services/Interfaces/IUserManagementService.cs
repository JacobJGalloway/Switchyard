using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Services.Interfaces
{
    public interface IUserManagementService
    {
        Task<IEnumerable<AppUser>> GetAllUsersAsync();
        Task<AppUser> CreateUserAsync(CreateUserRequest request);
        Task DeactivateUserAsync(string userId);
    }
}
