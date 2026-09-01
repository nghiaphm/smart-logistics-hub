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
import { useTrips, useDeleteTrip } from "@/hooks/use-trips"
import { TripFormModal } from "@/components/system-admin/logistic/TripFormModal"

type Trip = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_trip_dto.TripResponse"]

function getStatusBadge(status?: string) {
  switch (status?.toUpperCase()) {
    case "PLANNED":
      return <Badge variant="secondary">Đã lên kế hoạch</Badge>
    case "IN_TRANSIT":
      return <Badge variant="default">Đang chạy</Badge>
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

const columns: Column<Trip>[] = [
  { key: "trip_code", header: "Mã chuyến", cell: (item) => <span className="font-semibold">{item.trip_code}</span> },
  { key: "driver_id", header: "Tài xế", cell: (item) => item.driver_id ?? "—" },
  { key: "vehicle_license_plate", header: "Biển số", cell: (item) => item.vehicle_license_plate ?? "—" },
  { key: "status", header: "Trạng thái", cell: (item) => getStatusBadge(item.status) },
  {
    key: "total_distance_km",
    header: "Quãng đường (km)",
    cell: (item) => item.total_distance_km ?? "—",
    className: "text-right",
    headerClassName: "text-right",
  },
  {
    key: "estimated_duration_min",
    header: "Thời lượng (phút)",
    cell: (item) => item.estimated_duration_min ?? "—",
    className: "text-right",
    headerClassName: "text-right",
  },
]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export default function Page() {
  const [formOpen, setFormOpen] = useState(false)
  const [selectedTrip, setSelectedTrip] = useState<Trip | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [tripToDelete, setTripToDelete] = useState<Trip | undefined>()
  const { data, isLoading, isError, error, refetch } = useTrips()

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách chuyến xe",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteTrip(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Chuyến xe ${tripToDelete?.trip_code} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setTripToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!tripToDelete?.id) return
    deleteMutation.mutate(tripToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá chuyến xe này.",
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
          <p className="font-medium">Không thể tải danh sách chuyến xe</p>
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
    setSelectedTrip(undefined)
    setFormOpen(true)
  }

  const openEditForm = (trip: Trip) => {
    setSelectedTrip(trip)
    setFormOpen(true)
  }

  const handleDeleteClick = (trip: Trip) => {
    setTripToDelete(trip)
    setDeleteConfirmOpen(true)
  }

  const tableColumns: Column<Trip>[] = [
    ...columns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (trip) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa chuyến xe"
            onClick={() => openEditForm(trip)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá chuyến xe"
            onClick={() => handleDeleteClick(trip)}
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
      title="Chuyến xe"
      description="Danh sách chuyến vận chuyển"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/admin/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="size-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={openCreateForm}>Tạo chuyến xe</Button>
        </div>
      }
    >
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? item.trip_code ?? ""}
        emptyText="Chưa có dữ liệu chuyến xe"
        emptyDescription="Bấm “Tạo chuyến xe” để bắt đầu."
      />
      <TripFormModal
        key={`${formOpen}-${selectedTrip?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        trip={selectedTrip}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá chuyến xe</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn chuyến xe{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {tripToDelete?.trip_code}
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
