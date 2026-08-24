"use client"

import type { ComponentProps } from "react"

import { DialogFooter } from "@/components/ui/dialog"
import { cn } from "@/lib/utils"

export function AppModalActions({ className, ...props }: ComponentProps<typeof DialogFooter>) {
  return (
    <DialogFooter
      data-slot="app-modal-actions"
      className={cn("mt-1", className)}
      {...props}
    />
  )
}
