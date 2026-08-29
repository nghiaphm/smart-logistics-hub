import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Driver = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"]
export type PaginatedDrivers = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.PaginatedResponse"]
export type CreateDriverRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.CreateDriverRequest"]
export type UpdateDriverRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.UpdateDriverRequest"]

export function listDrivers(limit: number) {
  return apiClient<PaginatedDrivers>(`/drivers?limit=${limit}`)
}

export function createDriver(payload: CreateDriverRequest) {
  return apiClient<Driver>("/drivers", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateDriver(id: number, payload: UpdateDriverRequest) {
  return apiClient<Driver>(`/drivers/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteDriver(id: number) {
  return apiClient<void>(`/drivers/${id}`, { method: "DELETE" })
}
