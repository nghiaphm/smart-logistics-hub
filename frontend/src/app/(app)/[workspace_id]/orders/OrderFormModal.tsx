"use client"

import { useEffect, useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, Add01Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import { apiClient } from "@/lib/api-client"
import { useCreateOrder, useUpdateOrder } from "@/hooks/use-orders"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { Form, FormField, FormActions } from "@/components/shared/form/Form"
import type { components } from "@/types/api"

type OrderResponse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"]
type PaginatedWarehouses = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"]
type PaginatedProducts = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.PaginatedResponse"]

type FormItem = {
  productId: string
  productName: string
  quantity: number
  weightGram: number
}

interface OrderFormModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  order?: OrderResponse // If provided, we are editing
  onSuccess?: () => void
}

export function OrderFormModal({ open, onOpenChange, order, onSuccess }: OrderFormModalProps) {
  const isEdit = !!order

  // Form Fields
  const [orderCode, setOrderCode] = useState("")
  const [warehouseId, setWarehouseId] = useState("")
  const [senderName, setSenderName] = useState("")
  const [senderPhone, setSenderPhone] = useState("")
  const [senderAddress, setSenderAddress] = useState("")
  const [senderProvince, setSenderProvince] = useState("")
  const [senderDistrict, setSenderDistrict] = useState("")
  const [senderWard, setSenderWard] = useState("")
  const [senderPostalCode, setSenderPostalCode] = useState("")

  const [receiverName, setReceiverName] = useState("")
  const [receiverPhone, setReceiverPhone] = useState("")
  const [receiverAddress, setReceiverAddress] = useState("")
  const [receiverProvince, setReceiverProvince] = useState("")
  const [receiverDistrict, setReceiverDistrict] = useState("")
  const [receiverWard, setReceiverWard] = useState("")
  const [receiverPostalCode, setReceiverPostalCode] = useState("")

  const [items, setItems] = useState<FormItem[]>([])

  // Errors state
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({})

  // Fetch Warehouses
  const { data: warehousesData } = useQuery({
    queryKey: ["warehouses"],
    queryFn: () => apiClient<PaginatedWarehouses>("/warehouses?limit=100"),
    enabled: open,
  })

  // Fetch Products
  const { data: productsData } = useQuery({
    queryKey: ["products"],
    queryFn: () => apiClient<PaginatedProducts>("/products?limit=100"),
    enabled: open,
  })

  const warehouses = warehousesData?.items ?? []
  const products = productsData?.items ?? []

  // Initialize fields on open or change
  useEffect(() => {
    if (open) {
      setValidationErrors({})
      if (order) {
        setOrderCode(order.order_code ?? "")
        // Find warehouse_id associated with order; wait, backend order does not return warehouse_id directly? 
        // Let's check order fields: receiver_province, receiver_address etc. but wait! Does order have warehouse_id in GET /orders/:id?
        // Ah, order database only has assigned_driver_id, status etc., but let's default to the first available warehouse if not editable.
        setWarehouseId("") // Let them choose or default
        setSenderName(order.sender_name ?? "")
        setSenderPhone(order.sender_phone ?? "")
        setSenderAddress(order.sender_address ?? "")
        setSenderProvince(order.sender_province ?? "")
        setSenderDistrict(order.sender_district ?? "")
        setSenderWard(order.sender_ward ?? "")
        setSenderPostalCode(order.sender_postal_code ?? "")

        setReceiverName(order.receiver_name ?? "")
        setReceiverPhone(order.receiver_phone ?? "")
        setReceiverAddress(order.receiver_address ?? "")
        setReceiverProvince(order.receiver_province ?? "")
        setReceiverDistrict(order.receiver_district ?? "")
        setReceiverWard(order.receiver_ward ?? "")
        setReceiverPostalCode(order.receiver_postal_code ?? "")

        // Populate items
        const mappedItems = (order.items ?? []).map((it) => ({
          productId: String(it.product_id ?? ""),
          productName: it.product_name ?? "",
          quantity: it.quantity ?? 1,
          weightGram: it.weight_gram ?? 0,
        }))
        setItems(mappedItems.length > 0 ? mappedItems : [{ productId: "", productName: "", quantity: 1, weightGram: 0 }])
      } else {
        // Create new
        setOrderCode(`DH-${new Date().getFullYear().toString().slice(-2)}${(new Date().getMonth() + 1).toString().padStart(2, "0")}-${Math.floor(1000 + Math.random() * 9000)}`)
        setWarehouseId("")
        setSenderName("")
        setSenderPhone("")
        setSenderAddress("")
        setSenderProvince("")
        setSenderDistrict("")
        setSenderWard("")
        setSenderPostalCode("")

        setReceiverName("")
        setReceiverPhone("")
        setReceiverAddress("")
        setReceiverProvince("")
        setReceiverDistrict("")
        setReceiverWard("")
        setReceiverPostalCode("")

        setItems([{ productId: "", productName: "", quantity: 1, weightGram: 0 }])
      }
    }
  }, [open, order])

  // Populate first warehouse once loaded on creation
  useEffect(() => {
    if (open && !isEdit && warehouses.length > 0 && !warehouseId) {
      setWarehouseId(String(warehouses[0].id))
    }
  }, [open, isEdit, warehouses, warehouseId])

  const handleAddItem = () => {
    setItems((prev) => [...prev, { productId: "", productName: "", quantity: 1, weightGram: 0 }])
  }

  const handleRemoveItem = (index: number) => {
    if (items.length <= 1) return
    setItems((prev) => prev.filter((_, i) => i !== index))
  }

  const handleItemProductChange = (index: number, productIdStr: string) => {
    const selectedProd = products.find((p) => String(p.id) === productIdStr)
    setItems((prev) =>
      prev.map((item, i) =>
        i === index
          ? {
              ...item,
              productId: productIdStr,
              productName: selectedProd?.name ?? "",
              weightGram: selectedProd?.weight_gram ?? 0,
            }
          : item
      )
    )
  }

  const handleItemQtyChange = (index: number, qtyStr: string) => {
    const qty = Math.max(1, parseInt(qtyStr, 10) || 1)
    setItems((prev) => prev.map((item, i) => (i === index ? { ...item, quantity: qty } : item)))
  }

  const validate = (): boolean => {
    const errs: Record<string, string> = {}
    if (!orderCode.trim()) errs.orderCode = "Mã đơn hàng không được bỏ trống."
    if (!warehouseId && !isEdit) errs.warehouseId = "Vui lòng chọn kho bãi khởi tạo."
    
    if (!senderName.trim()) errs.senderName = "Tên người gửi không được bỏ trống."
    if (!senderPhone.trim()) errs.senderPhone = "SĐT người gửi không được bỏ trống."
    if (!senderAddress.trim()) errs.senderAddress = "Địa chỉ người gửi không được bỏ trống."

    if (!receiverName.trim()) errs.receiverName = "Tên người nhận không được bỏ trống."
    if (!receiverPhone.trim()) errs.receiverPhone = "SĐT người nhận không được bỏ trống."
    if (!receiverAddress.trim()) errs.receiverAddress = "Địa chỉ người nhận không được bỏ trống."

    if (!isEdit) {
      const hasInvalidItem = items.some((it) => !it.productId)
      if (hasInvalidItem) {
        errs.items = "Vui lòng chọn sản phẩm hợp lệ cho toàn bộ danh mục mặt hàng."
      }
    }

    setValidationErrors(errs)
    return Object.keys(errs).length === 0
  }

  const createMutation = useCreateOrder(() => {
    toast.add({
      title: "Thành công",
      description: "Khởi tạo đơn hàng thành công.",
      type: "success",
    })
    onOpenChange(false)
    if (onSuccess) onSuccess()
  })

  const updateMutation = useUpdateOrder(() => {
    toast.add({
      title: "Thành công",
      description: "Cập nhật thông tin đơn hàng thành công.",
      type: "success",
    })
    onOpenChange(false)
    if (onSuccess) onSuccess()
  })

  const isSubmitting = createMutation.isPending || updateMutation.isPending

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault()
    if (!validate()) {
      toast.add({
        title: "Thông tin không hợp lệ",
        description: "Vui lòng điền đầy đủ các thông tin bắt buộc.",
        type: "error",
      })
      return
    }

    if (isEdit) {
      updateMutation.mutate({
        id: order!.id!,
        payload: {
          order_code: orderCode,
          sender_name: senderName,
          sender_phone: senderPhone,
          sender_address: senderAddress,
          sender_province: senderProvince,
          sender_district: senderDistrict,
          sender_ward: senderWard,
          sender_postal_code: senderPostalCode,
          receiver_name: receiverName,
          receiver_phone: receiverPhone,
          receiver_address: receiverAddress,
          receiver_province: receiverProvince,
          receiver_district: receiverDistrict,
          receiver_ward: receiverWard,
          receiver_postal_code: receiverPostalCode,
        },
      }, {
        onError: (err: any) => {
          toast.add({
            title: "Cập nhật thất bại",
            description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi lưu thông tin đơn hàng.",
            type: "error",
          })
        }
      })
    } else {
      createMutation.mutate({
        order_code: orderCode,
        warehouse_id: parseInt(warehouseId, 10),
        sender_name: senderName,
        sender_phone: senderPhone,
        sender_address: senderAddress,
        sender_province: senderProvince,
        sender_district: senderDistrict,
        sender_ward: senderWard,
        sender_postal_code: senderPostalCode,
        receiver_name: receiverName,
        receiver_phone: receiverPhone,
        receiver_address: receiverAddress,
        receiver_province: receiverProvince,
        receiver_district: receiverDistrict,
        receiver_ward: receiverWard,
        receiver_postal_code: receiverPostalCode,
        items: items.map((it) => ({
          product_id: parseInt(it.productId, 10),
          product_name: it.productName,
          quantity: it.quantity,
          weight_gram: it.weightGram,
        })),
        status: "PENDING",
      }, {
        onError: (err: any) => {
          toast.add({
            title: "Khởi tạo thất bại",
            description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi lưu thông tin đơn hàng.",
            type: "error",
          })
        }
      })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Cập nhật đơn hàng" : "Khởi tạo đơn hàng mới"}</DialogTitle>
          <DialogDescription>
            {isEdit 
              ? "Chỉnh sửa thông tin hành trình và thông tin liên lạc của đơn hàng." 
              : "Hoàn thiện biểu mẫu thông tin dưới đây để tạo yêu cầu điều phối đơn hàng mới."
            }
          </DialogDescription>
        </DialogHeader>

        <Form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid gap-4 sm:grid-cols-2">
            <FormField label="Mã đơn hàng" error={validationErrors.orderCode} required>
              <Input
                value={orderCode}
                onChange={(e) => setOrderCode(e.target.value)}
                placeholder="VD: DH-2408-001"
                disabled={isSubmitting}
              />
            </FormField>

            {!isEdit && (
              <FormField label="Kho bãi khởi tạo" error={validationErrors.warehouseId} required>
                <select
                  value={warehouseId}
                  onChange={(e) => setWarehouseId(e.target.value)}
                  disabled={isSubmitting}
                  className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 py-1 text-sm outline-none transition-colors focus:border-ring focus:ring-1 focus:ring-ring"
                >
                  <option value="" disabled>-- Chọn kho bãi --</option>
                  {warehouses.map((w) => (
                    <option key={w.id} value={w.id}>{w.name} ({w.warehouse_code})</option>
                  ))}
                </select>
              </FormField>
            )}
          </div>

          <div className="grid gap-6 md:grid-cols-2">
            {/* Sender Information */}
            <div className="space-y-4 border border-border/60 rounded-2xl p-4 bg-muted/10">
              <h3 className="font-semibold text-sm text-neutral-800 dark:text-neutral-200 uppercase tracking-wider">Thông tin người gửi</h3>
              
              <FormField label="Họ và tên" error={validationErrors.senderName} required>
                <Input
                  value={senderName}
                  onChange={(e) => setSenderName(e.target.value)}
                  placeholder="Họ tên người gửi"
                  disabled={isSubmitting}
                />
              </FormField>

              <FormField label="Số điện thoại" error={validationErrors.senderPhone} required>
                <Input
                  value={senderPhone}
                  onChange={(e) => setSenderPhone(e.target.value)}
                  placeholder="SĐT người gửi"
                  disabled={isSubmitting}
                />
              </FormField>

              <FormField label="Địa chỉ cụ thể" error={validationErrors.senderAddress} required>
                <Input
                  value={senderAddress}
                  onChange={(e) => setSenderAddress(e.target.value)}
                  placeholder="Số nhà, tên đường..."
                  disabled={isSubmitting}
                />
              </FormField>

              <div className="grid gap-2 grid-cols-2">
                <FormField label="Phường/Xã">
                  <Input
                    value={senderWard}
                    onChange={(e) => setSenderWard(e.target.value)}
                    placeholder="Phường/Xã"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Quận/Huyện">
                  <Input
                    value={senderDistrict}
                    onChange={(e) => setSenderDistrict(e.target.value)}
                    placeholder="Quận/Huyện"
                    disabled={isSubmitting}
                  />
                </FormField>
              </div>

              <div className="grid gap-2 grid-cols-2">
                <FormField label="Tỉnh/Thành">
                  <Input
                    value={senderProvince}
                    onChange={(e) => setSenderProvince(e.target.value)}
                    placeholder="Tỉnh/Thành phố"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Mã bưu điện">
                  <Input
                    value={senderPostalCode}
                    onChange={(e) => setSenderPostalCode(e.target.value)}
                    placeholder="Mã bưu điện"
                    disabled={isSubmitting}
                  />
                </FormField>
              </div>
            </div>

            {/* Receiver Information */}
            <div className="space-y-4 border border-border/60 rounded-2xl p-4 bg-muted/10">
              <h3 className="font-semibold text-sm text-neutral-800 dark:text-neutral-200 uppercase tracking-wider">Thông tin người nhận</h3>

              <FormField label="Họ và tên" error={validationErrors.receiverName} required>
                <Input
                  value={receiverName}
                  onChange={(e) => setReceiverName(e.target.value)}
                  placeholder="Họ tên người nhận"
                  disabled={isSubmitting}
                />
              </FormField>

              <FormField label="Số điện thoại" error={validationErrors.receiverPhone} required>
                <Input
                  value={receiverPhone}
                  onChange={(e) => setReceiverPhone(e.target.value)}
                  placeholder="SĐT người nhận"
                  disabled={isSubmitting}
                />
              </FormField>

              <FormField label="Địa chỉ cụ thể" error={validationErrors.receiverAddress} required>
                <Input
                  value={receiverAddress}
                  onChange={(e) => setReceiverAddress(e.target.value)}
                  placeholder="Số nhà, tên đường..."
                  disabled={isSubmitting}
                />
              </FormField>

              <div className="grid gap-2 grid-cols-2">
                <FormField label="Phường/Xã">
                  <Input
                    value={receiverWard}
                    onChange={(e) => setReceiverWard(e.target.value)}
                    placeholder="Phường/Xã"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Quận/Huyện">
                  <Input
                    value={receiverDistrict}
                    onChange={(e) => setReceiverDistrict(e.target.value)}
                    placeholder="Quận/Huyện"
                    disabled={isSubmitting}
                  />
                </FormField>
              </div>

              <div className="grid gap-2 grid-cols-2">
                <FormField label="Tỉnh/Thành">
                  <Input
                    value={receiverProvince}
                    onChange={(e) => setReceiverProvince(e.target.value)}
                    placeholder="Tỉnh/Thành phố"
                    disabled={isSubmitting}
                  />
                </FormField>
                <FormField label="Mã bưu điện">
                  <Input
                    value={receiverPostalCode}
                    onChange={(e) => setReceiverPostalCode(e.target.value)}
                    placeholder="Mã bưu điện"
                    disabled={isSubmitting}
                  />
                </FormField>
              </div>
            </div>
          </div>

          {/* Order Items section (Only shown during Create, because editing order items in standard delivery flow isn't in scope unless asked) */}
          {!isEdit && (
            <div className="space-y-4 border border-border/60 rounded-2xl p-4 bg-muted/5">
              <div className="flex items-center justify-between">
                <h3 className="font-semibold text-sm text-neutral-800 dark:text-neutral-200 uppercase tracking-wider">Mặt hàng vận chuyển</h3>
                <Button 
                  type="button" 
                  variant="outline" 
                  size="sm" 
                  onClick={handleAddItem}
                  className="gap-1 text-xs"
                  disabled={isSubmitting}
                >
                  <HugeiconsIcon icon={Add01Icon} className="h-3 w-3" /> Thêm dòng
                </Button>
              </div>
              {validationErrors.items && (
                <div className="flex items-center gap-1.5 text-xs text-destructive">
                  <HugeiconsIcon icon={Alert02Icon} className="h-3.5 w-3.5" />
                  {validationErrors.items}
                </div>
              )}

              <div className="space-y-3">
                {items.map((item, index) => (
                  <div key={index} className="flex gap-2 items-end">
                    <div className="flex-1 min-w-0">
                      <label className="text-[11px] font-medium text-muted-foreground uppercase tracking-wider block mb-1">
                        Sản phẩm ({index + 1})
                      </label>
                      <select
                        value={item.productId}
                        onChange={(e) => handleItemProductChange(index, e.target.value)}
                        disabled={isSubmitting}
                        className="h-9 w-full rounded-4xl border border-input bg-input/30 px-3 py-1 text-sm outline-none transition-colors focus:border-ring focus:ring-1 focus:ring-ring"
                      >
                        <option value="">-- Chọn sản phẩm --</option>
                        {products.map((p) => (
                          <option key={p.id} value={p.id}>{p.name} ({p.sku})</option>
                        ))}
                      </select>
                    </div>

                    <div className="w-24">
                      <label className="text-[11px] font-medium text-muted-foreground uppercase tracking-wider block mb-1">
                        Số lượng
                      </label>
                      <Input
                        type="number"
                        min="1"
                        value={item.quantity}
                        onChange={(e) => handleItemQtyChange(index, e.target.value)}
                        disabled={isSubmitting}
                        className="text-right"
                      />
                    </div>

                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => handleRemoveItem(index)}
                      className="text-muted-foreground hover:text-destructive h-9 w-9 shrink-0"
                      disabled={isSubmitting || items.length <= 1}
                    >
                      <HugeiconsIcon icon={Delete01Icon} className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>
          )}

          <FormActions>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>
              Huỷ
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Đang xử lý..." : isEdit ? "Lưu thay đổi" : "Tạo đơn hàng"}
            </Button>
          </FormActions>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
