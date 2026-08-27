"use client"

import { useState } from "react"
import Link from "next/link"
import { useParams } from "next/navigation"

import { AppShell } from "@/components/shared/AppShell"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { Form, FormField } from "@/components/shared/form/Form"
import { AppModalActions, AppModalShell } from "@/components/shared/modal"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useOrders } from "@/hooks/use-orders"
import type { components } from "@/types/api"

type OrderResponse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_order_dto.OrderResponse"]

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

const columns: Column<OrderResponse>[] = [
  {
    key: "order_code",
    header: "Mã đơn",
    cell: (order) => <span className="font-semibold">{order.order_code}</span>,
  },
  {
    key: "receiver_name",
    header: "Người nhận",
    cell: (order) => <span className="font-medium text-sm">{order.receiver_name}</span>,
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
    header: "Cập nhật",
    cell: (order) => {
      if (!order.created_at) return "—"
      const date = new Date(order.created_at)
      return (
        <span className="text-xs text-muted-foreground">
          {date.toLocaleTimeString("vi-VN", { hour: "2-digit", minute: "2-digit" })}
        </span>
      )
    },
  },
]

export default function Page() {
  const [open, setOpen] = useState(false)
  const [code, setCode] = useState("")

  const params = useParams<{ workspace_id: string }>()
  const workspaceId = params.workspace_id

  const { data, isLoading } = useOrders(workspaceId, 3)

  const orders = data?.items ?? []

  return (
    <AppShell
      title="Tổng quan"
      description="Các đơn hàng mới nhất trong phiên điều hành"
      actions={
        <div className="flex items-center gap-2">
          <Link href={`/${workspaceId}/orders`}>
            <Button variant="outline" size="sm">
              Xem tất cả
            </Button>
          </Link>
          <Button size="sm" onClick={() => setOpen(true)}>Thêm đơn hàng</Button>
        </div>
      }
    >
      <DataTable
        columns={columns}
        rows={orders}
        rowKey={(order) => order.id ?? order.order_code ?? ""}
        loading={isLoading}
        emptyText="Chưa có đơn hàng nào được tạo."
      />

      <AppModalShell
        open={open}
        onOpenChange={setOpen}
        title="Thêm đơn hàng"
        description="Nhập thông tin đơn hàng để đưa vào tuyến điều phối"
        actions={
          <AppModalActions>
            <Button variant="outline" onClick={() => setOpen(false)}>
              Huỷ
            </Button>
            <Button onClick={() => setOpen(false)}>Lưu đơn hàng</Button>
          </AppModalActions>
        }
      >
        <Form onSubmit={(event) => event.preventDefault()}>
          <FormField label="Mã đơn hàng" htmlFor="order-code" required>
            <Input
              id="order-code"
              value={code}
              onChange={(event) => setCode(event.target.value)}
              placeholder="VD: DH-2408-004"
            />
          </FormField>
        </Form>
      </AppModalShell>
    </AppShell>
  )
}
