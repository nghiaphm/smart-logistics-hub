"use client"

import { useState } from "react"

import { AppShell } from "@/components/shared/AppShell"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { Form, FormField } from "@/components/shared/form/Form"
import { AppModalActions, AppModalShell } from "@/components/shared/modal"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

type Order = {
  code: string
  destination: string
  status: "đang xử lý" | "đang vận chuyển" | "đã giao"
  updatedAt: string
}

const sampleOrders: Order[] = [
  { code: "DH-2408-001", destination: "Kho VSIP 1 · Bình Dương", status: "đang vận chuyển", updatedAt: "14:20" },
  { code: "DH-2408-002", destination: "Kho Tân Thuận · TP.HCM", status: "đã giao", updatedAt: "13:05" },
  { code: "DH-2408-003", destination: "Kho Đà Nẵng", status: "đang xử lý", updatedAt: "11:48" },
]

const columns: Column<Order>[] = [
  { key: "code", header: "Mã đơn", cell: (order) => <span className="font-medium">{order.code}</span> },
  { key: "destination", header: "Điểm đến", cell: (order) => order.destination },
  {
    key: "status",
    header: "Trạng thái",
    cell: (order) => <Badge variant="outline">{order.status}</Badge>,
  },
  { key: "updatedAt", header: "Cập nhật", cell: (order) => order.updatedAt, className: "text-muted-foreground" },
]

export default function Page() {
  const [open, setOpen] = useState(false)
  const [code, setCode] = useState("")

  return (
    <AppShell
      title="Tổng quan"
      description="Các đơn hàng mới nhất trong phiên điều hành"
      actions={
        <Button onClick={() => setOpen(true)}>Thêm đơn hàng</Button>
      }
    >
      <DataTable columns={columns} rows={sampleOrders} rowKey={(order) => order.code} />

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
