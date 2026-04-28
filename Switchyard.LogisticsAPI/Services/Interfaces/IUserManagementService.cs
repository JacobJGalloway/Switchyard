using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Services.Interfaces
{
    public interface IUserManagementService
    {
        Task<IEnumerable<AppUser>> GetAllUsersAsync();
        Task<AppUser> CreateUserAsync(CreateUserRequest request);
        Task DeactivateUserAsync(string userId);
    }
}
