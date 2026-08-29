"use client"

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

import {
  listInventory,
  createInventory,
  updateInventory,
  deleteInventory,
  type CreateInventoryRequest,
  type UpdateInventoryRequest,
} from "@/services/inventory.service"
import { listWarehouses } from "@/services/warehouse.service"
import { listProducts } from "@/services/product.service"

export function useInventory() {
  return useQuery({
    queryKey: ["inventory", { limit: 100 }],
    queryFn: () => listInventory(),
  })
}

export function useInventoryFormOptions(enabled: boolean) {
  const warehouses = useQuery({
    queryKey: ["warehouses", { limit: 100 }],
    queryFn: () => listWarehouses(100),
    enabled,
  })
  const products = useQuery({
    queryKey: ["products", { limit: 100 }],
    queryFn: () => listProducts(100),
    enabled,
  })

  return { warehouses, products }
}

export function useCreateInventory(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateInventoryRequest) => createInventory(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inventory"] })
      onSuccess?.()
    },
  })
}

export function useUpdateInventory(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateInventoryRequest }) => updateInventory(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inventory"] })
      onSuccess?.()
    },
  })
}

export function useDeleteInventory(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteInventory(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inventory"] })
      onSuccess?.()
    },
  })
}
