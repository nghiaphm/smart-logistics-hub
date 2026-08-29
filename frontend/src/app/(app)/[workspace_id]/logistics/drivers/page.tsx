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
import { useDrivers, useDeleteDriver } from "@/hooks/use-drivers"
import { DriverFormModal } from "./DriverFormModal"

type Driver = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_driver_dto.DriverResponse"]

function getStatusBadge(status?: string) {
  switch (status?.toUpperCase()) {
    case "AVAILABLE":
      return (
        <Badge variant="outline" className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          Sẵn sàng
        </Badge>
      )
    case "BUSY":
      return <Badge variant="default">Đang làm việc</Badge>
    case "OFFLINE":
      return <Badge variant="secondary">Nghỉ</Badge>
    default:
      return <Badge variant="outline">{status || "—"}</Badge>
  }
}

const columns: Column<Driver>[] = [
  {
    key: "driver_code",
    header: "Mã tài xế",
    cell: (item) => <span className="font-semibold">{item.driver_code}</span>,
  },
  {
    key: "full_name",
    header: "Họ tên",
    cell: (item) => (
      <div className="flex flex-col">
        <span className="font-medium text-sm">{item.full_name}</span>
        <span className="text-xs text-muted-foreground">{item.phone}</span>
      </div>
    ),
  },
  { key: "vehicle_type", header: "Loại xe", cell: (item) => item.vehicle_type || "—" },
  { key: "license_plate", header: "Biển số", cell: (item) => item.license_plate || "—" },
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
  const [selectedDriver, setSelectedDriver] = useState<Driver | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [driverToDelete, setDriverToDelete] = useState<Driver | undefined>()
  const { data, isLoading, isError, error, refetch } = useDrivers()

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách tài xế",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteDriver(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Tài xế ${driverToDelete?.driver_code} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setDriverToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!driverToDelete?.id) return
    deleteMutation.mutate(driverToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá tài xế này.",
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
          <p className="font-medium">Không thể tải danh sách tài xế</p>
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
    setSelectedDriver(undefined)
    setFormOpen(true)
  }

  const openEditForm = (driver: Driver) => {
    setSelectedDriver(driver)
    setFormOpen(true)
  }

  const handleDeleteClick = (driver: Driver) => {
    setDriverToDelete(driver)
    setDeleteConfirmOpen(true)
  }

  const tableColumns: Column<Driver>[] = [
    ...columns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (driver) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa tài xế"
            onClick={() => openEditForm(driver)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá tài xế"
            onClick={() => handleDeleteClick(driver)}
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
      title="Tài xế"
      description="Danh sách tài xế vận tải"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/${workspaceId}/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="size-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={openCreateForm}>Thêm tài xế</Button>
        </div>
      }
    >
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? item.driver_code ?? ""}
        emptyText="Chưa có dữ liệu tài xế"
        emptyDescription="Bấm “Thêm tài xế” để đăng ký tài xế mới."
      />
      <DriverFormModal
        key={`${formOpen}-${selectedDriver?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        driver={selectedDriver}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá tài xế</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn tài xế{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {driverToDelete?.driver_code}
              </span>{" "}
              (họ tên{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {driverToDelete?.full_name}
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
