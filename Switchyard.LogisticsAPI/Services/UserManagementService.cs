using Auth0.AuthenticationApi;
using Auth0.AuthenticationApi.Models;
using Auth0.ManagementApi;
using Auth0.ManagementApi.Models;
using Switchyard.LogisticsAPI.Models;
using Switchyard.LogisticsAPI.Services.Interfaces;

namespace Switchyard.LogisticsAPI.Services
{
    // No unit tests for this service — it is a thin wrapper around the Auth0 Management SDK.
    // All logic lives in ManagementApiClient/AuthenticationApiClient; mocking those would only
    // test the mock, not this code. UserController tests cover the controller-layer behavior
    // (BadRequest on ArgumentException, etc.). Integration testing this requires a live Auth0 tenant.
    public class UserManagementService : IUserManagementService
    {
        private readonly string _domain;
        private readonly string _clientId;
        private readonly string _clientSecret;

        private string? _cachedToken;
        private DateTime _tokenExpiry = DateTime.MinValue;

        private Dictionary<string, string>? _roleNameToId;

        public UserManagementService(IConfiguration configuration)
        {
            _domain       = configuration["Auth0:Authority"]!.TrimEnd('/');
            _clientId     = configuration["Auth0:ManagementClientId"]!;
            _clientSecret = configuration["Auth0:ManagementClientSecret"]!;
        }

        private async Task<ManagementApiClient> GetClientAsync()
        {
            if (_cachedToken is null || DateTime.UtcNow >= _tokenExpiry)
            {
                var authClient = new AuthenticationApiClient(new Uri(_domain));
                var tokenResponse = await authClient.GetTokenAsync(new ClientCredentialsTokenRequest
                {
                    ClientId     = _clientId,
                    ClientSecret = _clientSecret,
                    Audience     = $"{_domain}/api/v2/",
                });
                _cachedToken = tokenResponse.AccessToken;
                _tokenExpiry = DateTime.UtcNow.AddSeconds(tokenResponse.ExpiresIn - 60);
            }

            return new ManagementApiClient(_cachedToken, new Uri($"{_domain}/api/v2"));
        }

        private async Task<Dictionary<string, string>> GetRoleMapAsync(ManagementApiClient client)
        {
            if (_roleNameToId is not null) return _roleNameToId;

            var roles = await client.Roles.GetAllAsync(new GetRolesRequest());
            _roleNameToId = roles.ToDictionary(r => r.Name, r => r.Id);
            return _roleNameToId;
        }

        public async Task<IEnumerable<AppUser>> GetAllUsersAsync()
        {
            var client = await GetClientAsync();

            var users = await client.Users.GetAllAsync(new GetUsersRequest());

            var roleTasks = users.Select(u => client.Users.GetRolesAsync(u.UserId));
            var allRoles = await Task.WhenAll(roleTasks);

            return users.Select((u, i) => new AppUser
            {
                UserId      = u.UserId,
                Email       = u.Email ?? string.Empty,
                Name        = u.FullName ?? u.Email ?? string.Empty,
                Role        = allRoles[i].FirstOrDefault()?.Name ?? string.Empty,
                WarehouseId = u.UserMetadata?.assigned_warehouse_id ?? string.Empty,
                StoreId     = u.UserMetadata?.assigned_store_id ?? string.Empty,
                Blocked     = u.Blocked ?? false,
            });
        }

        public async Task<AppUser> CreateUserAsync(CreateUserRequest request)
        {
            var client = await GetClientAsync();
            var roleMap = await GetRoleMapAsync(client);

            if (!roleMap.TryGetValue(request.Role, out var roleId))
                throw new ArgumentException($"Role '{request.Role}' does not exist.");

            var created = await client.Users.CreateAsync(new UserCreateRequest
            {
                Email      = request.Email,
                Password   = request.Password,
                FirstName  = request.FirstName,
                LastName   = request.LastName,
                FullName   = $"{request.FirstName} {request.LastName}",
                Connection = "Username-Password-Authentication",
                UserMetadata = new
                {
                    assigned_warehouse_id = request.WarehouseId,
                    assigned_store_id     = request.StoreId,
                },
            });

            await client.Users.AssignRolesAsync(created.UserId, new AssignRolesRequest
            {
                Roles = [roleId]
            });

            return new AppUser
            {
                UserId      = created.UserId,
                Email       = created.Email ?? string.Empty,
                Name        = created.FullName ?? string.Empty,
                Role        = request.Role,
                WarehouseId = request.WarehouseId,
                StoreId     = request.StoreId,
                Blocked     = false,
            };
        }

        public async Task DeactivateUserAsync(string userId)
        {
            var client = await GetClientAsync();
            await client.Users.UpdateAsync(userId, new UserUpdateRequest { Blocked = true });
        }
    }
}
