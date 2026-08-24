"use client";

import { Button } from "@/components/ui/button";
import { createAuthorizationUrl, saveOAuthState } from "@/lib/auth";

const REDIRECT_PATH = "/auth/callback";

export function SSOLoginButton() {
  function handleLogin() {
    const state = crypto.randomUUID();
    saveOAuthState(state);
    const redirectUri = `${window.location.origin}${REDIRECT_PATH}`;
    window.location.href = createAuthorizationUrl(redirectUri, state);
  }

  return (
    <Button variant="outline" size="lg" onClick={handleLogin}>
      Đăng nhập
    </Button>
  );
}
