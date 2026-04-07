using Microsoft.EntityFrameworkCore;
using WarehouseLogistics_Claude.Models;

namespace WarehouseLogistics_Claude.Data
{
    public class LogisticsContext(DbContextOptions<LogisticsContext> options) : DbContext(options)
    {
        public DbSet<BillOfLading> BillsOfLading => Set<BillOfLading>();
        public DbSet<LineEntry> LineEntries => Set<LineEntry>();
        public DbSet<Warehouse> Warehouses => Set<Warehouse>();
        public DbSet<Store> Stores => Set<Store>();

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<BillOfLading>().ToTable("BillsOfLading").HasKey(b => b.PartitionKey);
            modelBuilder.Entity<BillOfLading>().HasAlternateKey(b => b.TransactionId);
            modelBuilder.Entity<BillOfLading>().Ignore(b => b.LineEntries);

            modelBuilder.Entity<LineEntry>().ToTable("LineEntries").HasKey(le => le.PartitionKey);

            modelBuilder.Entity<Warehouse>().ToTable("Warehouses").HasKey(w => w.WarehouseId);
            modelBuilder.Entity<Warehouse>().Ignore(w => w.PartitionKey);

            modelBuilder.Entity<Store>().ToTable("Stores").HasKey(s => s.PartitionKey);
        }
    }
}
