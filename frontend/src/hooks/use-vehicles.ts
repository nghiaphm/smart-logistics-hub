"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listVehicles,
  createVehicle,
  updateVehicle,
  deleteVehicle,
  type CreateVehicleRequest,
  type UpdateVehicleRequest,
} from "@/services/vehicle.service"

export function useVehicles(limit = 100) {
  return useQuery({
    queryKey: ["vehicles", { limit }],
    queryFn: () => listVehicles(limit),
  })
}

export function useCreateVehicle(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateVehicleRequest) => createVehicle(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["vehicles"] })
      onSuccess?.()
    },
  })
}

export function useUpdateVehicle(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateVehicleRequest }) => updateVehicle(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["vehicles"] })
      onSuccess?.()
    },
  })
}

export function useDeleteVehicle(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteVehicle(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["vehicles"] })
      onSuccess?.()
    },
  })
}
