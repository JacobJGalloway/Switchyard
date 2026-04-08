namespace WarehouseLogistics_Claude.Models
{
    public class AppUser
    {
        public string UserId { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public string Name { get; set; } = string.Empty;
        public string Role { get; set; } = string.Empty;
        public string WarehouseId { get; set; } = string.Empty;
        public string StoreId { get; set; } = string.Empty;
        public bool Blocked { get; set; }
    }

    public class CreateUserRequest
    {
        public required string Email { get; set; }
        public required string Password { get; set; }
        public required string FirstName { get; set; }
        public required string LastName { get; set; }
        public required string Role { get; set; }
        public string WarehouseId { get; set; } = string.Empty;
        public string StoreId { get; set; } = string.Empty;
    }
}
