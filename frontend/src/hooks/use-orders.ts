"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listOrders,
  createOrder,
  updateOrder,
  deleteOrder,
  type CreateOrderRequest,
  type UpdateOrderRequest,
} from "@/services/order.service"

/**
 * Hook to retrieve orders list. By invalidating ["orders"], we can refetch this list.
 */
export function useOrders(workspaceId: string | undefined, limit = 100) {
  return useQuery({
    queryKey: ["orders", { workspaceId, limit }],
    queryFn: () => listOrders(limit, workspaceId),
  })
}

/**
 * Mutation to create a new order. Triggers full ["orders"] queries invalidation on success.
 */
export function useCreateOrder(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateOrderRequest) => createOrder(payload),
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
    mutationFn: ({ id, payload }: { id: number; payload: UpdateOrderRequest }) => updateOrder(id, payload),
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
    mutationFn: (id: number) => deleteOrder(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["orders"] })
      if (onSuccess) onSuccess()
    },
  })
}
