using System.Threading.Channels;
using Microsoft.EntityFrameworkCore;
using Switchyard.LogisticsAPI.Models;

namespace Switchyard.LogisticsAPI.Data.Sync
{
    public class LogisticsSyncWorker(
        ChannelReader<SyncJob> reader,
        IServiceScopeFactory scopeFactory,
        ILogger<LogisticsSyncWorker> logger) : BackgroundService
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
                    logger.LogError(ex, "Logistics sync failed for types: {Types}",
                        string.Join(", ", job.ChangedTypes.Select(t => t.Name)));
                }
            }
        }

        private async Task SyncAsync(SyncJob job, CancellationToken ct)
        {
            using var scope = scopeFactory.CreateScope();
            var writeCtx = scope.ServiceProvider.GetRequiredService<LogisticsContext>();
            var readCtx = scope.ServiceProvider.GetRequiredService<LogisticsReadContext>();

            foreach (var type in job.ChangedTypes)
            {
                if (type == typeof(BillOfLading))
                    await ResyncTableAsync(writeCtx.BillsOfLading, readCtx.BillsOfLading, readCtx, ct);
                else if (type == typeof(LineEntry))
                    await ResyncTableAsync(writeCtx.LineEntries, readCtx.LineEntries, readCtx, ct);
                else if (type == typeof(Warehouse))
                    await ResyncTableAsync(writeCtx.Warehouses, readCtx.Warehouses, readCtx, ct);
                else if (type == typeof(Store))
                    await ResyncTableAsync(writeCtx.Stores, readCtx.Stores, readCtx, ct);
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
