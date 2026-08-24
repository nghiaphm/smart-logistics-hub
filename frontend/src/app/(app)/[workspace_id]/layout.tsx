import type { ReactNode } from "react"

import { AppHeader } from "@/components/app-header"
import { AppSidebar, SidebarProvider } from "@/components/app-sidebar"

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <SidebarProvider>
      <div className="flex h-svh bg-background">
        <AppSidebar />
        <div className="flex min-w-0 flex-1 flex-col">
          <AppHeader />
          <main className="flex-1 overflow-y-auto p-4 lg:p-6">{children}</main>
        </div>
      </div>
    </SidebarProvider>
  )
}
