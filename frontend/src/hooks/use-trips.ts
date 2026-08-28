"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type Trip = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"]
type PaginatedTrips = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.PaginatedResponse"]
type CreateTripRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.CreateTripRequest"]
type UpdateTripRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.UpdateTripRequest"]

export function useTrips(limit = 100) {
  return useQuery({
    queryKey: ["trips", { limit }],
    queryFn: () => apiClient<PaginatedTrips>(`/trips?limit=${limit}`),
  })
}

export function useCreateTrip(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateTripRequest) =>
      apiClient<Trip>("/trips", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["trips"] })
      onSuccess?.()
    },
  })
}

export function useUpdateTrip(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateTripRequest }) =>
      apiClient<Trip>(`/trips/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["trips"] })
      onSuccess?.()
    },
  })
}

export function useDeleteTrip(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiClient<void>(`/trips/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["trips"] })
      onSuccess?.()
    },
  })
}
