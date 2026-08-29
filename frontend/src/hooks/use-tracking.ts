"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listTracking,
  createTrackingEvent,
  updateTrackingEvent,
  deleteTrackingEvent,
  type CreateTrackingEventRequest,
  type UpdateTrackingEventRequest,
} from "@/services/tracking.service"

export function useTracking(limit = 100) {
  return useQuery({
    queryKey: ["tracking-logs", { limit }],
    queryFn: () => listTracking(limit),
  })
}

export function useCreateTrackingEvent(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateTrackingEventRequest) => createTrackingEvent(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tracking-logs"] })
      onSuccess?.()
    },
  })
}

export function useUpdateTrackingEvent(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateTrackingEventRequest }) => updateTrackingEvent(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tracking-logs"] })
      onSuccess?.()
    },
  })
}

export function useDeleteTrackingEvent(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteTrackingEvent(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tracking-logs"] })
      onSuccess?.()
    },
  })
}
