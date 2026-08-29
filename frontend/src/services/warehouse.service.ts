import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Warehouse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"]
export type PaginatedWarehouses = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"]

export function listWarehouses(limit = 100) {
  return apiClient<PaginatedWarehouses>(`/warehouses?limit=${limit}`)
}
