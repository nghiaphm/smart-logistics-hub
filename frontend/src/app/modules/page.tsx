"use client"

import { useState } from "react"
import { HugeiconsIcon } from "@hugeicons/react"
import { DeliveryTruck01Icon, UserGroupIcon, SecurityValidationIcon } from "@hugeicons/core-free-icons"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { WorkspaceModal } from "@/components/auth/WorkspaceModal"

export default function ModulesPage() {
  const [showWorkspaceModal, setShowWorkspaceModal] = useState(false)

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-zinc-50 p-6 dark:bg-zinc-950">
      <div className="mx-auto w-full max-w-4xl space-y-8">
        {/* Header */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Chọn Phân Hệ Quản Trị</h1>
          <p className="text-muted-foreground text-sm max-w-lg mx-auto">
            Hệ sinh thái điều phối toàn diện. Vui lòng chọn module nghiệp vụ phù hợp để bắt đầu phiên làm việc.
          </p>
        </div>

        {/* Modules Grid */}
        <div className="grid gap-6 md:grid-cols-2">
          {/* Module 1: Logistics - ACTIVE */}
          <Card 
            className="group relative cursor-pointer border-border transition-all duration-200 hover:-translate-y-1 hover:border-indigo-500/50 hover:shadow-lg hover:shadow-indigo-500/5"
            onClick={() => setShowWorkspaceModal(true)}
          >
            <CardHeader className="space-y-1 pb-4">
              <div className="flex items-center justify-between">
                <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-indigo-50 text-indigo-600 transition-colors group-hover:bg-indigo-100 dark:bg-indigo-950/40 dark:text-indigo-400">
                  <HugeiconsIcon icon={DeliveryTruck01Icon} className="h-6 w-6" />
                </span>
                <Badge variant="secondary" className="bg-emerald-50 text-emerald-700 hover:bg-emerald-50 dark:bg-emerald-950/30 dark:text-emerald-400">
                  Đang hoạt động
                </Badge>
              </div>
              <CardTitle className="text-xl font-bold mt-4 transition-colors group-hover:text-indigo-600 dark:group-hover:text-indigo-400">
                Phân hệ Logistics
              </CardTitle>
              <CardDescription className="text-sm text-muted-foreground">
                Quản lý kho bãi, theo dõi hàng hoá tồn kho, danh mục sản phẩm và điều phối đội xe vận tải thông minh.
              </CardDescription>
            </CardHeader>
            <CardContent className="text-xs text-muted-foreground border-t border-border/50 pt-4 flex gap-4">
              <div className="flex items-center gap-1">
                <span className="h-1.5 w-1.5 rounded-full bg-indigo-500" />
                Quản lý kho bãi
              </div>
              <div className="flex items-center gap-1">
                <span className="h-1.5 w-1.5 rounded-full bg-indigo-500" />
                Điều phối đội xe
              </div>
            </CardContent>
          </Card>

          {/* Module 2: Nhân sự - INACTIVE (Chừa layout, không code rỗng/trực tiếp route) */}
          <Card className="relative opacity-60 border-dashed border-border bg-neutral-100/40 dark:bg-neutral-900/10">
            <CardHeader className="space-y-1 pb-4">
              <div className="flex items-center justify-between">
                <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-zinc-100 text-zinc-500 dark:bg-zinc-800 dark:text-zinc-400">
                  <HugeiconsIcon icon={UserGroupIcon} className="h-6 w-6" />
                </span>
                <Badge variant="outline" className="gap-1 bg-zinc-50 text-zinc-600 dark:bg-zinc-950/20 dark:text-zinc-400">
                  <HugeiconsIcon icon={SecurityValidationIcon} className="h-3 w-3" />
                  Sắp ra mắt
                </Badge>
              </div>
              <CardTitle className="text-xl font-bold mt-4 text-zinc-700 dark:text-zinc-300">
                Phân hệ Nhân sự
              </CardTitle>
              <CardDescription className="text-sm text-zinc-500 dark:text-zinc-400">
                Quản lý chấm công, bảng lương, đánh giá KPI và hồ sơ nhân sự dành cho bộ phận hành chính.
              </CardDescription>
            </CardHeader>
            <CardContent className="text-xs text-zinc-400 dark:text-zinc-500 border-t border-dashed border-border/50 pt-4 flex gap-4">
              <div className="flex items-center gap-1">
                <span className="h-1.5 w-1.5 rounded-full bg-zinc-400" />
                Quản lý nhân sự
              </div>
              <div className="flex items-center gap-1">
                <span className="h-1.5 w-1.5 rounded-full bg-zinc-400" />
                Chấm công & Lương
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Modal workspace selection */}
        {showWorkspaceModal && (
          <WorkspaceModal onClose={() => setShowWorkspaceModal(false)} />
        )}
      </div>
    </div>
  )
}
