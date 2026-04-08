using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using WarehouseInventory_Claude.Data;
using WarehouseInventory_Claude.Data.Interfaces;
using WarehouseInventory_Claude.Services;
using WarehouseInventory_Claude.Services.Interfaces;

var builder = WebApplication.CreateBuilder(args);

var httpsPort = builder.Configuration.GetValue<int>("Ports:Https", 7000);
builder.WebHost.UseUrls($"https://localhost:{httpsPort}");

builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
        policy.WithOrigins("http://localhost:5173")
              .AllowAnyHeader()
              .AllowAnyMethod());
});

builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme)
    .AddJwtBearer(options =>
    {
        options.Authority = builder.Configuration["Auth0:Authority"];
        options.Audience = builder.Configuration["Auth0:Audience"];
    });

builder.Services.AddAuthorization(options =>
{
    options.AddPolicy("ReadInventory",
        policy => policy.RequireClaim("permissions", "read:inventory"));
});

var dbPath = Path.Combine(builder.Environment.ContentRootPath, "..", "Sqlite 3 Implementation", "WarehouseData.db3");
builder.Services.AddDbContext<InventoryContext>(options => options.UseSqlite($"Data Source={dbPath}"));
builder.Services.AddScoped<IUnitOfWork, UnitOfWork>();
builder.Services.AddScoped<IClothingService, ClothingService>();
builder.Services.AddScoped<IPPEService, PPEService>();
builder.Services.AddScoped<IToolService, ToolService>();
builder.Services.AddControllers();

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    scope.ServiceProvider.GetRequiredService<InventoryContext>().Database.EnsureCreated();
}

app.UseRouting();
app.UseCors();
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();

app.Run();
