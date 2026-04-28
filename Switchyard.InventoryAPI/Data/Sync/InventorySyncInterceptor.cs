using System.Threading.Channels;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Diagnostics;

namespace Switchyard.InventoryAPI.Data.Sync
{
    public class InventorySyncInterceptor(ChannelWriter<SyncJob> writer) : SaveChangesInterceptor
    {
        private static readonly AsyncLocal<HashSet<Type>?> _pendingTypes = new();

        public override ValueTask<InterceptionResult<int>> SavingChangesAsync(
            DbContextEventData eventData,
            InterceptionResult<int> result,
            CancellationToken cancellationToken = default)
        {
            _pendingTypes.Value = eventData.Context?.ChangeTracker.Entries()
                .Where(e => e.State is EntityState.Added or EntityState.Modified or EntityState.Deleted)
                .Select(e => e.Entity.GetType())
                .ToHashSet();
            return ValueTask.FromResult(result);
        }

        public override async ValueTask<int> SavedChangesAsync(
            SaveChangesCompletedEventData eventData,
            int result,
            CancellationToken cancellationToken = default)
        {
            var types = _pendingTypes.Value;
            _pendingTypes.Value = null;

            if (types?.Count > 0)
                await writer.WriteAsync(new SyncJob(types), cancellationToken);

            return result;
        }
    }
}
