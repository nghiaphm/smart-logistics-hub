"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { consumeOAuthState, exchangeCodeForTokens, setTokens } from "@/lib/auth";

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
        router.replace("/workspaces");
      } catch {
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
