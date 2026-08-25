"use client"

import { useEffect } from "react"
import { useQueries } from "@tanstack/react-query"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  Alert02Icon,
  BoxesIcon,
  InvoiceIcon,
  PackageAddIcon,
  TruckIcon,
  UserMultipleIcon,
  WarehouseIcon,
} from "@hugeicons/core-free-icons"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { AppShell } from "@/components/shared/AppShell"

type CountResponse = Pick<
  components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"],
  "total"
>

type Metric = {
  label: string
  path: string
  icon: typeof WarehouseIcon
}

const metrics: Metric[] = [
  { label: "Kho bãi", path: "/warehouses", icon: WarehouseIcon },
  { label: "Sản phẩm", path: "/products", icon: BoxesIcon },
  { label: "Đơn hàng", path: "/orders", icon: InvoiceIcon },
  { label: "Chuyến xe", path: "/trips", icon: TruckIcon },
  { label: "Tài xế", path: "/drivers", icon: UserMultipleIcon },
  { label: "Phiếu nhập", path: "/inbounds", icon: PackageAddIcon },
]

export default function AdminDashboardPage() {
  const results = useQueries({
    queries: metrics.map((metric) => ({
      queryKey: ["admin", "metrics", metric.path],
      queryFn: () => apiClient<CountResponse>(`${metric.path}?limit=1`),
    })),
  })

  const hasError = results.some((result) => result.isError)
  const firstError = results.find((result) => result.isError)?.error ?? null

  useEffect(() => {
    if (!firstError) {
      return
    }
    toast.add({
      title: "Không thể tải số liệu tổng quan",
      description:
        firstError instanceof Error ? firstError.message : "Đã có lỗi xảy ra. Vui lòng thử lại.",
      type: "error",
      timeout: 6000,
    })
  }, [firstError])

  return (
    <AppShell
      title="Tổng quan hệ thống"
      description="Số liệu tổng hợp trên toàn hệ thống"
    >
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 xl:grid-cols-6">
        {metrics.map((metric, index) => {
          const result = results[index]
          return (
            <div key={metric.path} className="rounded-2xl border border-border bg-card p-5">
              <span className="flex size-9 items-center justify-center rounded-lg bg-admin-accent/10 text-admin-accent">
                <HugeiconsIcon icon={metric.icon} className="size-4.5" />
              </span>
              <p className="mt-4 text-[13px] font-medium text-muted-foreground">{metric.label}</p>
              <p className="mt-0.5 font-mono text-3xl font-semibold tabular-nums tracking-tight text-foreground">
                {result.isLoading ? (
                  <Skeleton className="h-8 w-14 rounded-lg" />
                ) : result.isError ? (
                  "—"
                ) : (
                  (result.data?.total ?? 0).toLocaleString("vi-VN")
                )}
              </p>
            </div>
          )
        })}
      </div>
      {hasError ? (
        <div className="flex items-center justify-between gap-3 rounded-2xl border border-border bg-card px-5 py-4">
          <div className="flex items-center gap-3">
            <HugeiconsIcon icon={Alert02Icon} className="size-5 shrink-0 text-destructive" />
            <p className="text-sm text-muted-foreground">
              Một số số liệu chưa tải được. Kiểm tra kết nối backend rồi thử lại.
            </p>
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              results.forEach((result) => void result.refetch())
            }}
          >
            Thử lại
          </Button>
        </div>
      ) : null}
    </AppShell>
  )
}
