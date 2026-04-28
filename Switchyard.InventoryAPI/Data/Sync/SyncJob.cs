namespace Switchyard.InventoryAPI.Data.Sync
{
    public record SyncJob(IReadOnlySet<Type> ChangedTypes);
}
