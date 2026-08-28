"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, Add01Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import { useCreateInbound, useUpdateInbound, useInboundFormOptions } from "@/hooks/use-inbounds"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type Inbound = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"]
type CreateInboundRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.CreateInboundRequest"]
type UpdateInboundRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.UpdateInboundRequest"]
type InboundItemRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundItemRequest"]

type InboundFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  inbound?: Inbound
  onSuccess?: () => void
}

type FormItem = {
  productId: string
  expectedQty: string
  receivedQty: string
  rejectedQty: string
  qcPassed: string
}

const STATUS_OPTIONS = ["PENDING", "RECEIVING", "COMPLETED"]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function emptyItem(): FormItem {
  return { productId: "", expectedQty: "", receivedQty: "", rejectedQty: "", qcPassed: "" }
}

function toQty(value: string): number | undefined {
  if (value.trim() === "") return undefined
  return Math.max(0, Number(value) || 0)
}

export function InboundFormModal({
  open,
  onOpenChange,
  inbound,
  onSuccess,
}: InboundFormModalProps) {
  const isEdit = Boolean(inbound)
  const [receiptCode, setReceiptCode] = useState(() => {
    if (inbound) return inbound.receipt_code ?? ""
    return `PN-${new Date().getFullYear().toString().slice(-2)}${(new Date().getMonth() + 1).toString().padStart(2, "0")}-${Math.floor(1000 + Math.random() * 9000)}`
  })
  const [supplierName, setSupplierName] = useState(inbound?.supplier_name ?? "")
  const [warehouseId, setWarehouseId] = useState(inbound ? String(inbound.warehouse_id ?? "") : "")
  const [status, setStatus] = useState(inbound?.status ?? "PENDING")
  const [createdBy, setCreatedBy] = useState(inbound?.created_by ?? "")
  const [items, setItems] = useState<FormItem[]>(() => {
    if (inbound) {
      const mapped = (inbound.items ?? []).map((item) => ({
        productId: String(item.product_id ?? ""),
        expectedQty: item.expected_qty != null ? String(item.expected_qty) : "",
        receivedQty: item.received_qty != null ? String(item.received_qty) : "",
        rejectedQty: item.rejected_qty != null ? String(item.rejected_qty) : "",
        qcPassed: item.qc_passed != null ? String(item.qc_passed) : "",
      }))
      return mapped.length > 0 ? mapped : [emptyItem()]
    }
    return [emptyItem()]
  })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const { warehouses, products } = useInboundFormOptions(open)

  const createMutation = useCreateInbound(() => {
    toast.add({ title: "Đã tạo phiếu nhập", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateInbound(() => {
    toast.add({ title: "Đã cập nhật phiếu nhập", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const handleAddItem = () => {
    setItems((current) => [...current, emptyItem()])
  }

  const handleRemoveItem = (index: number) => {
    if (items.length <= 1) return
    setItems((current) => current.filter((_, itemIndex) => itemIndex !== index))
  }

  const handleItemChange = (index: number, field: keyof FormItem, value: string) => {
    setItems((current) =>
      current.map((item, itemIndex) => (itemIndex === index ? { ...item, [field]: value } : item))
    )
  }

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!isEdit && !receiptCode.trim()) nextErrors.receiptCode = "Mã phiếu không được bỏ trống."
    if (!isEdit && !supplierName.trim()) nextErrors.supplierName = "Tên nhà cung cấp không được bỏ trống."
    if (!isEdit && !warehouseId) nextErrors.warehouseId = "Vui lòng chọn kho nhập."
    const hasInvalidItem = items.some((item) => !item.productId)
    const hasNegative = items.some((item) =>
      [item.expectedQty, item.receivedQty, item.rejectedQty, item.qcPassed].some(
        (value) => value.trim() !== "" && Number(value) < 0
      )
    )
    if (hasInvalidItem) nextErrors.items = "Vui lòng chọn sản phẩm cho toàn bộ mặt hàng."
    if (hasNegative) nextErrors.items = "Số lượng không được là số âm."
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

    const itemPayload: InboundItemRequest[] = items.map((item) => ({
      product_id: Number(item.productId),
      expected_qty: toQty(item.expectedQty),
      received_qty: toQty(item.receivedQty),
      rejected_qty: toQty(item.rejectedQty),
      qc_passed: toQty(item.qcPassed),
    }))

    if (isEdit && inbound?.id) {
      const payload: UpdateInboundRequest = { status, items: itemPayload }
      updateMutation.mutate(
        { id: inbound.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật phiếu nhập", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    const payload: CreateInboundRequest = {
      receipt_code: receiptCode.trim(),
      supplier_name: supplierName.trim(),
      warehouse_id: Number(warehouseId),
      status,
      created_by: createdBy.trim() || undefined,
      items: itemPayload,
    }
    createMutation.mutate(payload, {
      onError: (error) =>
        toast.add({ title: "Không thể tạo phiếu nhập", description: errorMessage(error), type: "error" }),
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa phiếu nhập" : "Tạo phiếu nhập mới"}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? "Cập nhật trạng thái và danh sách mặt hàng của phiếu nhập."
              : "Nhập thông tin phiếu nhập hàng vào kho."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Mã phiếu" htmlFor="inbound-receipt-code" error={errors.receiptCode} required={!isEdit}>
              <Input
                id="inbound-receipt-code"
                value={receiptCode}
                onChange={(event) => setReceiptCode(event.target.value)}
                placeholder="VD: PN-2408-001"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
            <FormField label="Nhà cung cấp" htmlFor="inbound-supplier" error={errors.supplierName} required={!isEdit}>
              <Input
                id="inbound-supplier"
                value={supplierName}
                onChange={(event) => setSupplierName(event.target.value)}
                placeholder="Tên nhà cung cấp"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Kho nhập" htmlFor="inbound-warehouse" error={errors.warehouseId} required={!isEdit}>
              {isEdit ? (
                <Input id="inbound-warehouse" value={warehouseId} disabled />
              ) : (
                <select
                  id="inbound-warehouse"
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
              )}
            </FormField>
            <FormField label="Trạng thái" htmlFor="inbound-status">
              <select
                id="inbound-status"
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
          </div>

          {!isEdit && (
            <FormField label="Người tạo" htmlFor="inbound-created-by">
              <Input
                id="inbound-created-by"
                value={createdBy}
                onChange={(event) => setCreatedBy(event.target.value)}
                placeholder="Tên người tạo"
                disabled={isSubmitting}
              />
            </FormField>
          )}

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-semibold uppercase tracking-wider text-neutral-800 dark:text-neutral-200">
                Mặt hàng
              </h3>
              <Button type="button" variant="outline" size="sm" onClick={handleAddItem} disabled={isSubmitting} className="gap-1 text-xs">
                <HugeiconsIcon icon={Add01Icon} className="size-3" /> Thêm dòng
              </Button>
            </div>
            {errors.items ? (
              <p className="flex items-center gap-1.5 text-xs text-destructive">
                <HugeiconsIcon icon={Alert02Icon} className="size-3.5" />
                {errors.items}
              </p>
            ) : null}
            {items.map((item, index) => (
              <div key={index} className="flex flex-wrap items-end gap-2 rounded-2xl border border-border/60 bg-muted/10 p-3">
                <FormField label={`Sản phẩm ${index + 1}`} className="min-w-40 flex-1">
                  <select
                    value={item.productId}
                    onChange={(event) => handleItemChange(index, "productId", event.target.value)}
                    disabled={isSubmitting}
                    className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 text-sm outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
                  >
                    <option value="">-- Chọn sản phẩm --</option>
                    {(products.data?.items ?? []).map((product) => (
                      <option key={product.id} value={product.id}>
                        {product.name} ({product.sku})
                      </option>
                    ))}
                  </select>
                </FormField>
                <FormField label="Dự kiến" className="w-24">
                  <Input type="number" min="0" value={item.expectedQty} onChange={(event) => handleItemChange(index, "expectedQty", event.target.value)} disabled={isSubmitting} />
                </FormField>
                <FormField label="Nhận" className="w-24">
                  <Input type="number" min="0" value={item.receivedQty} onChange={(event) => handleItemChange(index, "receivedQty", event.target.value)} disabled={isSubmitting} />
                </FormField>
                <FormField label="Loại" className="w-24">
                  <Input type="number" min="0" value={item.rejectedQty} onChange={(event) => handleItemChange(index, "rejectedQty", event.target.value)} disabled={isSubmitting} />
                </FormField>
                <FormField label="QC đạt" className="w-24">
                  <Input type="number" min="0" value={item.qcPassed} onChange={(event) => handleItemChange(index, "qcPassed", event.target.value)} disabled={isSubmitting} />
                </FormField>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => handleRemoveItem(index)}
                  className="h-9 w-9 shrink-0 text-muted-foreground hover:text-destructive"
                  disabled={isSubmitting || items.length <= 1}
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
              {isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Tạo phiếu nhập"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
