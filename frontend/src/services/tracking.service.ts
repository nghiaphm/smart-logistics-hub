import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type TrackingEvent = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"]
export type PaginatedTracking = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.PaginatedResponse"]
export type CreateTrackingEventRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.CreateTrackingEventRequest"]
export type UpdateTrackingEventRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.UpdateTrackingEventRequest"]

export function listTracking(limit: number) {
  return apiClient<PaginatedTracking>(`/tracking-logs?limit=${limit}`)
}

export function createTrackingEvent(payload: CreateTrackingEventRequest) {
  return apiClient<TrackingEvent>("/tracking-logs", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateTrackingEvent(id: number, payload: UpdateTrackingEventRequest) {
  return apiClient<TrackingEvent>(`/tracking-logs/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteTrackingEvent(id: number) {
  return apiClient<void>(`/tracking-logs/${id}`, { method: "DELETE" })
}
