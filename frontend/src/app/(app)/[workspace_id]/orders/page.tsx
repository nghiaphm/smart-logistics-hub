"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import Link from "next/link"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, ArrowLeft01Icon, Edit02Icon, Delete01Icon } from "@hugeicons/core-free-icons"

import { useOrders, useDeleteOrder } from "@/hooks/use-orders"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { AppShell } from "@/components/shared/AppShell"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { OrderFormModal } from "./OrderFormModal"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog"
import { FormActions } from "@/components/shared/form/Form"
import type { components } from "@/types/api"

type OrderResponse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function getStatusBadge(status?: string) {
  switch (status?.toUpperCase()) {
    case "PENDING":
      return (
        <Badge variant="outline" className="border-amber-500/30 bg-amber-500/10 text-amber-600 dark:text-amber-400">
          Chờ xử lý
        </Badge>
      )
    case "ASSIGNED":
      return <Badge variant="secondary">Đã phân tuyến</Badge>
    case "SHIPPED":
      return <Badge variant="default">Đang vận chuyển</Badge>
    case "DELIVERED":
      return (
        <Badge variant="outline" className="border-emerald-500/30 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          Đã giao
        </Badge>
      )
    default:
      return <Badge variant="outline">{status || "Không rõ"}</Badge>
  }
}

export default function Page() {
  const params = useParams<{ workspace_id: string }>()
  const workspaceId = params.workspace_id
  const [modalOpen, setModalOpen] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<OrderResponse | undefined>(undefined)

  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false)
  const [orderToDelete, setOrderToDelete] = useState<OrderResponse | undefined>(undefined)

  const { data, isLoading, isError, error, refetch } = useOrders(workspaceId, 100)

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách đơn hàng",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  const orders = data?.items ?? []

  const handleEditClick = (order: OrderResponse) => {
    setSelectedOrder(order)
    setModalOpen(true)
  }

  const handleCreateClick = () => {
    setSelectedOrder(undefined)
    setModalOpen(true)
  }

  const handleDeleteClick = (order: OrderResponse) => {
    setOrderToDelete(order)
    setDeleteConfirmOpen(true)
  }

  const deleteMutation = useDeleteOrder(() => {
    toast.add({
      title: "Xóa thành công",
      description: `Đơn hàng ${orderToDelete?.order_code} đã được loại bỏ khỏi hệ thống.`,
      type: "success",
    })
    setDeleteConfirmOpen(false)
    setOrderToDelete(undefined)
  })

  const confirmDelete = async () => {
    if (!orderToDelete?.id) return
    deleteMutation.mutate(orderToDelete.id, {
      onError: (err: unknown) => {
        toast.add({
          title: "Xoá thất bại",
          description: err instanceof Error ? err.message : "Đã xảy ra lỗi khi xoá đơn hàng này.",
          type: "error",
        })
      }
    })
  }

  const isDeleting = deleteMutation.isPending

  const columns: Column<OrderResponse>[] = [
    {
      key: "order_code",
      header: "Mã đơn",
      cell: (order) => <span className="font-semibold">{order.order_code}</span>,
    },
    {
      key: "sender_name",
      header: "Người gửi",
      cell: (order) => (
        <div className="flex flex-col">
          <span className="font-medium text-sm">{order.sender_name}</span>
          <span className="text-xs text-muted-foreground">{order.sender_phone}</span>
        </div>
      ),
    },
    {
      key: "receiver_name",
      header: "Người nhận",
      cell: (order) => (
        <div className="flex flex-col">
          <span className="font-medium text-sm">{order.receiver_name}</span>
          <span className="text-xs text-muted-foreground">{order.receiver_phone}</span>
        </div>
      ),
    },
    {
      key: "receiver_address",
      header: "Điểm đến",
      cell: (order) => (
        <span className="text-sm text-neutral-600 dark:text-neutral-400 line-clamp-1">
          {order.receiver_address}, {order.receiver_province}
        </span>
      ),
    },
    {
      key: "status",
      header: "Trạng thái",
      cell: (order) => getStatusBadge(order.status),
    },
    {
      key: "created_at",
      header: "Ngày tạo",
      cell: (order) => {
        if (!order.created_at) return "—"
        const date = new Date(order.created_at)
        return (
          <span className="text-xs text-muted-foreground">
            {date.toLocaleDateString("vi-VN", {
              day: "2-digit",
              month: "2-digit",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit",
            })}
          </span>
        )
      },
    },
    {
      key: "actions",
      header: "Thao tác",
      className: "text-right",
      headerClassName: "text-right",
      cell: (order) => (
        <div className="flex justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => handleEditClick(order)}
            title="Chỉnh sửa thông tin đơn hàng"
          >
            <HugeiconsIcon icon={Edit02Icon} className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => handleDeleteClick(order)}
            className="text-muted-foreground hover:text-destructive hover:bg-destructive/10"
            title="Xoá đơn hàng"
          >
            <HugeiconsIcon icon={Delete01Icon} className="h-4 w-4" />
          </Button>
        </div>
      ),
    },
  ]

  return (
    <AppShell
      title="Danh sách đơn hàng"
      description="Quản lý thông tin đơn hàng và theo dõi tuyến kết nối kho bãi"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/${workspaceId}/logistics`}>
            <Button variant="outline" size="sm" className="gap-1">
              <HugeiconsIcon icon={ArrowLeft01Icon} className="h-4 w-4" /> Quay lại
            </Button>
          </Link>
          <Button size="sm" onClick={handleCreateClick}>Thêm đơn hàng</Button>
        </div>
      }
    >
      {isError ? (
        <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-border bg-card px-6 py-14 text-center">
          <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
          <div>
            <p className="font-medium">Không thể tải danh sách đơn hàng</p>
            <p className="mt-1 text-sm text-muted-foreground">{errorMessage(error)}</p>
          </div>
          <Button variant="outline" size="sm" onClick={() => void refetch()}>
            Thử lại
          </Button>
        </div>
      ) : (
        <DataTable
          columns={columns}
          rows={orders}
          rowKey={(order) => order.id ?? order.order_code ?? ""}
          loading={isLoading}
          emptyText="Chưa có đơn hàng nào được tạo."
        />
      )}

      <OrderFormModal
        key={`${modalOpen}-${selectedOrder?.id ?? "new"}`}
        open={modalOpen}
        onOpenChange={setModalOpen}
        order={selectedOrder}
        onSuccess={() => void refetch()}
      />

      <Dialog open={deleteConfirmOpen} onOpenChange={setDeleteConfirmOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Xác nhận xoá đơn hàng</DialogTitle>
            <DialogDescription>
              Hành động này không thể hoàn tác. Bạn có chắc chắn muốn xoá vĩnh viễn đơn hàng{" "}
              <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                {orderToDelete?.order_code}
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
