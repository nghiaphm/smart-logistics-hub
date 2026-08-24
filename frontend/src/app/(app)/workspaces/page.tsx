"use client"

import Link from "next/link"
import { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon } from "@hugeicons/core-free-icons"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

type PaginatedWorkspaces = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.PaginatedResponse"]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

export default function Page() {
  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["workspaces"],
    queryFn: () => apiClient<PaginatedWorkspaces>("/workspaces"),
  })

  useEffect(() => {
    if (isError && error) {
      toast.add({
        title: "Không thể tải danh sách workspace",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error])

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
          <p className="font-medium">Không thể tải danh sách workspace</p>
          <p className="mt-1 text-sm text-muted-foreground">{errorMessage(error)}</p>
        </div>
        <Button variant="outline" size="sm" onClick={() => void refetch()}>
          Thử lại
        </Button>
      </div>
    )
  }

  const workspaces = data?.items ?? []

  return (
    <div className="flex flex-col gap-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Chọn workspace</h1>
        <p className="mt-1 text-sm text-muted-foreground">Chọn một workspace để vào khu vực điều hành</p>
      </div>
      {workspaces.length === 0 ? (
        <div className="flex flex-col items-center justify-center gap-2 rounded-2xl border border-dashed border-border bg-card px-6 py-14 text-center">
          <p className="font-medium">Chưa có workspace nào</p>
          <p className="text-sm text-muted-foreground">Workspace sẽ xuất hiện ở đây khi được tạo.</p>
        </div>
      ) : (
        <div className="rounded-2xl border border-border bg-card">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent">
                <TableHead>Mã workspace</TableHead>
                <TableHead>Tên</TableHead>
                <TableHead>Mô tả</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {workspaces.map((workspace, index) => (
                <TableRow key={workspace.id ?? workspace.workspace_code ?? index}>
                  <TableCell className="font-medium">{workspace.workspace_code}</TableCell>
                  <TableCell>
                    <Link href={`/${workspace.id}`} className="hover:underline">
                      {workspace.name}
                    </Link>
                  </TableCell>
                  <TableCell className="text-muted-foreground">{workspace.description}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  )
}
