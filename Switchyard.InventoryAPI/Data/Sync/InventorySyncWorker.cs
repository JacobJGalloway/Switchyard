using System.Threading.Channels;
using Microsoft.EntityFrameworkCore;
using Switchyard.InventoryAPI.Models;

namespace Switchyard.InventoryAPI.Data.Sync
{
    public class InventorySyncWorker(
        ChannelReader<SyncJob> reader,
        IServiceScopeFactory scopeFactory,
        ILogger<InventorySyncWorker> logger) : BackgroundService
    {
        protected override async Task ExecuteAsync(CancellationToken stoppingToken)
        {
            await foreach (var job in reader.ReadAllAsync(stoppingToken))
            {
                try
                {
                    await SyncAsync(job, stoppingToken);
                }
                catch (Exception ex)
                {
                    logger.LogError(ex, "Inventory sync failed for types: {Types}",
                        string.Join(", ", job.ChangedTypes.Select(t => t.Name)));
                }
            }
        }

        private async Task SyncAsync(SyncJob job, CancellationToken ct)
        {
            using var scope = scopeFactory.CreateScope();
            var writeCtx = scope.ServiceProvider.GetRequiredService<InventoryContext>();
            var readCtx = scope.ServiceProvider.GetRequiredService<InventoryReadContext>();

            foreach (var type in job.ChangedTypes)
            {
                if (type == typeof(Clothing))
                    await ResyncTableAsync(writeCtx.Clothing, readCtx.Clothing, readCtx, ct);
                else if (type == typeof(PPE))
                    await ResyncTableAsync(writeCtx.PPE, readCtx.PPE, readCtx, ct);
                else if (type == typeof(Tool))
                    await ResyncTableAsync(writeCtx.Tool, readCtx.Tool, readCtx, ct);
            }
        }

        private static async Task ResyncTableAsync<T>(
            DbSet<T> writeSet,
            DbSet<T> readSet,
            DbContext readCtx,
            CancellationToken ct) where T : class
        {
            var items = await writeSet.AsNoTracking().ToListAsync(ct);
            await readSet.ExecuteDeleteAsync(ct);
            readSet.AddRange(items);
            await readCtx.SaveChangesAsync(ct);
        }
    }
}
