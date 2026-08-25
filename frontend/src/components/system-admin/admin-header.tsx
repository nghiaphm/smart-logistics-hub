"use client"

import { useSyncExternalStore } from "react"
import { HugeiconsIcon } from "@hugeicons/react"
import { Clock01Icon, Logout01Icon } from "@hugeicons/core-free-icons"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { AdminSidebarTrigger } from "@/components/system-admin/admin-sidebar"
import { useAuth } from "@/contexts/auth.context"

let clockCache: { at: number; value: Date } | null = null

function getClockSnapshot(): Date {
  const at = Date.now()
  if (clockCache && at - clockCache.at < 1000) {
    return clockCache.value
  }
  clockCache = { at, value: new Date(at) }
  return clockCache.value
}

function subscribeClock(onStoreChange: () => void): () => void {
  const id = setInterval(onStoreChange, 1000)
  return () => clearInterval(id)
}

const clockFormatter = new Intl.DateTimeFormat("vi-VN", {
  timeZone: "Asia/Ho_Chi_Minh",
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
  hour12: false,
})

export function AdminHeader() {
  const { user, logout } = useAuth()
  const now = useSyncExternalStore(subscribeClock, getClockSnapshot, () => null)

  const username = user?.preferred_username ?? "user"
  const initials = username.slice(0, 2).toUpperCase()

  return (
    <header className="flex h-16 shrink-0 items-center gap-3 border-b border-border bg-background/80 px-4 backdrop-blur lg:px-6">
      <AdminSidebarTrigger />
      <span className="inline-flex shrink-0 items-center rounded-md bg-admin-accent px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.14em] text-admin-accent-foreground">
        Admin
      </span>
      <div className="min-w-0">
        <p className="truncate text-sm font-medium text-foreground">Quản trị hệ thống</p>
      </div>
      <div className="ml-auto flex items-center gap-2">
        <div className="hidden items-center gap-1.5 rounded-full border border-border px-3 py-1 font-mono text-xs tabular-nums text-muted-foreground sm:flex">
          <HugeiconsIcon icon={Clock01Icon} className="size-3.5" />
          <span>{now ? clockFormatter.format(now) : "—"}</span>
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger render={<Button variant="ghost" className="gap-2 rounded-full px-2" />}>
            <span className="flex size-7 items-center justify-center rounded-full bg-admin-accent text-[11px] font-semibold text-admin-accent-foreground">
              {initials}
            </span>
            <span className="hidden max-w-28 truncate text-sm sm:block">{username}</span>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-52">
            <DropdownMenuGroup>
              <DropdownMenuLabel>{username}</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={logout}>
                <HugeiconsIcon icon={Logout01Icon} />
                Đăng xuất
              </DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  )
}
