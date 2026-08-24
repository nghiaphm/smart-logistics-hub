import type { ReactNode } from "react"

import { cn } from "@/lib/utils"

type AppShellProps = {
  title: string
  description?: string
  actions?: ReactNode
  children: ReactNode
  className?: string
}

export function AppShell({ title, description, actions, children, className }: AppShellProps) {
  return (
    <div data-slot="app-shell" className={cn("flex flex-col gap-6", className)}>
      <header className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight text-foreground">{title}</h1>
          {description ? (
            <p className="mt-1 text-sm text-muted-foreground">{description}</p>
          ) : null}
        </div>
        {actions ? <div className="flex items-center gap-2">{actions}</div> : null}
      </header>
      <div className="flex flex-col gap-4">{children}</div>
    </div>
  )
}
