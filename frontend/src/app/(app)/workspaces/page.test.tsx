import type { ReactNode } from "react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { fireEvent, render, screen, waitFor, within } from "@testing-library/react"
import { beforeEach, describe, expect, it, vi } from "vitest"

import Page from "@/app/(app)/workspaces/page"
import { apiClient } from "@/lib/api-client"
import { Toaster } from "@/components/ui/toast"

vi.mock("@/lib/api-client", () => ({
  apiClient: vi.fn(),
}))

vi.mock("next/link", () => ({
  default: ({ href, children }: { href: string; children: ReactNode }) => (
    <a href={href}>{children}</a>
  ),
}))

type Workspace = {
  id: number
  workspace_code: string
  name: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
  created_by: string
}

function paginated(items: Workspace[]) {
  return { items, total: items.length, limit: 20, skip: 0 }
}

const sampleWorkspaces: Workspace[] = [
  {
    id: 1,
    workspace_code: "WS-001",
    name: "Kho miền Bắc",
    description: "Hà Nội",
    is_active: true,
    created_at: "2026-08-25T00:00:00Z",
    updated_at: "2026-08-25T00:00:00Z",
    created_by: "",
  },
  {
    id: 2,
    workspace_code: "WS-002",
    name: "Kho miền Nam",
    description: "TP.HCM",
    is_active: true,
    created_at: "2026-08-25T00:00:00Z",
    updated_at: "2026-08-25T00:00:00Z",
    created_by: "",
  },
]

describe("workspaces page", () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    })
    vi.mocked(apiClient).mockReset()
  })

  function renderPage(withToaster = false) {
    return render(
      <QueryClientProvider client={queryClient}>
        {withToaster ? <Toaster /> : null}
        <Page />
      </QueryClientProvider>
    )
  }

  it("loading hiển thị 4 skeleton", () => {
    vi.mocked(apiClient).mockReturnValue(new Promise(() => {}))
    const { container } = renderPage()

    expect(container.querySelectorAll('[data-slot="skeleton"]')).toHaveLength(4)
  })

  it("error hiển thị panel với thông báo rõ ràng + nút Thử lại", async () => {
    vi.mocked(apiClient).mockRejectedValue(new Error("Network error"))
    renderPage()

    expect(await screen.findByText("Không thể tải danh sách workspace")).toBeInTheDocument()
    expect(screen.getByText("Network error")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "Thử lại" })).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Thử lại" }))
    await waitFor(() => expect(vi.mocked(apiClient)).toHaveBeenCalledTimes(2))
  })

  it("error cũng hiển thị qua Toaster", async () => {
    vi.mocked(apiClient).mockRejectedValue(new Error("Network error"))
    renderPage(true)

    await waitFor(() => {
      const viewport = document.querySelector('[data-slot="toast-viewport"]')
      expect(viewport).not.toBeNull()
      expect(within(viewport as HTMLElement).getByText("Không thể tải danh sách workspace")).toBeInTheDocument()
    })
  })

  it("empty hiển thị thông báo rõ ràng", async () => {
    vi.mocked(apiClient).mockResolvedValue(paginated([]))
    renderPage()

    expect(await screen.findByText("Chưa có workspace nào")).toBeInTheDocument()
    expect(
      screen.getByText("Workspace sẽ xuất hiện ở đây khi được tạo.")
    ).toBeInTheDocument()
  })

  it("success hiển thị đúng danh sách workspace", async () => {
    vi.mocked(apiClient).mockResolvedValue(paginated(sampleWorkspaces))
    renderPage()

    expect(await screen.findByText("Kho miền Bắc")).toBeInTheDocument()
    expect(screen.getByText("Kho miền Nam")).toBeInTheDocument()
    expect(screen.getByText("WS-001")).toBeInTheDocument()
    expect(screen.getByText("WS-002")).toBeInTheDocument()
    expect(screen.getByRole("link", { name: "Kho miền Bắc" })).toHaveAttribute("href", "/1")
  })
})
