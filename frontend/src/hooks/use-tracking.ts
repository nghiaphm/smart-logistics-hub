"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type TrackingEvent = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"]
type PaginatedTracking = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.PaginatedResponse"]
type CreateTrackingEventRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.CreateTrackingEventRequest"]
type UpdateTrackingEventRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.UpdateTrackingEventRequest"]

export function useTracking(limit = 100) {
  return useQuery({
    queryKey: ["tracking-logs", { limit }],
    queryFn: () => apiClient<PaginatedTracking>(`/tracking-logs?limit=${limit}`),
  })
}

export function useCreateTrackingEvent(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateTrackingEventRequest) =>
      apiClient<TrackingEvent>("/tracking-logs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tracking-logs"] })
      onSuccess?.()
    },
  })
}

export function useUpdateTrackingEvent(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateTrackingEventRequest }) =>
      apiClient<TrackingEvent>(`/tracking-logs/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tracking-logs"] })
      onSuccess?.()
    },
  })
}

export function useDeleteTrackingEvent(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiClient<void>(`/tracking-logs/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tracking-logs"] })
      onSuccess?.()
    },
  })
}
