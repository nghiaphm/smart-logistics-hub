"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { apiClient } from "@/lib/api-client"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { toast } from "@/components/ui/toast"
import { HugeiconsIcon } from "@hugeicons/react"
import { UserIcon, TelephoneIcon } from "@hugeicons/core-free-icons"

export default function CompleteProfilePage() {
  const router = useRouter()
  const [name, setName] = useState("")
  const [phone, setPhone] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) {
      toast.add({
        title: "Thiếu thông tin",
        description: "Vui lòng nhập tên của bạn.",
        type: "error",
      })
      return
    }
    if (!phone.trim()) {
      toast.add({
        title: "Thiếu thông tin",
        description: "Vui lòng nhập số điện thoại của bạn.",
        type: "error",
      })
      return
    }

    setIsSubmitting(true)
    try {
      await apiClient("/profile", {
        method: "POST",
        body: JSON.stringify({
          name: name.trim(),
          phone: phone.trim(),
        }),
      })
      toast.add({
        title: "Thành công",
        description: "Hồ sơ của bạn đã được cập nhật.",
        type: "success",
      })
      router.replace("/modules")
    } catch (err: unknown) {
      console.error(err)
      toast.add({
        title: "Lỗi hệ thống",
        description: err instanceof Error ? err.message : "Không thể tạo hồ sơ. Vui lòng thử lại.",
        type: "error",
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4 dark:bg-zinc-950">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-bold">Hoàn tất hồ sơ</CardTitle>
          <CardDescription>
            Cung cấp thông tin cơ bản để bắt đầu sử dụng Smart Logistics Hub
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="name" className="text-sm font-medium text-neutral-700 dark:text-neutral-300">
                Họ và tên
              </label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-neutral-400">
                  <HugeiconsIcon icon={UserIcon} className="h-4 w-4" />
                </span>
                <Input
                  id="name"
                  type="text"
                  placeholder="Nguyễn Văn A"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="pl-9"
                  disabled={isSubmitting}
                />
              </div>
            </div>

            <div className="space-y-2">
              <label htmlFor="phone" className="text-sm font-medium text-neutral-700 dark:text-neutral-300">
                Số điện thoại
              </label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-neutral-400">
                  <HugeiconsIcon icon={TelephoneIcon} className="h-4 w-4" />
                </span>
                <Input
                  id="phone"
                  type="tel"
                  placeholder="0912345678"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  className="pl-9"
                  disabled={isSubmitting}
                />
              </div>
            </div>
          </CardContent>
          <CardFooter>
            <Button type="submit" className="w-full" disabled={isSubmitting}>
              {isSubmitting ? "Đang xử lý..." : "Lưu và tiếp tục"}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
