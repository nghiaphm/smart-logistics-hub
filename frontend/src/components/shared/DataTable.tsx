import type { ReactNode } from "react"

import { HugeiconsIcon } from "@hugeicons/react"
import { InboxIcon } from "@hugeicons/core-free-icons"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Skeleton } from "@/components/ui/skeleton"
import { cn } from "@/lib/utils"

export type Column<T> = {
  key: string
  header: ReactNode
  cell: (row: T) => ReactNode
  className?: string
  headerClassName?: string
}

type DataTableProps<T> = {
  columns: Column<T>[]
  rows: T[]
  rowKey: (row: T) => string | number
  loading?: boolean
  emptyText?: string
  emptyDescription?: string
  className?: string
}

export function DataTable<T>({
  columns,
  rows,
  rowKey,
  loading = false,
  emptyText = "Không có dữ liệu",
  emptyDescription,
  className,
}: DataTableProps<T>) {
  return (
    <div data-slot="data-table" className={cn("rounded-2xl border border-border bg-card", className)}>
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            {columns.map((column) => (
              <TableHead key={column.key} className={column.headerClassName}>
                {column.header}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {loading
            ? Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={index} className="hover:bg-transparent">
                  {columns.map((column) => (
                    <TableCell key={column.key}>
                      <Skeleton className="h-4 w-full max-w-40" />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            : rows.length === 0
              ? (
                  <TableRow className="hover:bg-transparent">
                    <TableCell colSpan={columns.length} className="h-24">
                      <div className="flex flex-col items-center justify-center gap-2 text-center">
                        <HugeiconsIcon icon={InboxIcon} className="size-6 text-muted-foreground/50" />
                        <p className="text-sm text-muted-foreground">{emptyText}</p>
                        {emptyDescription ? (
                          <p className="text-xs text-muted-foreground/70">{emptyDescription}</p>
                        ) : null}
                      </div>
                    </TableCell>
                  </TableRow>
                )
              : rows.map((row) => (
                  <TableRow key={rowKey(row)}>
                    {columns.map((column) => (
                      <TableCell key={column.key} className={column.className}>
                        {column.cell(row)}
                      </TableCell>
                    ))}
                  </TableRow>
                ))}
        </TableBody>
      </Table>
    </div>
  )
}
