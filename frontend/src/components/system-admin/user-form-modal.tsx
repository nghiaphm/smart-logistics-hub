"use client"

import { useState } from "react"
import type { FormEvent } from "react"

import { Form, FormField } from "@/components/shared/form/Form"
import { AppModalActions, AppModalShell } from "@/components/shared/modal"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import type { components } from "@/types/api"

type User = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_user_dto.UserResponse"]

export type UserFormValues = {
  username: string
  full_name: string
  email: string
  phone: string
  role: string
  is_active: boolean
}

type UserFormModalProps = {
  open: boolean
  mode: "create" | "edit"
  user?: User
  isSubmitting: boolean
  onOpenChange: (open: boolean) => void
  onSubmit: (values: UserFormValues) => void
}

const roleOptions = [
  { value: "user", label: "Nhân viên" },
  { value: "admin", label: "Quản trị viên" },
]

function toValues(user?: User): UserFormValues {
  return {
    username: user?.username ?? "",
    full_name: user?.full_name ?? "",
    email: user?.email ?? "",
    phone: user?.phone ?? "",
    role: user?.role || "user",
    is_active: user?.is_active ?? true,
  }
}

export function UserFormModal({
  open,
  mode,
  user,
  isSubmitting,
  onOpenChange,
  onSubmit,
}: UserFormModalProps) {
  const [values, setValues] = useState<UserFormValues>(() => toValues(user))
  const [usernameError, setUsernameError] = useState<string | undefined>(undefined)

  const isEdit = mode === "edit"

  const setField = (field: keyof UserFormValues, value: string | boolean) => {
    setValues((prev) => ({ ...prev, [field]: value }))
    if (field === "username") {
      setUsernameError(undefined)
    }
  }

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!isEdit && !values.username.trim()) {
      setUsernameError("Vui lòng nhập tên đăng nhập")
      return
    }
    onSubmit(values)
  }

  return (
    <AppModalShell
      open={open}
      onOpenChange={onOpenChange}
      title={isEdit ? "Sửa người dùng" : "Thêm người dùng"}
      description={
        isEdit ? `Cập nhật thông tin người dùng ${user?.username ?? ""}` : "Nhập thông tin người dùng mới"
      }
      actions={
        <AppModalActions>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isSubmitting}>
            Huỷ
          </Button>
          <Button type="submit" form="user-form" disabled={isSubmitting}>
            {isSubmitting ? "Đang lưu…" : isEdit ? "Lưu thay đổi" : "Thêm người dùng"}
          </Button>
        </AppModalActions>
      }
    >
      <Form id="user-form" onSubmit={handleSubmit}>
        {!isEdit ? (
          <FormField label="Tên đăng nhập" htmlFor="user-username" required error={usernameError}>
            <Input
              id="user-username"
              value={values.username}
              onChange={(event) => setField("username", event.target.value)}
              placeholder="VD: anh.nguyen"
              autoComplete="off"
            />
          </FormField>
        ) : null}
        <FormField label="Họ và tên" htmlFor="user-fullname">
          <Input
            id="user-fullname"
            value={values.full_name}
            onChange={(event) => setField("full_name", event.target.value)}
            placeholder="VD: Nguyễn Văn An"
          />
        </FormField>
        <div className="grid gap-5 sm:grid-cols-2">
          <FormField label="Email" htmlFor="user-email">
            <Input
              id="user-email"
              type="email"
              value={values.email}
              onChange={(event) => setField("email", event.target.value)}
              placeholder="VD: an.nguyen@example.com"
            />
          </FormField>
          <FormField label="Số điện thoại" htmlFor="user-phone">
            <Input
              id="user-phone"
              value={values.phone}
              onChange={(event) => setField("phone", event.target.value)}
              placeholder="VD: 0901234567"
            />
          </FormField>
        </div>
        <FormField label="Vai trò" htmlFor="user-role">
          <Select
            value={values.role}
            onValueChange={(value) => setField("role", value ?? "user")}
          >
            <SelectTrigger className="w-full" id="user-role">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {roleOptions.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </FormField>
        {isEdit ? (
          <label className="flex cursor-pointer items-center gap-2 text-sm text-foreground">
            <Checkbox
              checked={values.is_active}
              onCheckedChange={(checked) => setField("is_active", checked)}
            />
            Tài khoản đang hoạt động
          </label>
        ) : null}
      </Form>
    </AppModalShell>
  )
}
