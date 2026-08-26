"use client"

import Link from "next/link"
import { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, ArrowRight01Icon } from "@hugeicons/core-free-icons"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Skeleton } from "@/components/ui/skeleton"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

type PaginatedWorkspaces = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_workspace_dto.PaginatedResponse"]

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

interface WorkspaceModalProps {
  onClose: () => void
}

export function WorkspaceModal({ onClose }: WorkspaceModalProps) {
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

  const workspaces = data?.items ?? []

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <Card className="w-full max-w-lg relative animate-in fade-in zoom-in-95 duration-150 border-border shadow-2xl bg-card">
        <CardHeader className="pb-4">
          <div className="flex items-center justify-between">
            <CardTitle className="text-xl font-bold tracking-tight">Chọn Không Gian Làm Việc</CardTitle>
            <button 
              onClick={onClose}
              className="rounded-lg p-1.5 text-muted-foreground hover:bg-neutral-100 dark:hover:bg-neutral-800 transition-colors"
            >
              ✕
            </button>
          </div>
          <CardDescription>
            Vui lòng chọn một workspace phù hợp để bắt đầu quy trình vận hành Logistics
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 max-h-[350px] overflow-y-auto pr-1">
          {isLoading && (
            <div className="space-y-2 py-4">
              <Skeleton className="h-14 w-full rounded-xl" />
              <Skeleton className="h-14 w-full rounded-xl" />
              <Skeleton className="h-14 w-full rounded-xl" />
            </div>
          )}

          {isError && (
            <div className="flex flex-col items-center justify-center gap-3 py-6 text-center">
              <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
              <div>
                <p className="text-sm font-medium">Không thể tải danh sách workspace</p>
                <p className="mt-1 text-xs text-muted-foreground">{errorMessage(error)}</p>
              </div>
              <Button variant="outline" size="sm" onClick={() => void refetch()}>
                Thử lại
              </Button>
            </div>
          )}

          {!isLoading && !isError && workspaces.length === 0 && (
            <div className="text-center py-8 text-sm text-muted-foreground border border-dashed border-border rounded-xl">
              Chưa có workspace nào được thiết lập trên hệ thống.
            </div>
          )}

          {!isLoading && !isError && workspaces.length > 0 && (
            <div className="space-y-3">
              {workspaces.map((workspace) => (
                <Link 
                  key={workspace.id} 
                  href={`/${workspace.id}`}
                  className="flex items-center justify-between p-4 rounded-xl border border-border/60 bg-neutral-100/40 dark:bg-neutral-900/10 hover:border-indigo-500/50 hover:bg-indigo-500/[0.02] transition-all group"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-mono font-semibold px-2 py-0.5 rounded bg-zinc-500/10 text-neutral-800 dark:text-neutral-200">
                        {workspace.workspace_code}
                      </span>
                      <span className="font-semibold text-neutral-900 dark:text-neutral-100">
                        {workspace.name}
                      </span>
                    </div>
                    {workspace.description && (
                      <p className="text-xs text-muted-foreground line-clamp-1">
                        {workspace.description}
                      </p>
                    )}
                  </div>
                  <HugeiconsIcon 
                    icon={ArrowRight01Icon} 
                    className="h-5 w-5 text-muted-foreground group-hover:text-indigo-500 transition-colors" 
                  />
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
