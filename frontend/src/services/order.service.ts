import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type OrderResponse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"]
export type PaginatedOrders = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.PaginatedResponse"]
export type CreateOrderRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.CreateOrderRequest"]
export type UpdateOrderRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.UpdateOrderRequest"]

export function listOrders(limit: number) {
  return apiClient<PaginatedOrders>(`/orders?limit=${limit}`)
}

export function createOrder(payload: CreateOrderRequest) {
  return apiClient<OrderResponse>("/orders", {
    method: "POST",
    body: JSON.stringify(payload),
  })
}

export function updateOrder(id: number, payload: UpdateOrderRequest) {
  return apiClient<OrderResponse>(`/orders/${id}`, {
    method: "PATCH",
    body: JSON.stringify(payload),
  })
}

export function deleteOrder(id: number) {
  return apiClient<void>(`/orders/${id}`, { method: "DELETE" })
}
