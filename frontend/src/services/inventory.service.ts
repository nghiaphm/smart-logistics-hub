import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Inventory = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"]
export type PaginatedInventory = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.PaginatedResponse"]
export type CreateInventoryRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.CreateInventoryRequest"]
export type UpdateInventoryRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.UpdateInventoryRequest"]

export function listInventory() {
  return apiClient<PaginatedInventory>("/inventory?limit=100")
}

export function createInventory(payload: CreateInventoryRequest) {
  return apiClient<Inventory>("/inventory", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateInventory(id: number, payload: UpdateInventoryRequest) {
  return apiClient<Inventory>(`/inventory/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteInventory(id: number) {
  return apiClient<void>(`/inventory/${id}`, { method: "DELETE" })
}
