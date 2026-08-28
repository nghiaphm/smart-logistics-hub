"use client"

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type Inventory = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"]
type PaginatedInventory = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.PaginatedResponse"]
type CreateInventoryRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.CreateInventoryRequest"]
type UpdateInventoryRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.UpdateInventoryRequest"]

type PaginatedWarehouses = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"]
type PaginatedProducts = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse"]

export function useInventory() {
  return useQuery({
    queryKey: ["inventory", { limit: 100 }],
    queryFn: () => apiClient<PaginatedInventory>("/inventory?limit=100"),
  })
}

export function useInventoryFormOptions(enabled: boolean) {
  const warehouses = useQuery({
    queryKey: ["warehouses", { limit: 100 }],
    queryFn: () => apiClient<PaginatedWarehouses>("/warehouses?limit=100"),
    enabled,
  })
  const products = useQuery({
    queryKey: ["products", { limit: 100 }],
    queryFn: () => apiClient<PaginatedProducts>("/products?limit=100"),
    enabled,
  })

  return { warehouses, products }
}

export function useCreateInventory(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateInventoryRequest) =>
      apiClient<Inventory>("/inventory", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inventory"] })
      onSuccess?.()
    },
  })
}

export function useUpdateInventory(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateInventoryRequest }) =>
      apiClient<Inventory>(`/inventory/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inventory"] })
      onSuccess?.()
    },
  })
}

export function useDeleteInventory(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiClient<void>(`/inventory/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inventory"] })
      onSuccess?.()
    },
  })
}
