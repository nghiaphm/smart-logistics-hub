import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Trip = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"]
export type PaginatedTrips = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.PaginatedResponse"]
export type CreateTripRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.CreateTripRequest"]
export type UpdateTripRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.UpdateTripRequest"]

export function listTrips(limit: number) {
  return apiClient<PaginatedTrips>(`/trips?limit=${limit}`)
}

export function createTrip(payload: CreateTripRequest) {
  return apiClient<Trip>("/trips", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateTrip(id: number, payload: UpdateTripRequest) {
  return apiClient<Trip>(`/trips/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteTrip(id: number) {
  return apiClient<void>(`/trips/${id}`, { method: "DELETE" })
}
