"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type Product = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"]
type PaginatedProducts = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse"]
type CreateProductRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.CreateProductRequest"]
type UpdateProductRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.UpdateProductRequest"]

export function useProducts(limit = 100) {
  return useQuery({
    queryKey: ["products", { limit }],
    queryFn: () => apiClient<PaginatedProducts>(`/products?limit=${limit}`),
  })
}

export function useCreateProduct(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateProductRequest) =>
      apiClient<Product>("/products", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["products"] })
      onSuccess?.()
    },
  })
}

export function useUpdateProduct(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateProductRequest }) =>
      apiClient<Product>(`/products/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["products"] })
      onSuccess?.()
    },
  })
}

export function useDeleteProduct(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => apiClient<void>(`/products/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["products"] })
      onSuccess?.()
    },
  })
}
