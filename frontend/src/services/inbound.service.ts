import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

export type Inbound = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"]
export type PaginatedInbounds = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.PaginatedResponse"]
export type CreateInboundRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.CreateInboundRequest"]
export type UpdateInboundRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.UpdateInboundRequest"]

export function listInbounds(limit: number) {
  return apiClient<PaginatedInbounds>(`/inbounds?limit=${limit}`)
}

export function createInbound(payload: CreateInboundRequest) {
  return apiClient<Inbound>("/inbounds", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function updateInbound(id: number, payload: UpdateInboundRequest) {
  return apiClient<Inbound>(`/inbounds/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
}

export function deleteInbound(id: number) {
  return apiClient<void>(`/inbounds/${id}`, { method: "DELETE" })
}
