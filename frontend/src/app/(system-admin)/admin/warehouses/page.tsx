"use client"

import { useEffect, useState } from "react"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { HugeiconsIcon } from "@hugeicons/react"
import {
  Add01Icon,
  Alert02Icon,
  Delete01Icon,
  Edit02Icon,
} from "@hugeicons/core-free-icons"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { AppShell } from "@/components/shared/AppShell"
import { DataTable } from "@/components/shared/DataTable"
import type { Column } from "@/components/shared/DataTable"
import { AppModalActions, AppModalShell } from "@/components/shared/modal"
import { WarehouseFormModal } from "@/components/system-admin/warehouse-form-modal"
import type { WarehouseFormValues } from "@/components/system-admin/warehouse-form-modal"

type Warehouse = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.WarehouseResponse"]
type WarehouseList = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.PaginatedResponse"]
type CreateWarehouseRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.CreateWarehouseRequest"]
type UpdateWarehouseRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_warehouse_dto.UpdateWarehouseRequest"]

type ModalState =
  | { type: "form"; mode: "create" | "edit"; warehouse?: Warehouse }
  | { type: "delete"; warehouse: Warehouse }
  | null

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function toLocation(values: WarehouseFormValues): { lat: number; lng: number } | undefined {
  const lat = Number.parseFloat(values.lat)
  const lng = Number.parseFloat(values.lng)
  if (Number.isNaN(lat) || Number.isNaN(lng)) {
    return undefined
  }
  return { lat, lng }
}

const columns: Column<Warehouse>[] = [
  {
    key: "code",
    header: "Mã kho",
    cell: (warehouse) => <span className="font-medium">{warehouse.warehouse_code}</span>,
  },
  { key: "name", header: "Tên kho", cell: (warehouse) => warehouse.name },
  {
    key: "address",
    header: "Địa chỉ",
    cell: (warehouse) => (
      <span className="line-clamp-1 max-w-72 text-muted-foreground">{warehouse.address}</span>
    ),
  },
  { key: "manager", header: "Quản lý", cell: (warehouse) => warehouse.manager_name ?? "—" },
  {
    key: "phone",
    header: "Liên hệ",
    cell: (warehouse) => warehouse.contact_phone ?? "—",
    className: "text-muted-foreground",
  },
  {
    key: "status",
    header: "Trạng thái",
    cell: (warehouse) => (
      <Badge variant={warehouse.is_active ? "default" : "outline"}>
        {warehouse.is_active ? "Hoạt động" : "Ngừng"}
      </Badge>
    ),
  },
]

export default function AdminWarehousesPage() {
  const queryClient = useQueryClient()
  const [modal, setModal] = useState<ModalState>(null)

  const listQuery = useQuery({
    queryKey: ["admin", "warehouses"],
    queryFn: () => apiClient<WarehouseList>("/warehouses?limit=100"),
  })

  useEffect(() => {
    if (listQuery.isError && listQuery.error) {
      toast.add({
        title: "Không thể tải danh sách kho",
        description: errorMessage(listQuery.error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [listQuery.isError, listQuery.error])

  const invalidateList = () => {
    queryClient.invalidateQueries({ queryKey: ["admin", "warehouses"] })
  }

  const createMutation = useMutation({
    mutationFn: (values: WarehouseFormValues) => {
      const payload: CreateWarehouseRequest = {
        warehouse_code: values.warehouse_code.trim(),
        name: values.name.trim(),
        address: values.address.trim(),
        manager_name: values.manager_name.trim() || undefined,
        contact_phone: values.contact_phone.trim() || undefined,
        location: toLocation(values),
      }
      return apiClient<Warehouse>("/warehouses", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
    },
    onSuccess: () => {
      toast.add({ title: "Đã thêm kho mới", type: "success", timeout: 4000 })
      invalidateList()
      setModal(null)
    },
    onError: (error) => {
      toast.add({
        title: "Không thể thêm kho",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: number; values: WarehouseFormValues }) => {
      const payload: UpdateWarehouseRequest = {
        name: values.name.trim(),
        address: values.address.trim(),
        manager_name: values.manager_name.trim() || undefined,
        contact_phone: values.contact_phone.trim() || undefined,
        location: toLocation(values),
        is_active: values.is_active,
      }
      return apiClient<Warehouse>(`/warehouses/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
    },
    onSuccess: () => {
      toast.add({ title: "Đã cập nhật kho", type: "success", timeout: 4000 })
      invalidateList()
      setModal(null)
    },
    onError: (error) => {
      toast.add({
        title: "Không thể cập nhật kho",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiClient<void>(`/warehouses/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      toast.add({ title: "Đã xoá kho", type: "success", timeout: 4000 })
      invalidateList()
      setModal(null)
    },
    onError: (error) => {
      toast.add({
        title: "Không thể xoá kho",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    },
  })

  const deleteTarget = modal?.type === "delete" ? modal.warehouse : null
  const warehouses = listQuery.data?.items ?? []

  return (
    <AppShell
      title="Kho bãi"
      description="Danh sách kho bãi trên toàn hệ thống"
      actions={
        <Button onClick={() => setModal({ type: "form", mode: "create" })}>
          <HugeiconsIcon icon={Add01Icon} />
          Thêm kho
        </Button>
      }
    >
      {listQuery.isError ? (
        <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-border bg-card px-6 py-14 text-center">
          <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
          <div>
            <p className="font-medium">Không thể tải danh sách kho</p>
            <p className="mt-1 text-sm text-muted-foreground">{errorMessage(listQuery.error)}</p>
          </div>
          <Button variant="outline" size="sm" onClick={() => void listQuery.refetch()}>
            Thử lại
          </Button>
        </div>
      ) : (
        <DataTable
          columns={[
            ...columns,
            {
              key: "actions",
              header: "",
              className: "text-right",
              cell: (warehouse) => (
                <div className="flex items-center justify-end gap-1">
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label="Sửa kho"
                    onClick={() =>
                      setModal({ type: "form", mode: "edit", warehouse })
                    }
                  >
                    <HugeiconsIcon icon={Edit02Icon} className="size-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label="Xoá kho"
                    className="text-destructive hover:text-destructive"
                    onClick={() => setModal({ type: "delete", warehouse })}
                  >
                    <HugeiconsIcon icon={Delete01Icon} className="size-4" />
                  </Button>
                </div>
              ),
            },
          ]}
          rows={warehouses}
          rowKey={(warehouse) => warehouse.id ?? warehouse.warehouse_code ?? ""}
          loading={listQuery.isLoading}
          emptyText="Chưa có kho nào. Bấm “Thêm kho” để tạo mới."
        />
      )}

      {modal?.type === "form" ? (
        <WarehouseFormModal
          key={modal.mode === "create" ? "create" : `edit-${modal.warehouse?.id ?? "new"}`}
          open
          mode={modal.mode}
          warehouse={modal.warehouse}
          isSubmitting={createMutation.isPending || updateMutation.isPending}
          onOpenChange={(open) => {
            if (!open) {
              setModal(null)
            }
          }}
          onSubmit={(values) => {
            if (modal.mode === "create") {
              createMutation.mutate(values)
            } else if (modal.warehouse?.id != null) {
              updateMutation.mutate({ id: modal.warehouse.id, values })
            }
          }}
        />
      ) : null}

      <AppModalShell
        open={deleteTarget !== null}
        onOpenChange={(open) => {
          if (!open) {
            setModal(null)
          }
        }}
        title="Xoá kho"
        description={
          deleteTarget
            ? `Xoá kho ${deleteTarget.warehouse_code} — ${deleteTarget.name}?`
            : undefined
        }
        actions={
          <AppModalActions>
            <Button
              variant="outline"
              onClick={() => setModal(null)}
              disabled={deleteMutation.isPending}
            >
              Huỷ
            </Button>
            <Button
              variant="destructive"
              disabled={deleteMutation.isPending || deleteTarget == null}
              onClick={() => {
                if (deleteTarget?.id != null) {
                  deleteMutation.mutate(deleteTarget.id)
                }
              }}
            >
              {deleteMutation.isPending ? "Đang xoá…" : "Xoá"}
            </Button>
          </AppModalActions>
        }
      >
        <p className="text-sm text-muted-foreground">Hành động này không thể hoàn tác.</p>
      </AppModalShell>
    </AppShell>
  )
}
