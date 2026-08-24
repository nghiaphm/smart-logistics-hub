"use client"

import { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"

type Profile = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.ProfileResponse"]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export default function Page() {
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["profile"],
    queryFn: () => apiClient<Profile>("/profile"),
  })

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải thông tin hồ sơ",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

  if (isLoading) {
    return (
      <div className="flex flex-col gap-3">
        <Skeleton className="h-28 w-full rounded-2xl" />
        <Skeleton className="h-44 w-full rounded-2xl" />
      </div>
    )
  }

  if (isError) {
    return (
      <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-border bg-card px-6 py-14 text-center">
        <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
        <div>
          <p className="font-medium">Không thể tải thông tin hồ sơ</p>
          <p className="mt-1 text-sm text-muted-foreground">{errorMessage(error)}</p>
        </div>
        <Button variant="outline" size="sm" onClick={() => void refetch()}>
          Thử lại
        </Button>
      </div>
    )
  }

  if (!data) {
    return (
      <div className="flex flex-col items-center justify-center gap-2 rounded-2xl border border-dashed border-border bg-card px-6 py-14 text-center">
        <p className="font-medium">Chưa có thông tin hồ sơ</p>
        <p className="text-sm text-muted-foreground">Hồ sơ sẽ được tạo tự động khi bạn đăng nhập.</p>
      </div>
    )
  }

  const initials = (data.display_name || data.user_sub || "U").slice(0, 2).toUpperCase()

  return (
    <div className="flex max-w-xl flex-col gap-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Hồ sơ</h1>
        <p className="mt-1 text-sm text-muted-foreground">Thông tin tài khoản của bạn</p>
      </div>
      <div className="flex items-center gap-4 rounded-2xl border border-border bg-card p-5">
        <span className="flex size-14 shrink-0 items-center justify-center rounded-full bg-sidebar-primary text-lg font-semibold text-sidebar-primary-foreground">
          {initials}
        </span>
        <div className="min-w-0">
          <p className="truncate text-lg font-medium">{data.display_name || "Chưa đặt tên"}</p>
          <p className="truncate text-sm text-muted-foreground">{data.user_sub}</p>
        </div>
      </div>
      <div className="rounded-2xl border border-border bg-card p-5">
        <dl className="flex flex-col gap-4 text-sm">
          <div className="flex items-center justify-between gap-4">
            <dt className="text-muted-foreground">Tên hiển thị</dt>
            <dd className="font-medium">{data.display_name || "—"}</dd>
          </div>
          <div className="flex items-center justify-between gap-4">
            <dt className="text-muted-foreground">Điện thoại</dt>
            <dd className="font-medium">{data.phone || "—"}</dd>
          </div>
          <div className="flex items-center justify-between gap-4">
            <dt className="text-muted-foreground">Tham gia từ</dt>
            <dd className="font-medium">{data.created_at}</dd>
          </div>
        </dl>
      </div>
    </div>
  )
}
