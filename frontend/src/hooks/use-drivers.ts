"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listDrivers,
  createDriver,
  updateDriver,
  deleteDriver,
  type CreateDriverRequest,
  type UpdateDriverRequest,
} from "@/services/driver.service"

export function useDrivers(limit = 100) {
  return useQuery({
    queryKey: ["drivers", { limit }],
    queryFn: () => listDrivers(limit),
  })
}

export function useCreateDriver(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateDriverRequest) => createDriver(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["drivers"] })
      onSuccess?.()
    },
  })
}

export function useUpdateDriver(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateDriverRequest }) => updateDriver(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["drivers"] })
      onSuccess?.()
    },
  })
}

export function useDeleteDriver(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteDriver(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["drivers"] })
      onSuccess?.()
    },
  })
}
