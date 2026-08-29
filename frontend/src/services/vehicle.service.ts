import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Vehicle = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.VehicleResponse"]
export type PaginatedVehicles = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.PaginatedResponse"]
export type CreateVehicleRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.CreateVehicleRequest"]
export type UpdateVehicleRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.UpdateVehicleRequest"]

export function listVehicles(limit: number) {
  return apiClient<PaginatedVehicles>(`/vehicles?limit=${limit}`)
}

export function createVehicle(payload: CreateVehicleRequest) {
  return apiClient<Vehicle>("/vehicles", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateVehicle(id: number, payload: UpdateVehicleRequest) {
  return apiClient<Vehicle>(`/vehicles/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteVehicle(id: number) {
  return apiClient<void>(`/vehicles/${id}`, { method: "DELETE" })
}
