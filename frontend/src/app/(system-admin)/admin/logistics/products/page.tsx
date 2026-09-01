"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, ArrowLeft01Icon, Edit02Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { AppShell } from "@/components/shared/AppShell"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { FormActions } from "@/components/shared/form/Form"
import { formatDateTime } from "@/lib/format"
import { useProducts, useDeleteProduct } from "@/hooks/use-products"
import { ProductFormModal } from "@/components/system-admin/logistic/ProductFormModal"

type Product = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_product_dto.ProductResponse"]

const columns: Column<Product>[] = [
  { key: "sku", header: "Mã SKU", cell: (item) => <span className="font-semibold">{item.sku}</span> },
  { key: "name", header: "Tên sản phẩm", cell: (item) => item.name },
  { key: "category", header: "Danh mục", cell: (item) => item.category ?? "—" },
  {
    key: "price",
    header: "Giá",
    cell: (item) => (item.price != null ? `${item.price.toLocaleString("vi-VN")} ₫` : "—"),
    className: "text-right",
    headerClassName: "text-right",
  },
  {
    key: "weight_gram",
    header: "Cân nặng (g)",
    cell: (item) => item.weight_gram ?? "—",
    className: "text-right",
    headerClassName: "text-right",
  },
  {
    key: "updated_at",
    header: "Cập nhật",
    cell: (item) => formatDateTime(item.updated_at),
    className: "text-muted-foreground",
  },
]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export default function Page() {
  const [formOpen, setFormOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [productToDelete, setProductToDelete] = useState<Product | undefined>()
  const { data, isLoading, isError, error, refetch } = useProducts()

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách sản phẩm",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteProduct(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Sản phẩm ${productToDelete?.sku} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setProductToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!productToDelete?.id) return
    deleteMutation.mutate(productToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá sản phẩm này.",
          type: "error",
        })
      }
    })
  }

  if (isLoading) {
    return (
      <div className="flex flex-col gap-2">
        {Array.from({ length: 4 }).map((_, index) => (
          <Skeleton key={index} className="h-11 w-full rounded-2xl" />
        ))}
      </div>
    )
  }

  if (isError) {
    return (
      <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-border bg-card px-6 py-14 text-center">
        <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
        <div>
          <p className="font-medium">Không thể tải danh sách sản phẩm</p>
          <p className="mt-1 text-sm text-muted-foreground">{errorMessage(error)}</p>
        </div>
        <Button variant="outline" size="sm" onClick={() => void refetch()}>
          Thử lại
        </Button>
      </div>
    )
  }

  const items = data?.items ?? []

  const openCreateForm = () => {
    setSelectedProduct(undefined)
    setFormOpen(true)
  }

  const openEditForm = (product: Product) => {
    setSelectedProduct(product)
    setFormOpen(true)
  }

  const handleDeleteClick = (product: Product) => {
    setProductToDelete(product)
    setDeleteConfirmOpen(true)
  }

  const tableColumns: Column<Product>[] = [
    ...columns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (product) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa sản phẩm"
            onClick={() => openEditForm(product)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá sản phẩm"
            onClick={() => handleDeleteClick(product)}
            className="text-muted-foreground hover:text-destructive hover:bg-destructive/10"
          >
            <HugeiconsIcon icon={Delete01Icon} className="size-4" />
          </Button>
        </div>
      ),
    },
  ]

  return (
    <AppShell
      title="Sản phẩm"
      description="Danh sách sản phẩm trong hệ thống"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/admin/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="size-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={openCreateForm}>Thêm sản phẩm</Button>
        </div>
      }
    >
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? item.sku ?? ""}
        emptyText="Chưa có dữ liệu sản phẩm"
        emptyDescription="Bấm “Thêm sản phẩm” để tạo mới."
      />
      <ProductFormModal
        key={`${formOpen}-${selectedProduct?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        product={selectedProduct}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá sản phẩm</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn sản phẩm{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {productToDelete?.name}
              </span>{" "}
              (mã{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {productToDelete?.sku}
              </span>
              ) khỏi hệ thống?
            </DialogDescription>
          </DialogHeader>

          <FormActions>
            <Button variant="outline" onClick={() => setDeleteConfirmOpen(false)} disabled={isDeleting}>
              Không, quay lại
            </Button>
            <Button variant="destructive" onClick={confirmDelete} disabled={isDeleting}>
              {isDeleting ? "Đang xoá..." : "Xác nhận xoá"}
            </Button>
          </FormActions>
        </DialogContent>
      </Dialog>
    </AppShell>
  )
}
