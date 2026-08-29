"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listTrips,
  createTrip,
  updateTrip,
  deleteTrip,
  type CreateTripRequest,
  type UpdateTripRequest,
} from "@/services/trip.service"

export function useTrips(limit = 100) {
  return useQuery({
    queryKey: ["trips", { limit }],
    queryFn: () => listTrips(limit),
  })
}

export function useCreateTrip(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateTripRequest) => createTrip(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["trips"] })
      onSuccess?.()
    },
  })
}

export function useUpdateTrip(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateTripRequest }) => updateTrip(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["trips"] })
      onSuccess?.()
    },
  })
}

export function useDeleteTrip(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteTrip(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["trips"] })
      onSuccess?.()
    },
  })
}
