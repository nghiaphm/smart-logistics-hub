"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import { useCreateProduct, useUpdateProduct } from "@/hooks/use-products"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

type Product = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"]
type UpdateProductRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.UpdateProductRequest"]

type ProductFormModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  product?: Product
  onSuccess?: () => void
}

type DimensionFields = {
  lengthCm: string
  widthCm: string
  heightCm: string
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function parseDimension(value: string): number | undefined {
  if (!value.trim()) return undefined
  return Math.max(0, Number(value) || 0)
}

export function ProductFormModal({
  open,
  onOpenChange,
  product,
  onSuccess,
}: ProductFormModalProps) {
  const isEdit = Boolean(product)
  const [sku, setSku] = useState(product?.sku ?? "")
  const [name, setName] = useState(product?.name ?? "")
  const [category, setCategory] = useState(product?.category ?? "")
  const [price, setPrice] = useState(product?.price != null ? String(product.price) : "")
  const [weightGram, setWeightGram] = useState(product?.weight_gram != null ? String(product.weight_gram) : "")
  const [dimensions, setDimensions] = useState<DimensionFields>({
    lengthCm: product?.length_cm != null ? String(product.length_cm) : "",
    widthCm: product?.width_cm != null ? String(product.width_cm) : "",
    heightCm: product?.height_cm != null ? String(product.height_cm) : "",
  })
  const [createdBy, setCreatedBy] = useState(product?.created_by ?? "")
  const [errors, setErrors] = useState<Record<string, string>>({})

  const createMutation = useCreateProduct(() => {
    toast.add({ title: "Đã tạo sản phẩm", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const updateMutation = useUpdateProduct(() => {
    toast.add({ title: "Đã cập nhật sản phẩm", type: "success", timeout: 4000 })
    onOpenChange(false)
    onSuccess?.()
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const validate = () => {
    const nextErrors: Record<string, string> = {}
    if (!isEdit && !sku.trim()) nextErrors.sku = "Mã SKU không được bỏ trống."
    if (!name.trim()) nextErrors.name = "Tên sản phẩm không được bỏ trống."
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

    const payload: UpdateProductRequest = {
      name: name.trim(),
      category: category.trim() || undefined,
      price: price.trim() ? Number(price) : undefined,
      weight_gram: weightGram.trim() ? Math.max(0, Number.parseInt(weightGram, 10) || 0) : undefined,
      dimensions:
        dimensions.lengthCm || dimensions.widthCm || dimensions.heightCm
          ? {
              length: parseDimension(dimensions.lengthCm),
              width: parseDimension(dimensions.widthCm),
              height: parseDimension(dimensions.heightCm),
            }
          : undefined,
      created_by: createdBy.trim() || undefined,
    }

    if (isEdit && product?.id) {
      updateMutation.mutate(
        { id: product.id, payload },
        {
          onError: (error) =>
            toast.add({ title: "Không thể cập nhật sản phẩm", description: errorMessage(error), type: "error" }),
        }
      )
      return
    }

    createMutation.mutate(
      { ...payload, sku: sku.trim(), name: name.trim() },
      {
        onError: (error) =>
          toast.add({ title: "Không thể tạo sản phẩm", description: errorMessage(error), type: "error" }),
      }
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Sửa sản phẩm" : "Tạo sản phẩm mới"}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? "Cập nhật thông tin sản phẩm hiện tại."
              : "Nhập thông tin sản phẩm để thêm vào hệ thống."}
          </DialogDescription>
        </DialogHeader>
        <Form onSubmit={handleSubmit}>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Mã SKU" htmlFor="product-sku" error={errors.sku} required={!isEdit}>
              <Input
                id="product-sku"
                value={sku}
                onChange={(event) => setSku(event.target.value)}
                placeholder="VD: SP-0001"
                disabled={isSubmitting || isEdit}
              />
            </FormField>
            <FormField label="Tên sản phẩm" htmlFor="product-name" error={errors.name} required>
              <Input
                id="product-name"
                value={name}
                onChange={(event) => setName(event.target.value)}
                placeholder="Tên sản phẩm"
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <div className="grid gap-4 sm:grid-cols-3">
            <FormField label="Danh mục" htmlFor="product-category">
              <Input
                id="product-category"
                value={category}
                onChange={(event) => setCategory(event.target.value)}
                placeholder="VD: Điện tử"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Giá (₫)" htmlFor="product-price">
              <Input
                id="product-price"
                type="number"
                min="0"
                value={price}
                onChange={(event) => setPrice(event.target.value)}
                placeholder="0"
                disabled={isSubmitting}
              />
            </FormField>
            <FormField label="Cân nặng (g)" htmlFor="product-weight">
              <Input
                id="product-weight"
                type="number"
                min="0"
                value={weightGram}
                onChange={(event) => setWeightGram(event.target.value)}
                placeholder="0"
                disabled={isSubmitting}
              />
            </FormField>
          </div>

          <FormField label="Kích thước (cm)" hint="Có thể bỏ trống nếu chưa có.">
            <div className="grid gap-4 sm:grid-cols-3">
              <Input
                type="number"
                min="0"
                value={dimensions.lengthCm}
                onChange={(event) => setDimensions((current) => ({ ...current, lengthCm: event.target.value }))}
                placeholder="Dài"
                disabled={isSubmitting}
              />
              <Input
                type="number"
                min="0"
                value={dimensions.widthCm}
                onChange={(event) => setDimensions((current) => ({ ...current, widthCm: event.target.value }))}
                placeholder="Rộng"
                disabled={isSubmitting}
              />
              <Input
                type="number"
                min="0"
                value={dimensions.heightCm}
                onChange={(event) => setDimensions((current) => ({ ...current, heightCm: event.target.value }))}
                placeholder="Cao"
                disabled={isSubmitting}
              />
            </div>
          </FormField>

          <FormField label="Người tạo" htmlFor="product-created-by">
            <Input
              id="product-created-by"
              value={createdBy}
              onChange={(event) => setCreatedBy(event.target.value)}
              placeholder="Tên người tạo"
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
              {isSubmitting ? "Đang lưu..." : isEdit ? "Lưu thay đổi" : "Tạo sản phẩm"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
