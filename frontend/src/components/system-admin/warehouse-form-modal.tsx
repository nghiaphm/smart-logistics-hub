"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { Form, FormField } from "@/components/shared/form/Form"
import { AppModalActions, AppModalShell } from "@/components/shared/modal"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import type { components } from "@/types/api"

type Warehouse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"]

export type WarehouseFormValues = {
  warehouse_code: string
  name: string
  address: string
  manager_name: string
  contact_phone: string
  lat: string
  lng: string
  is_active: boolean
}

type WarehouseFormModalProps = {
  open: boolean
  mode: "create" | "edit"
  warehouse?: Warehouse
  isSubmitting: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (values: WarehouseFormValues) => void
}

type FormErrors = Partial<Record<"warehouse_code" | "name" | "address", string>>

function toValues(warehouse?: Warehouse): WarehouseFormValues {
  return {
    warehouse_code: warehouse?.warehouse_code ?? "",
    name: warehouse?.name ?? "",
    address: warehouse?.address ?? "",
    manager_name: warehouse?.manager_name ?? "",
    contact_phone: warehouse?.contact_phone ?? "",
    lat: warehouse?.lat != null ? String(warehouse.lat) : "",
    lng: warehouse?.lng != null ? String(warehouse.lng) : "",
    is_active: warehouse?.is_active ?? true,
  }
}

export function WarehouseFormModal({
  open,
  mode,
  warehouse,
  isSubmitting,
  onOpenChange,
  onSubmit,
}: WarehouseFormModalProps) {
  const [values, setValues] = useState<WarehouseFormValues>(() => toValues(warehouse))
  const [errors, setErrors] = useState<FormErrors>({})

  const isEdit = mode === "edit"

  const setField = (field: keyof WarehouseFormValues, value: string | boolean) => {
    setValues((prev) => ({ ...prev, [field]: value }))
    setErrors((prev) => ({ ...prev, [field]: undefined }))
  }

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const nextErrors: FormErrors = {}
    if (!values.warehouse_code.trim()) {
      nextErrors.warehouse_code = "Vui lòng nhập mã kho"
    }
    if (!values.name.trim()) {
      nextErrors.name = "Vui lòng nhập tên kho"
    }
    if (!values.address.trim()) {
      nextErrors.address = "Vui lòng nhập địa chỉ"
    }
    setErrors(nextErrors)
    if (Object.keys(nextErrors).length > 0) {
      return
    }
    onSubmit(values)
  }

  return (
    <AppModalShell
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? "Sửa kho" : "Thêm kho"}
      description={
        isEdit
          ? `Cập nhật thông tin kho ${warehouse?.warehouse_code ?? ""}`
          : "Nhập thông tin kho mới"
      }
      actions={
        <AppModalActions>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>
            Huỷ
          </Button>
          <Button type="submit" form="warehouse-form" disabled={isSubmitting}>
            {isSubmitting ? "Đang lưu…" : isEdit ? "Lưu thay đổi" : "Thêm kho"}
          </Button>
        </AppModalActions>
      }
    >
      <Form id="warehouse-form" onSubmit={handleSubmit}>
        {!isEdit ? (
          <FormField
            label="Mã kho"
            htmlFor="wh-code"
            required
            error={errors.warehouse_code}
          >
            <Input
              id="wh-code"
              value={values.warehouse_code}
              onChange={(event) => setField("warehouse_code", event.target.value)}
              placeholder="VD: WH-001"
            />
          </FormField>
        ) : null}
        <FormField label="Tên kho" htmlFor="wh-name" required error={errors.name}>
          <Input
            id="wh-name"
            value={values.name}
            onChange={(event) => setField("name", event.target.value)}
            placeholder="VD: Kho VSIP 1"
          />
        </FormField>
        <FormField label="Địa chỉ" htmlFor="wh-address" required error={errors.address}>
          <Input
            id="wh-address"
            value={values.address}
            onChange={(event) => setField("address", event.target.value)}
            placeholder="VD: Khu công nghiệp VSIP 1, Bình Dương"
          />
        </FormField>
        <div className="grid gap-5 sm:grid-cols-2">
          <FormField label="Người quản lý" htmlFor="wh-manager">
            <Input
              id="wh-manager"
              value={values.manager_name}
              onChange={(event) => setField("manager_name", event.target.value)}
            />
          </FormField>
          <FormField label="Số điện thoại" htmlFor="wh-phone">
            <Input
              id="wh-phone"
              value={values.contact_phone}
              onChange={(event) => setField("contact_phone", event.target.value)}
            />
          </FormField>
        </div>
        <div className="grid gap-5 sm:grid-cols-2">
          <FormField label="Vĩ độ (lat)" htmlFor="wh-lat" hint="Không bắt buộc">
            <Input
              id="wh-lat"
              value={values.lat}
              onChange={(event) => setField("lat", event.target.value)}
              placeholder="VD: 10.8231"
              inputMode="decimal"
            />
          </FormField>
          <FormField label="Kinh độ (lng)" htmlFor="wh-lng" hint="Không bắt buộc">
            <Input
              id="wh-lng"
              value={values.lng}
              onChange={(event) => setField("lng", event.target.value)}
              placeholder="VD: 106.6297"
              inputMode="decimal"
            />
          </FormField>
        </div>
        {isEdit ? (
          <label className="flex cursor-pointer items-center gap-2 text-sm text-foreground">
            <Checkbox
              checked={values.is_active}
              onCheckedChange={(checked) => setField("is_active", checked)}
            />
            Kho đang hoạt động
          </label>
        ) : null}
      </Form>
    </AppModalShell>
  )
}
