"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { consumeOAuthState, exchangeCodeForTokens, setTokens } from "@/lib/auth";
import { apiClient } from "@/lib/api-client";
import { ApiError } from "@/types/api";

export default function CallbackClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const started = useRef(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (started.current) {
      return;
    }
    started.current = true;

    async function handleCallback() {
      const code = searchParams.get("code");
      const callbackError = searchParams.get("error");
      if (callbackError) {
        setError(`Đăng nhập thất bại: ${callbackError}`);
        return;
      }
      if (!code) {
        setError("Thiếu mã xác thực từ Keycloak.");
        return;
      }
      const savedState = consumeOAuthState();
      const state = searchParams.get("state");
      if (!savedState || !state || state !== savedState) {
        setError("Trạng thái đăng nhập không hợp lệ. Vui lòng thử lại.");
        return;
      }
      try {
        const redirectUri = `${window.location.origin}/auth/callback`;
        const tokens = await exchangeCodeForTokens(code, redirectUri);
        setTokens(tokens);

        // Gọi API đọc profile để phân biệt user mới/cũ
        try {
          await apiClient("/profile");
          router.replace("/modules");
        } catch (apiErr) {
          if (apiErr instanceof ApiError && apiErr.status === 404) {
            router.replace("/auth/complete-profile");
          } else {
            console.error("Lỗi khi tải thông tin hồ sơ:", apiErr);
            setError("Lỗi hệ thống khi xác thực hồ sơ.");
          }
        }
      } catch (err) {
        console.error("Lỗi OAuth callback:", err);
        setError("Không thể hoàn tất đăng nhập. Vui lòng thử lại.");
      }
    }

    void handleCallback();
  }, [router, searchParams]);

  return (
    <main className="flex min-h-screen items-center justify-center">
      <p className="text-muted-foreground">{error ?? "Đang xử lý đăng nhập..."}</p>
    </main>
  );
}
