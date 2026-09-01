"use client"

import Link from "next/link"
import { useParams, usePathname } from "next/navigation"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  AlertCircleIcon,
  BankIcon,
  Invoice01Icon,
  KeyIcon,
  Location01Icon,
  Store01Icon,
  WalletIcon,
} from "@hugeicons/core-free-icons"

import { cn } from "@/lib/utils"

type NavItem = {
  label: string
  href: string
  icon: typeof Invoice01Icon
}

type NavSection = {
  label: string
  items: NavItem[]
}

const sections: NavSection[] = [
  {
    label: "Quản lý Vận đơn",
    items: [
      { label: "Danh sách đơn hàng", href: "/logistic/orders", icon: Invoice01Icon },
      { label: "Tra cứu hành trình", href: "/tracking", icon: Location01Icon },
      { label: "Xử lý sự cố", href: "/incidents", icon: AlertCircleIcon },
    ],
  },
  {
    label: "Tài chính & COD",
    items: [
      { label: "Ví tiền thu hộ COD", href: "/cod-wallet", icon: WalletIcon },
      { label: "Lịch sử cước phí & Báo cáo đối soát", href: "/billing", icon: BankIcon },
    ],
  },
  {
    label: "Cài đặt Cửa hàng",
    items: [
      { label: "Địa chỉ lấy hàng", href: "/settings/pickup", icon: Store01Icon },
      { label: "Thông tin tài khoản & Kết nối API", href: "/settings/api", icon: KeyIcon },
    ],
  },
]

export function NavMain() {
  const params = useParams<{ workspace_id: string }>()
  const pathname = usePathname()
  const base = `/${params.workspace_id}`
  const visibleSections = sections.filter((section) => section.items.length > 0)

  return (
    <nav className="flex flex-col gap-6 px-3">
      {visibleSections.map((section) => (
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
