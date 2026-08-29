"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import {
  listProducts,
  createProduct,
  updateProduct,
  deleteProduct,
  type CreateProductRequest,
  type UpdateProductRequest,
} from "@/services/product.service"

export function useProducts(limit = 100) {
  return useQuery({
    queryKey: ["products", { limit }],
    queryFn: () => listProducts(limit),
  })
}

export function useCreateProduct(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateProductRequest) => createProduct(payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["products"] })
      onSuccess?.()
    },
  })
}

export function useUpdateProduct(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateProductRequest }) => updateProduct(id, payload),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["products"] })
      onSuccess?.()
    },
  })
}

export function useDeleteProduct(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => deleteProduct(id),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["products"] })
      onSuccess?.()
    },
  })
}
