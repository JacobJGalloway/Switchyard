using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Data;
using WarehouseLogistics_Claude.Data.Interfaces;
using WarehouseLogistics_Claude.Services;
using WarehouseLogistics_Claude.Services.Interfaces;
using WarehouseLogistics_Claude.Models;

var builder = WebApplication.CreateBuilder(args);

var httpsPort = builder.Configuration.GetValue<int>("Ports:Https", 7001);
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
    options.AddPolicy("ReadBOL",
        policy => policy.RequireClaim("permissions", "read:bol"));
    options.AddPolicy("CreateBOL",
        policy => policy.RequireClaim("permissions", "create:bol"));
    options.AddPolicy("ModifyBOL",
        policy => policy.RequireClaim("permissions", "modify:bol"));
    options.AddPolicy("ManageUsers",
        policy => policy.RequireClaim("permissions", "manage:users"));
});

var dbPath = Path.Combine(builder.Environment.ContentRootPath, "..", "Sqlite 3 Implementation", "WarehouseData.db3");
builder.Services.AddDbContext<LogisticsContext>(options => options.UseSqlite($"Data Source={dbPath}"));

builder.Services.AddScoped<IUnitOfWork, UnitOfWork>();
builder.Services.AddScoped<IBillOfLadingService, BillOfLadingService>();
builder.Services.AddScoped<IUserManagementService, UserManagementService>();
builder.Services.AddControllers();

var app = builder.Build();

app.UseRouting();
app.UseCors();
app.UseAuthentication();
app.UseAuthorization();
app.MapControllers();

app.Run();
