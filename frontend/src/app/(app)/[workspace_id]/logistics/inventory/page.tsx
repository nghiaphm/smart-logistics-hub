"use client"

import { useEffect, useState } from "react"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, Edit02Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { FormActions } from "@/components/shared/form/Form"
import { useInventory, useDeleteInventory } from "@/hooks/use-inventory"
import { InventoryFormModal } from "./InventoryFormModal"

type Inventory = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inventory_dto.InventoryResponse"]

const columns: Column<Inventory>[] = [
  { key: "product_id", header: "Mã sản phẩm", cell: (item) => <span className="font-medium">{item.product_id}</span> },
  { key: "warehouse_id", header: "Kho", cell: (item) => item.warehouse_id },
  { key: "available_qty", header: "Tồn sẵn sàng", cell: (item) => item.available_qty ?? 0, className: "text-right" },
  { key: "reserved_qty", header: "Giữ chỗ", cell: (item) => item.reserved_qty ?? 0, className: "text-right" },
  { key: "damaged_qty", header: "Hư hỏng", cell: (item) => item.damaged_qty ?? 0, className: "text-right" },
  { key: "hold_qty", header: "Chờ duyệt", cell: (item) => item.hold_qty ?? 0, className: "text-right" },
  { key: "updated_at", header: "Cập nhật", cell: (item) => item.updated_at ?? "—", className: "text-muted-foreground" },
]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export default function Page() {
  const [formOpen, setFormOpen] = useState(false)
  const [selectedInventory, setSelectedInventory] = useState<Inventory | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [inventoryToDelete, setInventoryToDelete] = useState<Inventory | undefined>()
  const { data, isLoading, isError, error, refetch } = useInventory()

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách tồn kho",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteInventory(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Bản ghi tồn kho sản phẩm ${inventoryToDelete?.product_id} tại kho ${inventoryToDelete?.warehouse_id} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setInventoryToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!inventoryToDelete?.id) return
    deleteMutation.mutate(inventoryToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá bản ghi tồn kho này.",
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
          <p className="font-medium">Không thể tải danh sách tồn kho</p>
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
    setSelectedInventory(undefined)
    setFormOpen(true)
  }

  const openEditForm = (inventory: Inventory) => {
    setSelectedInventory(inventory)
    setFormOpen(true)
  }

  const handleDeleteClick = (inventory: Inventory) => {
    setInventoryToDelete(inventory)
    setDeleteConfirmOpen(true)
  }

  const tableColumns: Column<Inventory>[] = [
    ...columns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (inventory) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa tồn kho"
            onClick={() => openEditForm(inventory)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá tồn kho"
            onClick={() => handleDeleteClick(inventory)}
            className="text-muted-foreground hover:text-destructive hover:bg-destructive/10"
          >
            <HugeiconsIcon icon={Delete01Icon} className="size-4" />
          </Button>
        </div>
      ),
    },
  ]

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Tồn kho</h1>
          <p className="mt-1 text-sm text-muted-foreground">Danh sách tồn kho theo sản phẩm và kho</p>
        </div>
        <Button size="sm" onClick={openCreateForm}>Thêm tồn kho</Button>
      </div>
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? `${item.product_id}-${item.warehouse_id}`}
        loading={isLoading}
        emptyText="Chưa có dữ liệu tồn kho"
      />
      <InventoryFormModal
        key={`${formOpen}-${selectedInventory?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        inventory={selectedInventory}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá tồn kho</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn bản ghi tồn kho của sản phẩm{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {inventoryToDelete?.product_id}
              </span>{" "}
              tại kho{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {inventoryToDelete?.warehouse_id}
              </span>{" "}
              khỏi hệ thống?
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
    </div>
  )
}
