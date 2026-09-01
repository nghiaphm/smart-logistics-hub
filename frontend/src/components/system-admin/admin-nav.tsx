"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  BoxesIcon,
  Car01Icon,
  DashboardSquare01Icon,
  Invoice01Icon,
  Location01Icon,
  PackageAddIcon,
  PackageIcon,
  TruckIcon,
  UserCircleIcon,
  UserMultipleIcon,
  WarehouseIcon,
} from "@hugeicons/core-free-icons"

import { cn } from "@/lib/utils"

type AdminNavItem = {
  label: string
  href: string
  icon: typeof DashboardSquare01Icon
}

const systemItems: AdminNavItem[] = [
  { label: "Tổng quan", href: "/admin", icon: DashboardSquare01Icon },
  { label: "Người dùng", href: "/admin/users", icon: UserMultipleIcon },
  { label: "Kho bãi", href: "/admin/warehouses", icon: WarehouseIcon },
]

const logisticsItems: AdminNavItem[] = [
  { label: "Tổng quan điều hành", href: "/logistics", icon: DashboardSquare01Icon },
  { label: "Toàn bộ đơn hàng", href: "/logistics/orders", icon: Invoice01Icon },
  { label: "Theo dõi Đơn hàng", href: "/logistics/tracking", icon: Location01Icon },
  { label: "Chuyến xe", href: "/logistics/trips", icon: TruckIcon },
  { label: "Tồn kho", href: "/logistics/inventory", icon: BoxesIcon },
  { label: "Sản phẩm", href: "/logistics/products", icon: PackageIcon },
  { label: "Nhập kho", href: "/logistics/inbounds", icon: PackageAddIcon },
]

const fleetItems: AdminNavItem[] = [
  { label: "Phương tiện", href: "/logistics/vehicles", icon: Car01Icon },
  { label: "Tài xế", href: "/logistics/drivers", icon: UserCircleIcon },
]

function renderItems(items: AdminNavItem[], pathname: string, base?: string) {
  return items.map((item) => {
    const href = base ? `${base}${item.href}` : item.href
    const isActive =
      pathname === href ||
      (item.href !== "/admin" && pathname.startsWith(`${href}/`))
    return (
      <Link
        key={href}
        href={href}
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
  })
}

export function AdminNav() {
  const pathname = usePathname()

  return (
    <nav className="flex flex-col gap-6 px-3">
      <div className="flex flex-col gap-1">
        <p className="px-3 pb-1 text-[11px] font-medium uppercase tracking-[0.16em] text-admin-sidebar-muted/80">
          Hệ thống
        </p>
        {renderItems(systemItems, pathname)}
      </div>

      <div className="flex flex-col gap-1">
        <p className="px-3 pb-1 text-[11px] font-medium uppercase tracking-[0.16em] text-admin-sidebar-muted/80">
          Điều hành vận tải
        </p>
        {renderItems(logisticsItems, pathname, "/admin/logistics")}
      </div>
      <div className="flex flex-col gap-1">
        <p className="px-3 pb-1 text-[11px] font-medium uppercase tracking-[0.16em] text-admin-sidebar-muted/80">
          Đội xe & Tài xế
        </p>
        {renderItems(fleetItems, pathname, "/admin/logistics")}
      </div>
    </nav>
  )
}
