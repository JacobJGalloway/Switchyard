# Switchyard — Logistics API

This is the API to manipulate the inventory of Items in a Warehouse or Store. It will process either individual requests or a BillofLading from a shipment to update the Warehouse or Store Inventory by increasing and removing quantities as available. The UnloadedDate will show if the Inventory will be updated in the future or is part of the current Inventory. You will also be able to start checking a Warehouse or Store Item by SKU and UnloadedDate to see current Inventory. 

BillOfLading processing will handle one or more Warehouses or Stores in total. Each stop can only be for one Warehouse or Store at a time. There will be a lookup table for logistics transfer times between warehouses. Warehouse-to-Store and Store-to-Store with the same BaseWarehouseId transfers are considered same day.
