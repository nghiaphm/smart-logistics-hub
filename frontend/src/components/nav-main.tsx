"use client"

import Link from "next/link"
import { useParams, usePathname } from "next/navigation"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  BoxesIcon,
  DashboardSquare01Icon,
  Location01Icon,
  PackageAddIcon,
  PackageIcon,
  TruckIcon,
} from "@hugeicons/core-free-icons"

import { cn } from "@/lib/utils"

type NavItem = {
  label: string
  href: string
  icon: typeof DashboardSquare01Icon
}

type NavSection = {
  label: string
  items: NavItem[]
}

const sections: NavSection[] = [
  {
    label: "Điều hành",
    items: [
      { label: "Tổng quan", href: "/logistics", icon: DashboardSquare01Icon },
      { label: "Hàng nhập", href: "/logistics/inbounds", icon: PackageAddIcon },
      { label: "Tồn kho", href: "/logistics/inventory", icon: BoxesIcon },
      { label: "Sản phẩm", href: "/logistics/products", icon: PackageIcon },
    ],
  },
  {
    label: "Vận tải",
    items: [
      { label: "Theo dõi đơn", href: "/logistics/tracking", icon: Location01Icon },
      { label: "Chuyến xe", href: "/logistics/trips", icon: TruckIcon },
    ],
  },
]

export function NavMain() {
  const params = useParams<{ workspace_id: string }>()
  const pathname = usePathname()
  const base = `/${params.workspace_id}`

  return (
    <nav className="flex flex-col gap-6 px-3">
      {sections.map((section) => (
        <div key={section.label} className="flex flex-col gap-1">
          <p className="px-3 pb-1 text-[11px] font-medium uppercase tracking-[0.16em] text-muted-foreground">
            {section.label}
          </p>
          <div className="relative flex flex-col">
            <span
              aria-hidden
              className="absolute top-3 bottom-3 left-[19px] w-px bg-border/70"
            />
            {section.items.map((item) => {
              const href = `${base}${item.href}`
              const isActive =
                pathname === href ||
                (item.href !== "" && pathname.startsWith(`${href}/`))
              return (
                <Link
                  key={item.href}
                  href={href}
                  aria-current={isActive ? "page" : undefined}
                  className={cn(
                    "group relative flex items-center gap-3 rounded-xl px-3 py-2 text-[13px] outline-none transition-colors",
                    "focus-visible:ring-2 focus-visible:ring-ring/60",
                    isActive
                      ? "bg-sidebar-accent font-medium text-sidebar-accent-foreground"
                      : "text-sidebar-foreground/70 hover:bg-sidebar-accent/70 hover:text-sidebar-foreground"
                  )}
                >
                  <span aria-hidden className="relative flex size-4 shrink-0 items-center justify-center">
                    <span
                      className={cn(
                        "size-2 rounded-full transition-colors",
                        isActive
                          ? "bg-sidebar-primary"
                          : "bg-muted-foreground/45 group-hover:bg-muted-foreground/75"
                      )}
                    />
                  </span>
                  <HugeiconsIcon icon={item.icon} className="size-4 shrink-0 text-muted-foreground" />
                  <span className="truncate">{item.label}</span>
                </Link>
              )
            })}
          </div>
        </div>
      ))}
    </nav>
  )
}
