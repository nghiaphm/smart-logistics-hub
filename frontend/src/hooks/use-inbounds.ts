"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type Inbound = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"]
type PaginatedInbounds = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.PaginatedResponse"]
type CreateInboundRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.CreateInboundRequest"]
type UpdateInboundRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.UpdateInboundRequest"]

type PaginatedWarehouses = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"]
type PaginatedProducts = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse"]

export function useInbounds(limit = 100) {
  return useQuery({
    queryKey: ["inbounds", { limit }],
    queryFn: () => apiClient<PaginatedInbounds>(`/inbounds?limit=${limit}`),
  })
}

export function useInboundFormOptions(enabled: boolean) {
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

export function useCreateInbound(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateInboundRequest) =>
      apiClient<Inbound>("/inbounds", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inbounds"] })
      onSuccess?.()
    },
  })
}

export function useUpdateInbound(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateInboundRequest }) =>
      apiClient<Inbound>(`/inbounds/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inbounds"] })
      onSuccess?.()
    },
  })
}

export function useDeleteInbound(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiClient<void>(`/inbounds/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inbounds"] })
      onSuccess?.()
    },
  })
}
