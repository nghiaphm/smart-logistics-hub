"use client";

import { Button } from "@/components/ui/button";
import { createAuthorizationUrl, saveOAuthState } from "@/lib/auth";

const REDIRECT_PATH = "/auth/callback";

interface SSOLoginButtonProps {
  isRegister?: boolean;
}

export function SSOLoginButton({ isRegister = false }: SSOLoginButtonProps) {
  function handleLogin() {
    const state = crypto.randomUUID();
    saveOAuthState(state);
    const redirectUri = `${window.location.origin}${REDIRECT_PATH}`;
    window.location.href = createAuthorizationUrl(redirectUri, state, isRegister);
  }

  return (
    <Button variant={isRegister ? "outline" : "default"} size="sm" onClick={handleLogin}>
      {isRegister ? "Đăng ký" : "Đăng nhập"}
    </Button>
  );
}
