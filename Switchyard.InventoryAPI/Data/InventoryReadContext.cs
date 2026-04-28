using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data
{
    public class InventoryReadContext(DbContextOptions<InventoryReadContext> options) : DbContext(options)
    {
        public DbSet<Clothing> Clothing => Set<Clothing>();
        public DbSet<PPE> PPE => Set<PPE>();
        public DbSet<Tool> Tool => Set<Tool>();

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<Clothing>().ToTable("Clothing").HasKey(c => c.PartitionKey);
            modelBuilder.Entity<PPE>().ToTable("PPE").HasKey(p => p.PartitionKey);
            modelBuilder.Entity<Tool>().ToTable("Tools").HasKey(t => t.PartitionKey);
        }
    }
}
