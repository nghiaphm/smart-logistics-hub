"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  DashboardSquare01Icon,
  UserMultipleIcon,
  WarehouseIcon,
} from "@hugeicons/core-free-icons"

import { cn } from "@/lib/utils"

type AdminNavItem = {
  label: string
  href: string
  icon: typeof DashboardSquare01Icon
}

const items: AdminNavItem[] = [
  { label: "Tổng quan", href: "/admin", icon: DashboardSquare01Icon },
  { label: "Người dùng", href: "/admin/users", icon: UserMultipleIcon },
  { label: "Kho bãi", href: "/admin/warehouses", icon: WarehouseIcon },
]

export function AdminNav() {
  const pathname = usePathname()

  return (
    <nav className="flex flex-col gap-1 px-3">
      <p className="px-3 pb-1 text-[11px] font-medium uppercase tracking-[0.16em] text-admin-sidebar-muted/80">
        Hệ thống
      </p>
      {items.map((item) => {
        const isActive =
          pathname === item.href ||
          (item.href !== "/admin" && pathname.startsWith(`${item.href}/`))
        return (
          <Link
            key={item.href}
            href={item.href}
            aria-current={isActive ? "page" : undefined}
            className={cn(
              "group relative flex items-center gap-3 rounded-xl px-3 py-2 text-[13px] outline-none transition-colors",
              "focus-visible:ring-2 focus-visible:ring-admin-accent/60",
              isActive
                ? "bg-admin-accent/15 font-medium text-admin-sidebar-foreground"
                : "text-admin-sidebar-muted hover:bg-white/5 hover:text-admin-sidebar-foreground"
            )}
          >
            <span aria-hidden className="relative flex size-4 shrink-0 items-center justify-center">
              <span
                className={cn(
                  "size-2 rounded-[3px] transition-colors",
                  isActive ? "bg-admin-accent" : "bg-white/25 group-hover:bg-white/50"
                )}
              />
            </span>
            <HugeiconsIcon
              icon={item.icon}
              className={cn(
                "size-4 shrink-0 transition-colors",
                isActive ? "text-admin-accent" : "text-admin-sidebar-muted"
              )}
            />
            <span className="truncate">{item.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}
