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
import { useTracking, useDeleteTrackingEvent } from "@/hooks/use-tracking"
import { useOrders } from "@/hooks/use-orders"
import { TrackingFormModal } from "./TrackingFormModal"

type TrackingEvent = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_tracking_dto.TrackingEventResponse"]
type OrderResponse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"]

const baseColumns: Column<TrackingEvent>[] = [
  { key: "order_code", header: "Mã đơn", cell: (item) => <span className="font-semibold">{item.order_code}</span> },
  { key: "driver_code", header: "Tài xế", cell: (item) => item.driver_code ?? "—" },
  {
    key: "status_update",
    header: "Trạng thái",
    cell: (item) => (item.status_update ? <Badge variant="outline">{item.status_update}</Badge> : "—"),
  },
  { key: "note", header: "Ghi chú", cell: (item) => item.note ?? "—" },
  {
    key: "timestamp",
    header: "Thời điểm",
    cell: (item) => formatDateTime(item.timestamp),
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
  const [selectedEvent, setSelectedEvent] = useState<TrackingEvent | undefined>()
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [eventToDelete, setEventToDelete] = useState<TrackingEvent | undefined>()
  const { data, isLoading, isError, error, refetch } = useTracking()
  const { data: ordersData } = useOrders(workspaceId, 100)

  // Join theo order_id (ưu tiên) — fallback order_code cho tracking cũ chưa có link (xem WN-043).
  const ordersById = new Map<number, OrderResponse>()
  const ordersByCode = new Map<string, OrderResponse>()
  for (const order of ordersData?.items ?? []) {
    if (order.id != null) ordersById.set(order.id, order)
    if (order.order_code) ordersByCode.set(order.order_code, order)
  }
  const findOrder = (event: TrackingEvent): OrderResponse | undefined => {
    if (event.order_id != null) return ordersById.get(event.order_id)
    if (event.order_code) return ordersByCode.get(event.order_code)
    return undefined
  }

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách theo dõi",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const deleteMutation = useDeleteTrackingEvent(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Sự kiện theo dõi đơn ${eventToDelete?.order_code} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setEventToDelete(undefined)
  })

  const isDeleting = deleteMutation.isPending

  const confirmDelete = async () => {
    if (!eventToDelete?.id) return
    deleteMutation.mutate(eventToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá sự kiện theo dõi này.",
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
          <p className="font-medium">Không thể tải danh sách theo dõi</p>
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
    setSelectedEvent(undefined)
    setFormOpen(true)
  }

  const openEditForm = (event: TrackingEvent) => {
    setSelectedEvent(event)
    setFormOpen(true)
  }

  const handleDeleteClick = (event: TrackingEvent) => {
    setEventToDelete(event)
    setDeleteConfirmOpen(true)
  }

  const orderColumns: Column<TrackingEvent>[] = [
    {
      key: "order_sender",
      header: "Người gửi",
      cell: (event) => {
        const order = findOrder(event)
        if (!order) return "—"
        return (
          <div className="flex flex-col">
            <span className="font-medium text-sm">{order.sender_name}</span>
            <span className="text-xs text-muted-foreground">{order.sender_phone}</span>
          </div>
        )
      },
    },
    {
      key: "order_receiver",
      header: "Người nhận",
      cell: (event) => {
        const order = findOrder(event)
        if (!order) return "—"
        return (
          <div className="flex flex-col">
            <span className="font-medium text-sm">{order.receiver_name}</span>
            <span className="text-xs text-muted-foreground">{order.receiver_phone}</span>
          </div>
        )
      },
    },
    {
      key: "order_destination",
      header: "Điểm đến",
      cell: (event) => {
        const order = findOrder(event)
        if (!order) return "—"
        const destination = [order.receiver_address, order.receiver_province].filter(Boolean).join(", ")
        return (
          <span className="text-sm text-neutral-600 dark:text-neutral-400 line-clamp-1">
            {destination || "—"}
          </span>
        )
      },
    },
  ]

  const tableColumns: Column<TrackingEvent>[] = [
    ...baseColumns,
    ...orderColumns,
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (event) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            title="Sửa sự kiện theo dõi"
            onClick={() => openEditForm(event)}
          >
            <HugeiconsIcon icon={Edit02Icon} className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            title="Xoá sự kiện theo dõi"
            onClick={() => handleDeleteClick(event)}
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
      title="Theo dõi đơn"
      description="Nhật ký theo dõi vận chuyển theo đơn"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/${workspaceId}/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="size-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={openCreateForm}>Ghi nhận sự kiện</Button>
        </div>
      }
    >
      <DataTable
        columns={tableColumns}
        rows={items}
        rowKey={(item) => item.id ?? `${item.order_code}-${item.timestamp ?? ""}`}
        emptyText="Chưa có dữ liệu theo dõi"
        emptyDescription="Ghi nhận sự kiện khi có hoạt động vận chuyển."
      />
      <TrackingFormModal
        key={`${formOpen}-${selectedEvent?.id ?? "new"}`}
        open={formOpen}
        onOpenChange={setFormOpen}
        event={selectedEvent}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá sự kiện theo dõi</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn sự kiện theo dõi của đơn{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {eventToDelete?.order_code}
              </span>{" "}
              (trạng thái{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {eventToDelete?.status_update}
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
