"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
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
import { useVehicles, useDeleteVehicle } from "@/hooks/use-vehicles"
import { VehicleFormModal } from "@/components/system-admin/logistic/VehicleFormModal"

type Vehicle = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_vehicle_dto.VehicleResponse"]

function getStatusBadge(status?: string) {
  switch (status?.toUpperCase()) {
    case "ACTIVE":
      return (
        <Badge variant="outline" className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          Đang hoạt động
        </Badge>
      )
    case "MAINTENANCE":
      return (
        <Badge variant="outline" className="border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400">
          Đang bảo trì
        </Badge>
      )
    case "INACTIVE":
      return <Badge variant="secondary">Ngừng hoạt động</Badge>
    default:
      return <Badge variant="outline">{status || "—"}</Badge>
  }
}

const columns: Column<Vehicle>[] = [
  {
    key: "license_plate",
    header: "Biển số",
    cell: (item) => <span className="font-semibold">{item.license_plate}</span>,
  },
  { key: "type", header: "Loại xe", cell: (item) => item.type || "—" },
  {
    key: "capacity",
    header: "Tải trọng (kg)",
    cell: (item) => (item.capacity != null ? item.capacity.toLocaleString("vi-VN") : "—"),
    className: "text-right",
    headerClassName: "text-right",
  },
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
  const [formOpen, setFormOpen] = useState(false)
  const [selectedVehicle, setSelectedVehicle] = useState<Vehicle | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [vehicleToDelete, setVehicleToDelete] = useState<Vehicle | undefined>()
  const { data, isLoading, isError, error, refetch } = useVehicles()

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách phương tiện",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteVehicle(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Phương tiện ${vehicleToDelete?.license_plate} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setVehicleToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!vehicleToDelete?.id) return
    deleteMutation.mutate(vehicleToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá phương tiện này.",
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
          <p className="font-medium">Không thể tải danh sách phương tiện</p>
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
    setSelectedVehicle(undefined)
    setFormOpen(true)
  }

  const openEditForm = (vehicle: Vehicle) => {
    setSelectedVehicle(vehicle)
    setFormOpen(true)
  }

  const handleDeleteClick = (vehicle: Vehicle) => {
    setVehicleToDelete(vehicle)
    setDeleteConfirmOpen(true)
  }

  const tableColumns: Column<Vehicle>[] = [
    ...columns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (vehicle) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa phương tiện"
            onClick={() => openEditForm(vehicle)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá phương tiện"
            onClick={() => handleDeleteClick(vehicle)}
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
      title="Phương tiện"
      description="Danh sách phương tiện vận tải (Fleet)"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/admin/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="size-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={openCreateForm}>Thêm phương tiện</Button>
        </div>
      }
    >
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? item.license_plate ?? ""}
        emptyText="Chưa có dữ liệu phương tiện"
        emptyDescription="Bấm “Thêm phương tiện” để đăng ký xe mới."
      />
      <VehicleFormModal
        key={`${formOpen}-${selectedVehicle?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        vehicle={selectedVehicle}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá phương tiện</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn phương tiện{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {vehicleToDelete?.license_plate}
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
