import type { ReactNode } from "react"

import { AdminHeader } from "@/components/system-admin/admin-header"
import {
  AdminSidebar,
  AdminSidebarProvider,
} from "@/components/system-admin/admin-sidebar"

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <AdminSidebarProvider>
      <div className="flex h-svh bg-background">
        <AdminSidebar />
        <div className="flex min-w-0 flex-1 flex-col">
          <AdminHeader />
          <main className="flex-1 overflow-y-auto p-4 lg:p-6">{children}</main>
        </div>
      </div>
    </AdminSidebarProvider>
  )
}
