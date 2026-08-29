"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import { useCreateVehicle, useUpdateVehicle } from "@/hooks/use-vehicles"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select } from "@/components/shared/form/Select"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type Vehicle = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.VehicleResponse"]
type CreateVehicleRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.CreateVehicleRequest"]
type UpdateVehicleRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.UpdateVehicleRequest"]

type VehicleFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  vehicle?: Vehicle
  onSuccess?: () => void
}

const STATUS_OPTIONS = ["ACTIVE", "MAINTENANCE", "INACTIVE"] as const
type VehicleStatus = (typeof STATUS_OPTIONS)[number]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function toCapacity(value: string): number | undefined {
  if (value.trim() === "") return undefined
  return Number(value)
}

export function VehicleFormModal({ open, onOpenChange, vehicle, onSuccess }: VehicleFormModalProps) {
  const isEdit = Boolean(vehicle)
  const [licensePlate, setLicensePlate] = useState(vehicle?.license_plate ?? "")
  const [type, setType] = useState(vehicle?.type ?? "")
  const [capacity, setCapacity] = useState(vehicle?.capacity != null ? String(vehicle.capacity) : "")
  const [status, setStatus] = useState<VehicleStatus>((vehicle?.status ?? "ACTIVE") as VehicleStatus)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const createMutation = useCreateVehicle(() => {
    toast.add({ title: "Đã thêm phương tiện", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateVehicle(() => {
    toast.add({ title: "Đã cập nhật phương tiện", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!licensePlate.trim()) nextErrors.licensePlate = "Biển số không được bỏ trống."
    if (capacity.trim() !== "" && (Number.isNaN(Number(capacity)) || Number(capacity) < 0)) {
      nextErrors.capacity = "Tải trọng phải là số không âm."
    }
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
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

    const payload: CreateVehicleRequest & UpdateVehicleRequest = {
      license_plate: licensePlate.trim(),
      type: type.trim() || undefined,
      capacity: toCapacity(capacity),
      status,
    }

    if (isEdit && vehicle?.id) {
      updateMutation.mutate(
        { id: vehicle.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật phương tiện", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    createMutation.mutate(payload, {
      onError: (error) =>
        toast.add({ title: "Không thể thêm phương tiện", description: errorMessage(error), type: "error" }),
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa phương tiện" : "Thêm phương tiện mới"}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? "Cập nhật thông tin phương tiện vận tải."
              : "Đăng ký phương tiện mới vào đội xe (Fleet)."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Biển số" htmlFor="vehicle-license-plate" error={errors.licensePlate} required>
              <Input
                id="vehicle-license-plate"
                value={licensePlate}
                onChange={(event) => setLicensePlate(event.target.value)}
                placeholder="VD: 51F-123.45"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Loại xe" htmlFor="vehicle-type">
              <Input
                id="vehicle-type"
                value={type}
                onChange={(event) => setType(event.target.value)}
                placeholder="VD: TRUCK, VAN, MOTORBIKE..."
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Tải trọng (kg)" htmlFor="vehicle-capacity" error={errors.capacity} hint="Khối lượng hàng tối đa, đơn vị kg.">
              <Input
                id="vehicle-capacity"
                type="number"
                min="0"
                step="any"
                value={capacity}
                onChange={(event) => setCapacity(event.target.value)}
                placeholder="VD: 1500"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Trạng thái" htmlFor="vehicle-status">
              <Select
                id="vehicle-status"
                value={status}
                onChange={(event) => setStatus(event.target.value as VehicleStatus)}
                disabled={isSubmitting}
              >
                {STATUS_OPTIONS.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </Select>
            </FormField>
          </div>

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
              {isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Thêm phương tiện"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
