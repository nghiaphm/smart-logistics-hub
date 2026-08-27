"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type OrderResponse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"]
type PaginatedOrders = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.PaginatedResponse"]
type CreateOrderRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.CreateOrderRequest"]
type UpdateOrderRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.UpdateOrderRequest"]

/**
 * Hook to retrieve orders list. By invalidating ["orders"], we can refetch this list.
 */
export function useOrders(workspaceId: string, limit = 100) {
  return useQuery({
    queryKey: ["orders", { workspaceId, limit }],
    queryFn: () => apiClient<PaginatedOrders>(`/orders?limit=${limit}`),
  })
}

/**
 * Mutation to create a new order. Triggers full ["orders"] queries invalidation on success.
 */
export function useCreateOrder(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateOrderRequest) =>
      apiClient<OrderResponse>("/orders", {
        method: "POST",
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["orders"] })
      if (onSuccess) onSuccess()
    },
  })
}

/**
 * Mutation to update an existing order. Triggers full ["orders"] queries invalidation on success.
 */
export function useUpdateOrder(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateOrderRequest }) =>
      apiClient<OrderResponse>(`/orders/${id}`, {
        method: "PATCH",
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["orders"] })
      if (onSuccess) onSuccess()
    },
  })
}

/**
 * Mutation to delete an existing order. Triggers full ["orders"] queries invalidation on success.
 */
export function useDeleteOrder(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) =>
      apiClient<void>(`/orders/${id}`, {
        method: "DELETE",
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["orders"] })
      if (onSuccess) onSuccess()
    },
  })
}
