"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import { useCreateDriver, useUpdateDriver } from "@/hooks/use-drivers"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type Driver = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"]
type CreateDriverRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.CreateDriverRequest"]
type UpdateDriverRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.UpdateDriverRequest"]

type DriverFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  driver?: Driver
  onSuccess?: () => void
}

const STATUS_OPTIONS = ["AVAILABLE", "BUSY", "OFFLINE"] as const
type DriverStatus = (typeof STATUS_OPTIONS)[number]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export function DriverFormModal({ open, onOpenChange, driver, onSuccess }: DriverFormModalProps) {
  const isEdit = Boolean(driver)
  const [driverCode, setDriverCode] = useState(driver?.driver_code ?? "")
  const [fullName, setFullName] = useState(driver?.full_name ?? "")
  const [phone, setPhone] = useState(driver?.phone ?? "")
  const [vehicleType, setVehicleType] = useState(driver?.vehicle_type ?? "")
  const [licensePlate, setLicensePlate] = useState(driver?.license_plate ?? "")
  const [status, setStatus] = useState<DriverStatus>((driver?.status ?? "AVAILABLE") as DriverStatus)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const createMutation = useCreateDriver(() => {
    toast.add({ title: "Đã thêm tài xế", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateDriver(() => {
    toast.add({ title: "Đã cập nhật tài xế", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!driverCode.trim()) nextErrors.driverCode = "Mã tài xế không được bỏ trống."
    if (!fullName.trim()) nextErrors.fullName = "Họ tên không được bỏ trống."
    if (!phone.trim()) nextErrors.phone = "Điện thoại không được bỏ trống."
    if (!vehicleType.trim()) nextErrors.vehicleType = "Loại xe không được bỏ trống."
    if (!licensePlate.trim()) nextErrors.licensePlate = "Biển số không được bỏ trống."
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

    const payload: CreateDriverRequest & UpdateDriverRequest = {
      driver_code: driverCode.trim(),
      full_name: fullName.trim(),
      phone: phone.trim(),
      vehicle_type: vehicleType.trim(),
      license_plate: licensePlate.trim(),
      status,
    }

    if (isEdit && driver?.id) {
      updateMutation.mutate(
        { id: driver.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật tài xế", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    createMutation.mutate(payload, {
      onError: (error) =>
        toast.add({ title: "Không thể thêm tài xế", description: errorMessage(error), type: "error" }),
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa tài xế" : "Thêm tài xế mới"}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? "Cập nhật thông tin tài xế vận tải."
              : "Đăng ký tài xế mới vào đội xe."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Mã tài xế" htmlFor="driver-code" error={errors.driverCode} required>
              <Input
                id="driver-code"
                value={driverCode}
                onChange={(event) => setDriverCode(event.target.value)}
                placeholder="VD: DRV-001"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Họ tên" htmlFor="driver-name" error={errors.fullName} required>
              <Input
                id="driver-name"
                value={fullName}
                onChange={(event) => setFullName(event.target.value)}
                placeholder="VD: Nguyễn Văn An"
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Điện thoại" htmlFor="driver-phone" error={errors.phone} required>
              <Input
                id="driver-phone"
                value={phone}
                onChange={(event) => setPhone(event.target.value)}
                placeholder="VD: 0901234567"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Loại xe" htmlFor="driver-vehicle-type" error={errors.vehicleType} required>
              <Input
                id="driver-vehicle-type"
                value={vehicleType}
                onChange={(event) => setVehicleType(event.target.value)}
                placeholder="VD: TRUCK, VAN..."
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Biển số xe" htmlFor="driver-license-plate" error={errors.licensePlate} required>
              <Input
                id="driver-license-plate"
                value={licensePlate}
                onChange={(event) => setLicensePlate(event.target.value)}
                placeholder="VD: 51A-123.45"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Trạng thái" htmlFor="driver-status">
              <select
                id="driver-status"
                value={status}
                onChange={(event) => setStatus(event.target.value as DriverStatus)}
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
              {isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Thêm tài xế"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
