import { beforeEach, describe, expect, it, vi } from "vitest"

import { apiClient, ApiError } from "@/lib/api-client"
import { ensureFreshAccessToken } from "@/lib/auth"

vi.mock("@/lib/auth", () => ({
  ensureFreshAccessToken: vi.fn(),
}))

const mockFetch = vi.fn()

function jsonResponse(body: unknown, status: number) {
  return new Response(JSON.stringify(body), { status })
}

describe("apiClient", () => {
  beforeEach(() => {
    vi.unstubAllEnvs()
    vi.unstubAllGlobals()
    vi.mocked(ensureFreshAccessToken).mockReset()
    mockFetch.mockReset()
    vi.stubGlobal("fetch", mockFetch)
    vi.stubEnv("NEXT_PUBLIC_API_URL", "http://localhost:8000/api/v1")
  })

  it("gọi API thành công và trả đúng data", async () => {
    vi.mocked(ensureFreshAccessToken).mockResolvedValue("token")
    mockFetch.mockResolvedValue(jsonResponse({ items: [], total: 0 }, 200))

    const data = await apiClient<{ items: unknown[]; total: number }>("/workspaces")

    expect(data).toEqual({ items: [], total: 0 })
    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8000/api/v1/workspaces",
      expect.anything()
    )
  })

  it("đính JWT vào Authorization header khi có token", async () => {
    vi.mocked(ensureFreshAccessToken).mockResolvedValue("jwt-token")
    mockFetch.mockResolvedValue(jsonResponse({}, 200))

    await apiClient("/profile")

    const [, options] = mockFetch.mock.calls[0] as [string, RequestInit]
    expect(new Headers(options.headers).get("Authorization")).toBe("Bearer jwt-token")
  })

  it("không đính Authorization header khi không có token", async () => {
    vi.mocked(ensureFreshAccessToken).mockResolvedValue(null)
    mockFetch.mockResolvedValue(jsonResponse({}, 200))

    await apiClient("/profile")

    const [, options] = mockFetch.mock.calls[0] as [string, RequestInit]
    expect(new Headers(options.headers).get("Authorization")).toBeNull()
  })

  it("ném ApiError đúng format { error: { code, message } } khi API lỗi", async () => {
    vi.mocked(ensureFreshAccessToken).mockResolvedValue("token")
    mockFetch.mockResolvedValue(jsonResponse({ error: { code: 400, message: "Bad Request" } }, 400))

    const error = await apiClient("/workspaces").catch((caught: unknown) => caught)

    expect(error).toBeInstanceOf(ApiError)
    expect((error as ApiError).status).toBe(400)
    expect((error as ApiError).code).toBe(400)
    expect((error as ApiError).message).toBe("Bad Request")
  })

  it("fallback code = HTTP status khi body lỗi không đúng format", async () => {
    vi.mocked(ensureFreshAccessToken).mockResolvedValue(null)
    mockFetch.mockResolvedValue(new Response("không phải json", { status: 500 }))

    const error = await apiClient("/workspaces").catch((caught: unknown) => caught)

    expect(error).toBeInstanceOf(ApiError)
    expect((error as ApiError).status).toBe(500)
    expect((error as ApiError).code).toBe(500)
  })

  it("trả undefined khi response 204", async () => {
    vi.mocked(ensureFreshAccessToken).mockResolvedValue("token")
    mockFetch.mockResolvedValue(new Response(null, { status: 204 }))

    await expect(
      apiClient<void>("/warehouses/1", { method: "DELETE" })
    ).resolves.toBeUndefined()
  })

  it("ném lỗi rõ ràng khi thiếu NEXT_PUBLIC_API_URL và không gọi fetch", async () => {
    vi.unstubAllEnvs()
    vi.mocked(ensureFreshAccessToken).mockResolvedValue("token")

    await expect(apiClient("/workspaces")).rejects.toThrow("NEXT_PUBLIC_API_URL is not set")
    expect(mockFetch).not.toHaveBeenCalled()
  })
})
