import { beforeEach, describe, expect, it, vi } from "vitest"

import {
  clearTokens,
  consumeOAuthState,
  createAuthorizationUrl,
  decodeJwt,
  ensureFreshAccessToken,
  getAccessToken,
  getAccessTokenExpiry,
  getRefreshToken,
  isAccessTokenExpired,
  isAccessTokenExpiringSoon,
  refreshAccessToken,
  saveOAuthState,
  setTokens,
  subscribeTokenChanges,
} from "@/lib/auth"

const mockFetch = vi.fn()

function base64Url(input: string): string {
  return Buffer.from(input, "utf8")
    .toString("base64")
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/g, "")
}

function makeToken(payload: Record<string, unknown>): string {
  return `${base64Url(JSON.stringify({ alg: "HS256", typ: "JWT" }))}.${base64Url(JSON.stringify(payload))}.sig`
}

function tokenExpiringIn(ms: number): string {
  return makeToken({ exp: Math.floor((Date.now() + ms) / 1000) })
}

describe("auth lib", () => {
  beforeEach(() => {
    vi.unstubAllEnvs()
    vi.unstubAllGlobals()
    window.localStorage.clear()
    window.sessionStorage.clear()
    document.cookie = "slh_access_token=; path=/; max-age=0"
    mockFetch.mockReset()
    vi.stubGlobal("fetch", mockFetch)
    vi.stubEnv("NEXT_PUBLIC_KEYCLOAK_URL", "http://localhost:8180")
    vi.stubEnv("NEXT_PUBLIC_KEYCLOAK_REALM", "smart-logistics")
    vi.stubEnv("NEXT_PUBLIC_KEYCLOAK_CLIENT_ID", "frontend-web")
  })

  describe("lưu/đọc token", () => {
    it("setTokens lưu access + refresh vào localStorage và cookie", () => {
      setTokens({ access_token: "at-1", refresh_token: "rt-1" })
      expect(getAccessToken()).toBe("at-1")
      expect(getRefreshToken()).toBe("rt-1")
      expect(document.cookie).toContain("slh_access_token=at-1")
    })

    it("setTokens không có refresh_token thì không ghi refresh key", () => {
      setTokens({ access_token: "at-2" })
      expect(getAccessToken()).toBe("at-2")
      expect(getRefreshToken()).toBeNull()
    })

    it("clearTokens xoá sạch localStorage và cookie", () => {
      setTokens({ access_token: "at-3", refresh_token: "rt-3" })
      clearTokens()
      expect(getAccessToken()).toBeNull()
      expect(getRefreshToken()).toBeNull()
      expect(document.cookie).not.toContain("slh_access_token")
    })

    it("subscribeTokenChanges nhận thông báo khi setTokens/clearTokens và hủy được", () => {
      const listener = vi.fn()
      const unsubscribe = subscribeTokenChanges(listener)
      setTokens({ access_token: "at-4" })
      expect(listener).toHaveBeenCalledTimes(1)
      clearTokens()
      expect(listener).toHaveBeenCalledTimes(2)
      unsubscribe()
      setTokens({ access_token: "at-5" })
      expect(listener).toHaveBeenCalledTimes(2)
    })
  })

  describe("decodeJwt và hết hạn", () => {
    it("decodeJwt decode đúng payload kể cả UTF-8", () => {
      expect(
        decodeJwt<{ sub: string; name: string }>(makeToken({ sub: "u1", name: "Nguyễn Văn An" }))
      ).toEqual({ sub: "u1", name: "Nguyễn Văn An" })
    })

    it("decodeJwt trả null khi token không hợp lệ", () => {
      expect(decodeJwt("không-phải-jwt")).toBeNull()
    })

    it("getAccessTokenExpiry trả exp dạng ms, null khi không có token", () => {
      expect(getAccessTokenExpiry()).toBeNull()
      setTokens({ access_token: makeToken({ exp: 1_800_000_000 }) })
      expect(getAccessTokenExpiry()).toBe(1_800_000_000_000)
    })

    it("isAccessTokenExpired: true khi không có token hoặc quá hạn, false khi còn hạn", () => {
      expect(isAccessTokenExpired()).toBe(true)
      setTokens({ access_token: makeToken({ exp: Math.floor((Date.now() - 1000) / 1000) }) })
      expect(isAccessTokenExpired()).toBe(true)
      setTokens({ access_token: tokenExpiringIn(600_000) })
      expect(isAccessTokenExpired()).toBe(false)
    })

    it("isAccessTokenExpiringSoon: true khi trong buffer 60s, false khi còn xa", () => {
      setTokens({ access_token: tokenExpiringIn(10_000) })
      expect(isAccessTokenExpiringSoon()).toBe(true)
      setTokens({ access_token: tokenExpiringIn(600_000) })
      expect(isAccessTokenExpiringSoon()).toBe(false)
    })
  })

  describe("refresh và ensureFreshAccessToken", () => {
    it("refreshAccessToken trả null khi không có refresh_token (chưa đăng nhập)", async () => {
      await expect(refreshAccessToken()).resolves.toBeNull()
      expect(mockFetch).not.toHaveBeenCalled()
    })

    it("refreshAccessToken gọi token endpoint và lưu token mới khi thành công", async () => {
      setTokens({ access_token: "old-at", refresh_token: "old-rt" })
      mockFetch.mockResolvedValue(
        new Response(
          JSON.stringify({ access_token: "new-at", refresh_token: "new-rt", expires_in: 300 }),
          { status: 200 }
        )
      )

      const token = await refreshAccessToken()

      expect(token).toBe("new-at")
      expect(getAccessToken()).toBe("new-at")
      expect(getRefreshToken()).toBe("new-rt")
      expect(mockFetch).toHaveBeenCalledTimes(1)
      const [, options] = mockFetch.mock.calls[0] as [string, RequestInit]
      expect(options.method).toBe("POST")
      const body = options.body as URLSearchParams
      expect(body.get("grant_type")).toBe("refresh_token")
      expect(body.get("refresh_token")).toBe("old-rt")
      expect(body.get("client_id")).toBe("frontend-web")
    })

    it("refreshAccessToken xoá token và trả null khi thất bại", async () => {
      setTokens({ access_token: "old-at", refresh_token: "old-rt" })
      mockFetch.mockResolvedValue(new Response("lỗi", { status: 400 }))

      await expect(refreshAccessToken()).resolves.toBeNull()
      expect(getAccessToken()).toBeNull()
      expect(getRefreshToken()).toBeNull()
    })

    it("ensureFreshAccessToken trả null khi chưa đăng nhập, không gọi fetch", async () => {
      await expect(ensureFreshAccessToken()).resolves.toBeNull()
      expect(mockFetch).not.toHaveBeenCalled()
    })

    it("ensureFreshAccessToken trả token hiện tại khi còn hạn, không gọi fetch", async () => {
      const token = tokenExpiringIn(600_000)
      setTokens({ access_token: token })
      await expect(ensureFreshAccessToken()).resolves.toBe(token)
      expect(mockFetch).not.toHaveBeenCalled()
    })

    it("ensureFreshAccessToken refresh khi token sắp hết hạn", async () => {
      setTokens({ access_token: tokenExpiringIn(10_000), refresh_token: "old-rt" })
      mockFetch.mockResolvedValue(
        new Response(JSON.stringify({ access_token: "renewed-at", refresh_token: "new-rt" }), {
          status: 200,
        })
      )

      await expect(ensureFreshAccessToken()).resolves.toBe("renewed-at")
      expect(getAccessToken()).toBe("renewed-at")
    })
  })

  describe("OAuth state và login URL", () => {
    it("save/consume state dùng 1 lần", () => {
      saveOAuthState("state-1")
      expect(consumeOAuthState()).toBe("state-1")
      expect(consumeOAuthState()).toBeNull()
    })

    it("createAuthorizationUrl build đúng URL login", () => {
      const url = createAuthorizationUrl("http://localhost:3000/auth/callback", "state-2")
      expect(url).toContain(
        "http://localhost:8180/realms/smart-logistics/protocol/openid-connect/auth?"
      )
      expect(url).toContain("client_id=frontend-web")
      expect(url).toContain("redirect_uri=")
      expect(url).toContain("state=state-2")
    })
  })
})
