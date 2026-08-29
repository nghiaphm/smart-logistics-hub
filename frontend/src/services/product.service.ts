import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Product = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"]
export type PaginatedProducts = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse"]
export type CreateProductRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.CreateProductRequest"]
export type UpdateProductRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.UpdateProductRequest"]

export function listProducts(limit: number) {
  return apiClient<PaginatedProducts>(`/products?limit=${limit}`)
}

export function createProduct(payload: CreateProductRequest) {
  return apiClient<Product>("/products", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateProduct(id: number, payload: UpdateProductRequest) {
  return apiClient<Product>(`/products/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteProduct(id: number) {
  return apiClient<void>(`/products/${id}`, { method: "DELETE" })
}
