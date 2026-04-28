namespace Switchyard.LogisticsAPI.Data.Sync
{
    public record SyncJob(IReadOnlySet<Type> ChangedTypes);
}
