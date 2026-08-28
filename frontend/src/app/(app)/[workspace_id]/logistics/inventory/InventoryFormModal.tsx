"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import {
  useCreateInventory,
  useInventoryFormOptions,
  useUpdateInventory,
} from "@/hooks/use-inventory"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type Inventory = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"]
type CreateInventoryRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.CreateInventoryRequest"]
type UpdateInventoryRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.UpdateInventoryRequest"]

type InventoryFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  inventory?: Inventory
  onSuccess?: () => void
}

type QuantityFields = {
  availableQty: number
  reservedQty: number
  damagedQty: number
  holdQty: number
}

const initialQuantities: QuantityFields = {
  availableQty: 0,
  reservedQty: 0,
  damagedQty: 0,
  holdQty: 0,
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export function InventoryFormModal({
  open,
  onOpenChange,
  inventory,
  onSuccess,
}: InventoryFormModalProps) {
  const isEdit = Boolean(inventory)
  const [productId, setProductId] = useState(() => (inventory ? String(inventory.product_id ?? "") : ""))
  const [warehouseId, setWarehouseId] = useState(() => (inventory ? String(inventory.warehouse_id ?? "") : ""))
  const [quantities, setQuantities] = useState<QuantityFields>(() =>
    inventory
      ? {
          availableQty: inventory.available_qty ?? 0,
          reservedQty: inventory.reserved_qty ?? 0,
          damagedQty: inventory.damaged_qty ?? 0,
          holdQty: inventory.hold_qty ?? 0,
        }
      : initialQuantities
  )
  const [updatedBy, setUpdatedBy] = useState(() => inventory?.updated_by ?? "")
  const [errors, setErrors] = useState<Record<string, string>>({})
  const { warehouses, products } = useInventoryFormOptions(open && !isEdit)

  const createMutation = useCreateInventory(() => {
    toast.add({ title: "Đã tạo bản ghi tồn kho", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateInventory(() => {
    toast.add({ title: "Đã cập nhật tồn kho", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const updateQuantity = (field: keyof QuantityFields, value: string) => {
    setQuantities((current) => ({
      ...current,
      [field]: Math.max(0, Number.parseInt(value, 10) || 0),
    }))
  }

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!isEdit && !productId) nextErrors.productId = "Vui lòng chọn sản phẩm."
    if (!isEdit && !warehouseId) nextErrors.warehouseId = "Vui lòng chọn kho."
    if (Object.values(quantities).some((quantity) => quantity < 0)) {
      nextErrors.quantities = "Số lượng không được là số âm."
    }
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

    if (isEdit && inventory?.id) {
      const payload: UpdateInventoryRequest = {
        available_qty: quantities.availableQty,
        reserved_qty: quantities.reservedQty,
        damaged_qty: quantities.damagedQty,
        hold_qty: quantities.holdQty,
        updated_by: updatedBy.trim() || undefined,
      }
      updateMutation.mutate(
        { id: inventory.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật tồn kho", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    const payload: CreateInventoryRequest = {
      product_id: Number(productId),
      warehouse_id: Number(warehouseId),
      available_qty: quantities.availableQty,
      reserved_qty: quantities.reservedQty,
      damaged_qty: quantities.damagedQty,
      hold_qty: quantities.holdQty,
      updated_by: updatedBy.trim() || undefined,
    }
    createMutation.mutate(payload, {
      onError: (error) =>
        toast.add({ title: "Không thể tạo tồn kho", description: errorMessage(error), type: "error" }),
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa bản ghi tồn kho" : "Tạo bản ghi tồn kho"}</DialogTitle>
          <DialogDescription>
            {isEdit ? "Cập nhật số lượng tồn kho hiện tại." : "Chọn sản phẩm, kho và nhập số lượng ban đầu."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          {!isEdit ? (
            <div className="grid gap-4 sm:grid-cols-2">
              <FormField label="Sản phẩm" htmlFor="inventory-product" error={errors.productId} required>
                <select
                  id="inventory-product"
                  value={productId}
                  onChange={(event) => setProductId(event.target.value)}
                  disabled={isSubmitting}
                  className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
                >
                  <option value="">Chọn sản phẩm</option>
                  {(products.data?.items ?? []).map((product) => (
                    <option key={product.id} value={product.id}>
                      {product.name} ({product.sku})
                    </option>
                  ))}
                </select>
              </FormField>
              <FormField label="Kho" htmlFor="inventory-warehouse" error={errors.warehouseId} required>
                <select
                  id="inventory-warehouse"
                  value={warehouseId}
                  onChange={(event) => setWarehouseId(event.target.value)}
                  disabled={isSubmitting}
                  className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
                >
                  <option value="">Chọn kho</option>
                  {(warehouses.data?.items ?? []).map((warehouse) => (
                    <option key={warehouse.id} value={warehouse.id}>
                      {warehouse.name} ({warehouse.warehouse_code})
                    </option>
                  ))}
                </select>
              </FormField>
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              <FormField label="Sản phẩm">
                <Input value={productId} disabled />
              </FormField>
              <FormField label="Kho">
                <Input value={warehouseId} disabled />
              </FormField>
            </div>
          )}

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Tồn sẵn sàng" htmlFor="available-qty" error={errors.quantities}>
              <Input id="available-qty" type="number" min="0" value={quantities.availableQty} onChange={(event) => updateQuantity("availableQty", event.target.value)} disabled={isSubmitting} />
            </FormField>
            <FormField label="Giữ chỗ" htmlFor="reserved-qty">
              <Input id="reserved-qty" type="number" min="0" value={quantities.reservedQty} onChange={(event) => updateQuantity("reservedQty", event.target.value)} disabled={isSubmitting} />
            </FormField>
            <FormField label="Hư hỏng" htmlFor="damaged-qty">
              <Input id="damaged-qty" type="number" min="0" value={quantities.damagedQty} onChange={(event) => updateQuantity("damagedQty", event.target.value)} disabled={isSubmitting} />
            </FormField>
            <FormField label="Chờ duyệt" htmlFor="hold-qty">
              <Input id="hold-qty" type="number" min="0" value={quantities.holdQty} onChange={(event) => updateQuantity("holdQty", event.target.value)} disabled={isSubmitting} />
            </FormField>
          </div>

          <FormField label="Người cập nhật" htmlFor="updated-by">
            <Input id="updated-by" value={updatedBy} onChange={(event) => setUpdatedBy(event.target.value)} placeholder="Tên người cập nhật" disabled={isSubmitting} />
          </FormField>

          {Object.keys(errors).length > 0 && errors.quantities ? (
            <p className="flex items-center gap-1 text-xs text-destructive">
              <HugeiconsIcon icon={Alert02Icon} className="size-3.5" />
              {errors.quantities}
            </p>
          ) : null}

          <FormActions>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>Huỷ</Button>
            <Button type="submit" disabled={isSubmitting}>{isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Tạo bản ghi"}</Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
