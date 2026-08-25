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
import { UserFormModal } from "@/components/system-admin/user-form-modal"
import type { UserFormValues } from "@/components/system-admin/user-form-modal"

type User = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse"]
type UserList = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.PaginatedResponse"]
type CreateUserRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.CreateUserRequest"]
type UpdateUserRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UpdateUserRequest"]

type ModalState =
  | { type: "form"; mode: "create" | "edit"; user?: User }
  | { type: "delete"; user: User }
  | null

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

const roleLabels: Record<string, string> = {
  admin: "Quản trị viên",
  user: "Nhân viên",
}

function roleLabel(role: string | undefined): string {
  return (role && roleLabels[role]) || role || "—"
}

const columns: Column<User>[] = [
  {
    key: "username",
    header: "Tên đăng nhập",
    cell: (user) => <span className="font-medium">{user.username}</span>,
  },
  { key: "fullname", header: "Họ và tên", cell: (user) => user.full_name || "—" },
  {
    key: "email",
    header: "Email",
    cell: (user) => <span className="line-clamp-1 max-w-56">{user.email || "—"}</span>,
  },
  {
    key: "phone",
    header: "Số điện thoại",
    cell: (user) => user.phone || "—",
    className: "text-muted-foreground",
  },
  {
    key: "role",
    header: "Vai trò",
    cell: (user) => (
      <Badge variant={user.role === "admin" ? "default" : "outline"}>
        {roleLabel(user.role)}
      </Badge>
    ),
  },
  {
    key: "status",
    header: "Trạng thái",
    cell: (user) => (
      <Badge variant={user.is_active ? "default" : "outline"}>
        {user.is_active ? "Hoạt động" : "Ngừng"}
      </Badge>
    ),
  },
]

export default function AdminUsersPage() {
  const queryClient = useQueryClient()
  const [modal, setModal] = useState<ModalState>(null)

  const listQuery = useQuery({
    queryKey: ["admin", "users"],
    queryFn: () => apiClient<UserList>("/users?limit=100"),
  })

  useEffect(() => {
    if (listQuery.isError && listQuery.error) {
      toast.add({
        title: "Không thể tải danh sách người dùng",
        description: errorMessage(listQuery.error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [listQuery.isError, listQuery.error])

  const invalidateList = () => {
    queryClient.invalidateQueries({ queryKey: ["admin", "users"] })
  }

  const createMutation = useMutation({
    mutationFn: (values: UserFormValues) => {
      const payload: CreateUserRequest = {
        username: values.username.trim(),
        full_name: values.full_name.trim() || undefined,
        email: values.email.trim() || undefined,
        phone: values.phone.trim() || undefined,
        role: values.role,
      }
      return apiClient<User>("/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
    },
    onSuccess: () => {
      toast.add({ title: "Đã thêm người dùng", type: "success", timeout: 4000 })
      invalidateList()
      setModal(null)
    },
    onError: (error) => {
      toast.add({
        title: "Không thể thêm người dùng",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, values }: { id: number; values: UserFormValues }) => {
      const payload: UpdateUserRequest = {
        full_name: values.full_name.trim() || undefined,
        email: values.email.trim() || undefined,
        phone: values.phone.trim() || undefined,
        role: values.role,
        is_active: values.is_active,
      }
      return apiClient<User>(`/users/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      })
    },
    onSuccess: () => {
      toast.add({ title: "Đã cập nhật người dùng", type: "success", timeout: 4000 })
      invalidateList()
      setModal(null)
    },
    onError: (error) => {
      toast.add({
        title: "Không thể cập nhật người dùng",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => apiClient<void>(`/users/${id}`, { method: "DELETE" }),
    onSuccess: () => {
      toast.add({ title: "Đã xoá người dùng", type: "success", timeout: 4000 })
      invalidateList()
      setModal(null)
    },
    onError: (error) => {
      toast.add({
        title: "Không thể xoá người dùng",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    },
  })

  const deleteTarget = modal?.type === "delete" ? modal.user : null
  const users = listQuery.data?.items ?? []

  return (
    <AppShell
      title="Người dùng"
      description="Quản lý tài khoản người dùng trên toàn hệ thống"
      actions={
        <Button onClick={() => setModal({ type: "form", mode: "create" })}>
          <HugeiconsIcon icon={Add01Icon} />
          Thêm người dùng
        </Button>
      }
    >
      {listQuery.isError ? (
        <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-border bg-card px-6 py-14 text-center">
          <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
          <div>
            <p className="font-medium">Không thể tải danh sách người dùng</p>
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
              cell: (user) => (
                <div className="flex items-center justify-end gap-1">
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label="Sửa người dùng"
                    onClick={() => setModal({ type: "form", mode: "edit", user })}
                  >
                    <HugeiconsIcon icon={Edit02Icon} className="size-4" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon-sm"
                    aria-label="Xoá người dùng"
                    className="text-destructive hover:text-destructive"
                    onClick={() => setModal({ type: "delete", user })}
                  >
                    <HugeiconsIcon icon={Delete01Icon} className="size-4" />
                  </Button>
                </div>
              ),
            },
          ]}
          rows={users}
          rowKey={(user) => user.id ?? user.username ?? ""}
          loading={listQuery.isLoading}
          emptyText="Chưa có người dùng nào. Bấm “Thêm người dùng” để tạo mới."
        />
      )}

      {modal?.type === "form" ? (
        <UserFormModal
          key={modal.mode === "create" ? "create" : `edit-${modal.user?.id ?? "new"}`}
          open
          mode={modal.mode}
          user={modal.user}
          isSubmitting={createMutation.isPending || updateMutation.isPending}
          onOpenChange={(open) => {
            if (!open) {
              setModal(null)
            }
          }}
          onSubmit={(values) => {
            if (modal.mode === "create") {
              createMutation.mutate(values)
            } else if (modal.user?.id != null) {
              updateMutation.mutate({ id: modal.user.id, values })
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
        title="Xoá người dùng"
        description={
          deleteTarget
            ? `Xoá người dùng ${deleteTarget.username}${deleteTarget.full_name ? ` — ${deleteTarget.full_name}` : ""}?`
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
