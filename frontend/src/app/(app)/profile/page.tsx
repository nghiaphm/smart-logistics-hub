"use client"

import { useEffect, useState } from "react"
import type { FormEvent } from "react"
import { HugeiconsIcon } from "@hugeicons/react"
import { Alert02Icon, Edit02Icon } from "@hugeicons/core-free-icons"

import { ApiError } from "@/lib/api-client"
import { useProfile, useCreateProfile, useUpdateProfile } from "@/hooks/use-profile"
import { toast } from "@/components/ui/toast"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Skeleton } from "@/components/ui/skeleton"
import { Form, FormActions, FormField } from "@/components/shared/form/Form"

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : "Đã có lỗi xảy ra. Vui lòng thử lại."
}

function isNotFound(error: unknown): boolean {
  return error instanceof ApiError && error.status === 404
}

export default function Page() {
  const { data, isLoading, isError, error, refetch } = useProfile()
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState("")
  const [phone, setPhone] = useState("")
  const [errors, setErrors] = useState<Record<string, string>>({})

  const profile = data
  const profileMissing = isError && isNotFound(error)
  const showForm = editing || profileMissing

  const createMutation = useCreateProfile(() => {
    toast.add({ title: "Đã tạo hồ sơ", type: "success", timeout: 4000 })
    setEditing(false)
  })
  const updateMutation = useUpdateProfile(() => {
    toast.add({ title: "Đã cập nhật hồ sơ", type: "success", timeout: 4000 })
    setEditing(false)
  })
  const isSubmitting = createMutation.isPending || updateMutation.isPending

  useEffect(() => {
    if (isError && error && !profileMissing) {
      toast.add({
        title: "Không thể tải thông tin hồ sơ",
        description: errorMessage(error),
        type: "error",
        timeout: 6000,
      })
    }
  }, [isError, error, profileMissing])

  const startEdit = () => {
    if (profile) {
      setName(profile.name ?? profile.display_name ?? "")
      setPhone(profile.phone ?? "")
    }
    setEditing(true)
  }

  const cancelEdit = () => {
    setErrors({})
    setEditing(false)
  }

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    const nextErrors: Record<string, string> = {}
    if (!name.trim()) nextErrors.name = "Tên không được bỏ trống."
    if (!phone.trim()) nextErrors.phone = "Số điện thoại không được bỏ trống."
    setErrors(nextErrors)
    if (Object.keys(nextErrors).length > 0) {
      toast.add({
        title: "Thông tin không hợp lệ",
        description: "Vui lòng kiểm tra các trường bắt buộc.",
        type: "error",
      })
      return
    }

    const payload = { name: name.trim(), phone: phone.trim() }
    const onError = (err: unknown) => {
      toast.add({
        title: profile ? "Không thể cập nhật hồ sơ" : "Không thể tạo hồ sơ",
        description: errorMessage(err),
        type: "error",
      })
    }

    if (profile) {
      updateMutation.mutate(payload, { onError })
    } else {
      createMutation.mutate(payload, { onError })
    }
  }

  if (isLoading) {
    return (
      <div className="flex max-w-xl flex-col gap-3">
        <Skeleton className="h-28 w-full rounded-2xl" />
        <Skeleton className="h-44 w-full rounded-2xl" />
      </div>
    )
  }

  if (isError && !profileMissing) {
    return (
      <div className="flex flex-col items-center justify-center gap-3 rounded-2xl border border-border bg-card px-6 py-14 text-center">
        <HugeiconsIcon icon={Alert02Icon} className="size-8 text-destructive" />
        <div>
          <p className="font-medium">Không thể tải thông tin hồ sơ</p>
          <p className="mt-1 text-sm text-muted-foreground">{errorMessage(error)}</p>
        </div>
        <Button variant="outline" size="sm" onClick={() => void refetch()}>
          Thử lại
        </Button>
      </div>
    )
  }

  const initials = (profile?.display_name || profile?.user_sub || "U").slice(0, 2).toUpperCase()

  return (
    <div className="flex max-w-xl flex-col gap-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Hồ sơ</h1>
        <p className="mt-1 text-sm text-muted-foreground">Thông tin tài khoản của bạn</p>
      </div>

      {showForm ? (
        <div className="rounded-2xl border border-border bg-card p-5">
          <Form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-4">
              <FormField label="Tên hiển thị" htmlFor="profile-name" error={errors.name} required>
                <Input
                  id="profile-name"
                  value={name}
                  onChange={(event) => setName(event.target.value)}
                  placeholder="Tên của bạn"
                  disabled={isSubmitting}
                />
              </FormField>
              <FormField label="Số điện thoại" htmlFor="profile-phone" error={errors.phone} required>
                <Input
                  id="profile-phone"
                  value={phone}
                  onChange={(event) => setPhone(event.target.value)}
                  placeholder="Số điện thoại"
                  disabled={isSubmitting}
                />
              </FormField>
              <FormActions>
                {!profileMissing ? (
                  <Button type="button" variant="outline" onClick={cancelEdit} disabled={isSubmitting}>
                    Huỷ
                  </Button>
                ) : null}
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? "Đang lưu..." : profile ? "Lưu thay đổi" : "Tạo hồ sơ"}
                </Button>
              </FormActions>
            </div>
          </Form>
        </div>
      ) : (
        <>
          <div className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-4 rounded-2xl border border-border bg-card p-5">
              <span className="flex size-14 shrink-0 items-center justify-center rounded-full bg-sidebar-primary text-lg font-semibold text-sidebar-primary-foreground">
                {initials}
              </span>
              <div className="min-w-0">
                <p className="truncate text-lg font-medium">{profile?.display_name || "Chưa đặt tên"}</p>
                <p className="truncate text-sm text-muted-foreground">{profile?.user_sub}</p>
              </div>
            </div>
            <Button variant="outline" size="sm" className="gap-1" onClick={startEdit}>
              <HugeiconsIcon icon={Edit02Icon} className="size-4" /> Sửa
            </Button>
          </div>
          <div className="rounded-2xl border border-border bg-card p-5">
            <dl className="flex flex-col gap-4 text-sm">
              <div className="flex items-center justify-between gap-4">
                <dt className="text-muted-foreground">Tên hiển thị</dt>
                <dd className="font-medium">{profile?.display_name || "—"}</dd>
              </div>
              <div className="flex items-center justify-between gap-4">
                <dt className="text-muted-foreground">Điện thoại</dt>
                <dd className="font-medium">{profile?.phone || "—"}</dd>
              </div>
              <div className="flex items-center justify-between gap-4">
                <dt className="text-muted-foreground">Tham gia từ</dt>
                <dd className="font-medium">{profile?.created_at || "—"}</dd>
              </div>
            </dl>
          </div>
        </>
      )}
    </div>
  )
}
