"use client"

import { useState } from "react"
import type { FormEvent } from "react"
import { useParams } from "next/navigation"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import { useCreateTrackingEvent, useUpdateTrackingEvent } from "@/hooks/use-tracking"
import { useOrders } from "@/hooks/use-orders"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type TrackingEvent = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"]
type CreateTrackingEventRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.CreateTrackingEventRequest"]
type UpdateTrackingEventRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.UpdateTrackingEventRequest"]

type TrackingFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  event?: TrackingEvent
  onSuccess?: () => void
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function toCoordinate(value: string): number | undefined {
  if (value.trim() === "") return undefined
  return Number(value) || undefined
}

export function TrackingFormModal({
  open,
  onOpenChange,
  event,
  onSuccess,
}: TrackingFormModalProps) {
  const params = useParams<{ workspace_id: string }>()
  const workspaceId = params.workspace_id
  const { data: ordersData } = useOrders(workspaceId, 100)
  const orders = ordersData?.items ?? []

  const isEdit = Boolean(event)
  const [orderCode, setOrderCode] = useState(event?.order_code ?? "")
  const [selectedOrderId, setSelectedOrderId] = useState(
    event?.order_id != null ? String(event.order_id) : ""
  )
  const [driverCode, setDriverCode] = useState(event?.driver_code ?? "")
  const [statusUpdate, setStatusUpdate] = useState(event?.status_update ?? "")
  const [lat, setLat] = useState(event?.lat != null ? String(event.lat) : "")
  const [lng, setLng] = useState(event?.lng != null ? String(event.lng) : "")
  const [note, setNote] = useState(event?.note ?? "")
  const [errors, setErrors] = useState<Record<string, string>>({})

  const createMutation = useCreateTrackingEvent(() => {
    toast.add({ title: "Đã ghi nhận sự kiện theo dõi", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateTrackingEvent(() => {
    toast.add({ title: "Đã cập nhật sự kiện theo dõi", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!orderCode.trim()) nextErrors.orderCode = "Mã đơn không được bỏ trống."
    if (!driverCode.trim()) nextErrors.driverCode = "Mã tài xế không được bỏ trống."
    if (!statusUpdate.trim()) nextErrors.statusUpdate = "Trạng thái không được bỏ trống."
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const handleOrderSelect = (value: string) => {
    setSelectedOrderId(value)
    const order = orders.find((o) => o.id != null && String(o.id) === value)
    if (order?.order_code) setOrderCode(order.order_code)
  }

  const handleSubmit = (eventForm: FormEvent<HTMLFormElement>) => {
    eventForm.preventDefault()
    if (!validate()) {
      toast.add({
        title: "Thông tin không hợp lệ",
        description: "Vui lòng kiểm tra các trường bắt buộc.",
        type: "error",
      })
      return
    }

    const payload: CreateTrackingEventRequest & UpdateTrackingEventRequest = {
      order_code: orderCode.trim(),
      order_id: selectedOrderId ? Number(selectedOrderId) : undefined,
      driver_code: driverCode.trim(),
      status_update: statusUpdate.trim(),
      lat: toCoordinate(lat),
      lng: toCoordinate(lng),
      note: note.trim() || undefined,
    }

    if (isEdit && event?.id) {
      updateMutation.mutate(
        { id: event.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật sự kiện theo dõi", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    createMutation.mutate(payload, {
      onError: (error) =>
        toast.add({ title: "Không thể ghi nhận sự kiện theo dõi", description: errorMessage(error), type: "error" }),
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa sự kiện theo dõi" : "Ghi nhận sự kiện theo dõi"}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? "Cập nhật thông tin sự kiện theo dõi vận chuyển."
              : "Ghi lại trạng thái vận chuyển mới nhất của đơn hàng."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Mã đơn" htmlFor="tracking-order-code" error={errors.orderCode} required>
              <Input
                id="tracking-order-code"
                value={orderCode}
                onChange={(event) => {
                  setOrderCode(event.target.value)
                  setSelectedOrderId("")
                }}
                placeholder="VD: DH-2408-001"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Mã tài xế" htmlFor="tracking-driver-code" error={errors.driverCode} required>
              <Input
                id="tracking-driver-code"
                value={driverCode}
                onChange={(event) => setDriverCode(event.target.value)}
                placeholder="VD: DRV-001"
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <FormField
            label="Đơn hàng liên kết"
            htmlFor="tracking-order"
            hint="Chọn đơn hàng để tự động điền mã đơn và liên kết theo dõi với đơn đó."
          >
            <select
              id="tracking-order"
              value={selectedOrderId}
              onChange={(event) => handleOrderSelect(event.target.value)}
              disabled={isSubmitting}
              className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
            >
              <option value="">Chưa liên kết (nhập mã đơn tay)</option>
              {orders.map((order) => (
                <option key={order.id ?? order.order_code} value={order.id ?? ""}>
                  {order.order_code}
                </option>
              ))}
            </select>
          </FormField>

          <FormField label="Trạng thái" htmlFor="tracking-status" error={errors.statusUpdate} required>
            <Input
              id="tracking-status"
              value={statusUpdate}
              onChange={(event) => setStatusUpdate(event.target.value)}
              placeholder="VD: PICKED_UP, IN_TRANSIT, DELIVERED..."
              disabled={isSubmitting}
            />
          </FormField>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Vĩ độ (lat)" htmlFor="tracking-lat" hint="Toạ độ GPS, có thể âm.">
              <Input
                id="tracking-lat"
                type="number"
                step="any"
                value={lat}
                onChange={(event) => setLat(event.target.value)}
                placeholder="VD: 10.8231"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Kinh độ (lng)" htmlFor="tracking-lng" hint="Toạ độ GPS, có thể âm.">
              <Input
                id="tracking-lng"
                type="number"
                step="any"
                value={lng}
                onChange={(event) => setLng(event.target.value)}
                placeholder="VD: 106.6297"
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <FormField label="Ghi chú" htmlFor="tracking-note">
            <Input
              id="tracking-note"
              value={note}
              onChange={(event) => setNote(event.target.value)}
              placeholder="Ghi chú thêm về sự kiện"
              disabled={isSubmitting}
            />
          </FormField>

          {Object.keys(errors).length > 0 ? (
            <p className="flex items-center gap-1 text-xs text-destructive">
              <HugeiconsIcon icon={Alert02Icon} className="size-3.5" />
              Vui lòng kiểm tra các trường bắt buộc.
            </p>
          ) : null}

          <FormActions>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>
              Huỷ
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Ghi nhận"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
