"use client"

import { createContext, useContext, useState } from "react"
import type { ReactNode } from "react"
import { Drawer } from "@base-ui/react/drawer"
import { HugeiconsIcon } from "@hugeicons/react"
import { CommandIcon, Menu01Icon } from "@hugeicons/core-free-icons"

import { Button } from "@/components/ui/button"
import { AdminNav } from "@/components/system-admin/admin-nav"

type AdminSidebarContextValue = {
  openMobile: boolean
  setOpenMobile: (open: boolean) => void
}

const AdminSidebarContext = createContext<AdminSidebarContextValue | null>(null)

export function AdminSidebarProvider({ children }: { children: ReactNode }) {
  const [openMobile, setOpenMobile] = useState(false)

  return (
    <AdminSidebarContext.Provider value={{ openMobile, setOpenMobile }}>
      {children}
    </AdminSidebarContext.Provider>
  )
}

function useAdminSidebar(): AdminSidebarContextValue {
  const context = useContext(AdminSidebarContext)
  if (!context) {
    throw new Error("useAdminSidebar must be used within AdminSidebarProvider")
  }
  return context
}

function AdminSidebarContent() {
  return (
    <div className="flex h-full flex-col">
      <div className="flex h-16 shrink-0 items-center gap-2.5 border-b border-admin-sidebar-border px-5">
        <span className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-admin-accent text-admin-accent-foreground">
          <HugeiconsIcon icon={CommandIcon} className="size-4.5" />
        </span>
        <div className="leading-tight">
          <p className="text-sm font-semibold text-admin-sidebar-foreground">Smart Logistics</p>
          <p className="text-[11px] text-admin-sidebar-muted">Bảng điều khiển hệ thống</p>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto py-4">
        <AdminNav />
      </div>
    </div>
  )
}

export function AdminSidebarTrigger() {
  const { setOpenMobile } = useAdminSidebar()

  return (
    <Button
      variant="ghost"
      size="icon"
      aria-label="Mở menu quản trị"
      className="lg:hidden"
      onClick={() => setOpenMobile(true)}
    >
      <HugeiconsIcon icon={Menu01Icon} />
    </Button>
  )
}

export function AdminSidebar() {
  const { openMobile, setOpenMobile } = useAdminSidebar()

  return (
    <>
      <aside className="hidden w-64 shrink-0 border-r border-admin-sidebar-border bg-admin-sidebar lg:block">
        <AdminSidebarContent />
      </aside>

      <Drawer.Root open={openMobile} onOpenChange={setOpenMobile}>
        <Drawer.Portal>
          <Drawer.Backdrop className="fixed inset-0 z-40 bg-black/40 lg:hidden" />
          <Drawer.Popup className="fixed inset-y-0 left-0 z-50 w-72 max-w-[85vw] bg-admin-sidebar shadow-2xl outline-none transition-transform duration-200 motion-reduce:transition-none data-[open]:translate-x-0 data-[closed]:-translate-x-full lg:hidden">
            <Drawer.Title className="sr-only">Menu quản trị</Drawer.Title>
            <AdminSidebarContent />
          </Drawer.Popup>
        </Drawer.Portal>
      </Drawer.Root>
    </>
  )
}
