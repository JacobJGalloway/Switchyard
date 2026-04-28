using System.Linq;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data

{
    public class InventoryContext(DbContextOptions<InventoryContext> options) : DbContext(options)
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

        public async Task<List<Clothing>> GetClothingBySKUIdsync(string skuId)
        {
            var response = await Clothing.Where(c => c.SKUMarker == skuId).ToListAsync();

            if (response.Count == 0) return new List<Clothing>();

            return response;
        }

        public async Task<List<PPE>> GetPPEBySKUIdsync(string skuId)
        {
            var response = await PPE.Where(p => p.SKUMarker == skuId).ToListAsync();

            if (response.Count == 0) return new List<PPE>();

            return response;
        }

        public async Task<List<Tool>> GetToolBySKUIdsync(string skuId)
        {
            var response = await Tool.Where(t => t.SKUMarker == skuId).ToListAsync();

            if (response.Count == 0) return new List<Tool>();

            return response;
        }
    }
}