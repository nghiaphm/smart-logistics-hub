import type { ComponentProps, ReactNode } from "react"

import { cn } from "@/lib/utils"

export function Form({ className, ...props }: ComponentProps<"form">) {
  return (
    <form data-slot="form" className={cn("flex flex-col gap-5", className)} {...props} />
  )
}

type FormFieldProps = {
  label: string
  htmlFor?: string
  hint?: string
  error?: string
  required?: boolean
  children: ReactNode
  className?: string
}

export function FormField({
  label,
  htmlFor,
  hint,
  error,
  required,
  children,
  className,
}: FormFieldProps) {
  return (
    <div data-slot="form-field" className={cn("flex flex-col gap-1.5", className)}>
      <label htmlFor={htmlFor} className="text-sm font-medium text-foreground">
        {label}
        {required ? <span className="ml-0.5 text-destructive">*</span> : null}
      </label>
      {children}
      {hint && !error ? <p className="text-xs text-muted-foreground">{hint}</p> : null}
      {error ? <p className="text-xs text-destructive">{error}</p> : null}
    </div>
  )
}

export function FormActions({ className, ...props }: ComponentProps<"div">) {
  return (
    <div
      data-slot="form-actions"
      className={cn(
        "flex flex-col-reverse gap-2 pt-1 sm:flex-row sm:justify-end",
        className
      )}
      {...props}
    />
  )
}
