"use client"

import { createContext, useContext, useState } from "react"
import type { ReactNode } from "react"
import { Drawer } from "@base-ui/react/drawer"
import { HugeiconsIcon } from "@hugeicons/react"
import { DeliveryTruck01Icon, Menu01Icon } from "@hugeicons/core-free-icons"

import { Button } from "@/components/ui/button"
import { NavMain } from "@/components/nav-main"

type SidebarContextValue = {
  openMobile: boolean
  setOpenMobile: (open: boolean) => void
}

const SidebarContext = createContext<SidebarContextValue | null>(null)

export function SidebarProvider({ children }: { children: ReactNode }) {
  const [openMobile, setOpenMobile] = useState(false)

  return (
    <SidebarContext.Provider value={{ openMobile, setOpenMobile }}>
      {children}
    </SidebarContext.Provider>
  )
}

function useSidebar(): SidebarContextValue {
  const context = useContext(SidebarContext)
  if (!context) {
    throw new Error("useSidebar must be used within SidebarProvider")
  }
  return context
}

function SidebarContent() {
  return (
    <div className="flex h-full flex-col">
      <div className="flex h-16 shrink-0 items-center gap-2.5 border-b border-sidebar-border px-5">
        <HugeiconsIcon icon={DeliveryTruck01Icon} className="size-5 text-sidebar-primary" />
        <div className="leading-tight">
          <p className="text-sm font-semibold text-sidebar-foreground">Smart Logistics</p>
          <p className="text-[11px] text-muted-foreground">Trung tâm điều hành</p>
        </div>
      </div>
      <div className="flex-1 overflow-y-auto py-4">
        <NavMain />
      </div>
    </div>
  )
}

export function SidebarTrigger() {
  const { setOpenMobile } = useSidebar()

  return (
    <Button
      variant="ghost"
      size="icon"
      aria-label="Mở menu điều hướng"
      className="lg:hidden"
      onClick={() => setOpenMobile(true)}
    >
      <HugeiconsIcon icon={Menu01Icon} />
    </Button>
  )
}

export function AppSidebar() {
  const { openMobile, setOpenMobile } = useSidebar()

  return (
    <>
      <aside className="hidden w-64 shrink-0 border-r border-sidebar-border bg-sidebar lg:block">
        <SidebarContent />
      </aside>

      <Drawer.Root open={openMobile} onOpenChange={setOpenMobile}>
        <Drawer.Portal>
          <Drawer.Backdrop className="fixed inset-0 z-40 bg-black/40 lg:hidden" />
          <Drawer.Popup className="fixed inset-y-0 left-0 z-50 w-72 max-w-[85vw] bg-sidebar shadow-2xl outline-none transition-transform duration-200 motion-reduce:transition-none data-[open]:translate-x-0 data-[closed]:-translate-x-full lg:hidden">
            <Drawer.Title className="sr-only">Menu điều hướng</Drawer.Title>
            <SidebarContent />
          </Drawer.Popup>
        </Drawer.Portal>
      </Drawer.Root>
    </>
  )
}
