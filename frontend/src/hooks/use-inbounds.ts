"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listInbounds,
  createInbound,
  updateInbound,
  deleteInbound,
  type CreateInboundRequest,
  type UpdateInboundRequest,
} from "@/services/inbound.service"
import { listWarehouses } from "@/services/warehouse.service"
import { listProducts } from "@/services/product.service"

export function useInbounds(limit = 100) {
  return useQuery({
    queryKey: ["inbounds", { limit }],
    queryFn: () => listInbounds(limit),
  })
}

export function useInboundFormOptions(enabled: boolean) {
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

export function useCreateInbound(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateInboundRequest) => createInbound(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inbounds"] })
      onSuccess?.()
    },
  })
}

export function useUpdateInbound(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateInboundRequest }) => updateInbound(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inbounds"] })
      onSuccess?.()
    },
  })
}

export function useDeleteInbound(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteInbound(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["inbounds"] })
      onSuccess?.()
    },
  })
}
