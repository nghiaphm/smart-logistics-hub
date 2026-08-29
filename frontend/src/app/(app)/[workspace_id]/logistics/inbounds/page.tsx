"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, ArrowLeft01Icon, Edit02Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { AppShell } from "@/components/shared/AppShell"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { FormActions } from "@/components/shared/form/Form"
import { formatDateTime } from "@/lib/format"
import { useInbounds, useDeleteInbound } from "@/hooks/use-inbounds"
import { InboundFormModal } from "./InboundFormModal"

type Inbound = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_inbound_dto.InboundResponse"]

function getStatusBadge(status?: string) {
  switch (status?.toUpperCase()) {
    case "PENDING":
      return (
        <Badge variant="outline" className="border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400">
          Chờ duyệt
        </Badge>
      )
    case "RECEIVING":
      return <Badge variant="secondary">Đang nhập</Badge>
    case "COMPLETED":
      return (
        <Badge variant="outline" className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          Hoàn thành
        </Badge>
      )
    default:
      return <Badge variant="outline">{status || "—"}</Badge>
  }
}

const columns: Column<Inbound>[] = [
  { key: "receipt_code", header: "Mã phiếu", cell: (item) => <span className="font-semibold">{item.receipt_code}</span> },
  { key: "supplier_name", header: "Nhà cung cấp", cell: (item) => item.supplier_name },
  { key: "warehouse_id", header: "Kho", cell: (item) => item.warehouse_id },
  { key: "status", header: "Trạng thái", cell: (item) => getStatusBadge(item.status) },
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
  const params = useParams<{ workspace_id: string }>()
  const workspaceId = params.workspace_id
  const [formOpen, setFormOpen] = useState(false)
  const [selectedInbound, setSelectedInbound] = useState<Inbound | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [inboundToDelete, setInboundToDelete] = useState<Inbound | undefined>()
  const { data, isLoading, isError, error, refetch } = useInbounds()

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách hàng nhập",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteInbound(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Phiếu nhập ${inboundToDelete?.receipt_code} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setInboundToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!inboundToDelete?.id) return
    deleteMutation.mutate(inboundToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá phiếu nhập này.",
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
          <p className="font-medium">Không thể tải danh sách hàng nhập</p>
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
    setSelectedInbound(undefined)
    setFormOpen(true)
  }

  const openEditForm = (inbound: Inbound) => {
    setSelectedInbound(inbound)
    setFormOpen(true)
  }

  const handleDeleteClick = (inbound: Inbound) => {
    setInboundToDelete(inbound)
    setDeleteConfirmOpen(true)
  }

  const tableColumns: Column<Inbound>[] = [
    ...columns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (inbound) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa phiếu nhập"
            onClick={() => openEditForm(inbound)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá phiếu nhập"
            onClick={() => handleDeleteClick(inbound)}
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
      title="Hàng nhập"
      description="Danh sách phiếu nhập hàng vào kho"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/${workspaceId}/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="size-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={openCreateForm}>Thêm phiếu nhập</Button>
        </div>
      }
    >
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? item.receipt_code ?? ""}
        emptyText="Chưa có dữ liệu hàng nhập"
        emptyDescription="Bấm “Thêm phiếu nhập” để tạo mới."
      />
      <InboundFormModal
        key={`${formOpen}-${selectedInbound?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        inbound={selectedInbound}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá phiếu nhập</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn phiếu nhập{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {inboundToDelete?.receipt_code}
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
    </AppShell>
  )
}
