"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, Add01Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import { useCreateTrip, useUpdateTrip } from "@/hooks/use-trips"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type Trip = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"]
type CreateTripRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.CreateTripRequest"]
type UpdateTripRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.UpdateTripRequest"]
type TripStopRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripStopRequest"]

type TripFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  trip?: Trip
  onSuccess?: () => void
}

type FormStop = {
  orderCode: string
  stopType: string
  address: string
  status: string
  lat: string
  lng: string
}

const STATUS_OPTIONS = ["PLANNED", "IN_TRANSIT", "COMPLETED"]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function emptyStop(): FormStop {
  return { orderCode: "", stopType: "", address: "", status: "", lat: "", lng: "" }
}

function toCoordinate(value: string): number | undefined {
  if (value.trim() === "") return undefined
  return Number(value) || undefined
}

export function TripFormModal({ open, onOpenChange, trip, onSuccess }: TripFormModalProps) {
  const isEdit = Boolean(trip)
  const [tripCode, setTripCode] = useState(() => {
    if (trip) return trip.trip_code ?? ""
    return `TRIP-${new Date().getFullYear().toString().slice(-2)}${(new Date().getMonth() + 1).toString().padStart(2, "0")}-${Math.floor(1000 + Math.random() * 9000)}`
  })
  const [driverCode, setDriverCode] = useState(() => (trip ? String(trip.driver_id ?? "") : ""))
  const [vehicleLicensePlate, setVehicleLicensePlate] = useState(trip?.vehicle_license_plate ?? "")
  const [totalDistanceKm, setTotalDistanceKm] = useState(
    trip?.total_distance_km != null ? String(trip.total_distance_km) : ""
  )
  const [status, setStatus] = useState(trip?.status ?? "PLANNED")
  const [stops, setStops] = useState<FormStop[]>(() => {
    if (trip) {
      const mapped = (trip.stops ?? []).map((stop) => ({
        orderCode: stop.order_code ?? "",
        stopType: stop.stop_type ?? "",
        address: stop.address ?? "",
        status: stop.status ?? "",
        lat: stop.lat != null ? String(stop.lat) : "",
        lng: stop.lng != null ? String(stop.lng) : "",
      }))
      return mapped.length > 0 ? mapped : [emptyStop()]
    }
    return [emptyStop()]
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const createMutation = useCreateTrip(() => {
    toast.add({ title: "Đã tạo chuyến xe", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateTrip(() => {
    toast.add({ title: "Đã cập nhật chuyến xe", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const handleAddStop = () => {
    setStops((current) => [...current, emptyStop()])
  }

  const handleRemoveStop = (index: number) => {
    if (stops.length <= 1) return
    setStops((current) => current.filter((_, stopIndex) => stopIndex !== index))
  }

  const handleStopChange = (index: number, field: keyof FormStop, value: string) => {
    setStops((current) =>
      current.map((stop, stopIndex) => (stopIndex === index ? { ...stop, [field]: value } : stop))
    )
  }

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!isEdit && !tripCode.trim()) nextErrors.tripCode = "Mã chuyến không được bỏ trống."
    if (!isEdit && !driverCode.trim()) nextErrors.driverCode = "Mã tài xế không được bỏ trống."
    const hasInvalidStop = stops.some((stop) => !stop.orderCode.trim() || !stop.address.trim())
    if (hasInvalidStop) nextErrors.stops = "Mỗi điểm dừng phải có mã đơn và địa chỉ."
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!validate()) {
      toast.add({
        title: "Thông tin không hợp lệ",
        description: "Vui lòng kiểm tra các trường bắt buộc.",
        type: "error",
      })
      return
    }

    const stopPayload: TripStopRequest[] = stops.map((stop) => ({
      order_code: stop.orderCode.trim(),
      stop_type: stop.stopType.trim() || undefined,
      address: stop.address.trim(),
      status: stop.status.trim() || undefined,
      location:
        stop.lat.trim() || stop.lng.trim()
          ? { lat: toCoordinate(stop.lat), lng: toCoordinate(stop.lng) }
          : undefined,
    }))

    if (isEdit && trip?.id) {
      const payload: UpdateTripRequest = { status, stops: stopPayload }
      updateMutation.mutate(
        { id: trip.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật chuyến xe", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    const payload: CreateTripRequest = {
      trip_code: tripCode.trim(),
      driver_code: driverCode.trim(),
      vehicle_license_plate: vehicleLicensePlate.trim() || undefined,
      status,
      total_distance_km: totalDistanceKm.trim() ? Number(totalDistanceKm) : undefined,
      stops: stopPayload,
    }
    createMutation.mutate(payload, {
      onError: (error) =>
        toast.add({ title: "Không thể tạo chuyến xe", description: errorMessage(error), type: "error" }),
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa chuyến xe" : "Tạo chuyến xe mới"}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? "Cập nhật trạng thái và danh sách điểm dừng của chuyến xe."
              : "Nhập thông tin chuyến vận chuyển và các điểm dừng."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Mã chuyến" htmlFor="trip-code" error={errors.tripCode} required={!isEdit}>
              <Input
                id="trip-code"
                value={tripCode}
                onChange={(event) => setTripCode(event.target.value)}
                placeholder="VD: TRIP-2408-001"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
            <FormField label="Mã tài xế" htmlFor="trip-driver" error={errors.driverCode} required={!isEdit}>
              <Input
                id="trip-driver"
                value={driverCode}
                onChange={(event) => setDriverCode(event.target.value)}
                placeholder="VD: DRV-001"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Biển số xe" htmlFor="trip-plate">
              <Input
                id="trip-plate"
                value={vehicleLicensePlate}
                onChange={(event) => setVehicleLicensePlate(event.target.value)}
                placeholder="VD: 51A-12345"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
            <FormField label="Quãng đường (km)" htmlFor="trip-distance">
              <Input
                id="trip-distance"
                type="number"
                min="0"
                value={totalDistanceKm}
                onChange={(event) => setTotalDistanceKm(event.target.value)}
                placeholder="0"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
          </div>

          <FormField label="Trạng thái" htmlFor="trip-status">
            <select
              id="trip-status"
              value={status}
              onChange={(event) => setStatus(event.target.value)}
              disabled={isSubmitting}
              className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
            >
              {STATUS_OPTIONS.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </select>
          </FormField>

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold uppercase tracking-wider text-neutral-800 dark:text-neutral-200">
                Điểm dừng
              </h3>
              <Button type="button" variant="outline" size="sm" onClick={handleAddStop} disabled={isSubmitting} className="gap-1 text-xs">
                <HugeiconsIcon icon={Add01Icon} className="size-3" /> Thêm điểm
              </Button>
            </div>
            {errors.stops ? (
              <p className="flex items-center gap-1.5 text-xs text-destructive">
                <HugeiconsIcon icon={Alert02Icon} className="size-3.5" />
                {errors.stops}
              </p>
            ) : null}
            {stops.map((stop, index) => (
              <div key={index} className="flex flex-wrap items-end gap-2 rounded-2xl border border-border/60 bg-muted/10 p-3">
                <FormField label={`Mã đơn ${index + 1}`} className="min-w-36 flex-1">
                  <Input
                    value={stop.orderCode}
                    onChange={(event) => handleStopChange(index, "orderCode", event.target.value)}
                    placeholder="VD: DH-2408-001"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Địa chỉ" className="min-w-44 flex-1">
                  <Input
                    value={stop.address}
                    onChange={(event) => handleStopChange(index, "address", event.target.value)}
                    placeholder="Địa chỉ điểm dừng"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Loại điểm" className="w-28">
                  <Input
                    value={stop.stopType}
                    onChange={(event) => handleStopChange(index, "stopType", event.target.value)}
                    placeholder="PICKUP"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Trạng thái" className="w-28">
                  <Input
                    value={stop.status}
                    onChange={(event) => handleStopChange(index, "status", event.target.value)}
                    placeholder="PENDING"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Lat" className="w-24">
                  <Input
                    type="number"
                    step="any"
                    value={stop.lat}
                    onChange={(event) => handleStopChange(index, "lat", event.target.value)}
                    placeholder="Vĩ độ"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Lng" className="w-24">
                  <Input
                    type="number"
                    step="any"
                    value={stop.lng}
                    onChange={(event) => handleStopChange(index, "lng", event.target.value)}
                    placeholder="Kinh độ"
                    disabled={isSubmitting}
                  />
                </FormField>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => handleRemoveStop(index)}
                  className="h-9 w-9 shrink-0 text-muted-foreground hover:text-destructive"
                  disabled={isSubmitting || stops.length <= 1}
                >
                  <HugeiconsIcon icon={Delete01Icon} className="size-4" />
                </Button>
              </div>
            ))}
          </div>

          <FormActions>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>
              Huỷ
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Tạo chuyến xe"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
