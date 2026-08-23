# Frontend Architecture — Smart Logistics Hub

Tài liệu mô tả kiến trúc **hiện tại** của frontend (`frontend/`), dựa trên code thực tế trong repo. Không bao gồm đề xuất thay đổi.

---

## 1. Tổng quan hệ thống

Frontend là một ứng dụng **Next.js 16 (App Router)** + **React 19** + **TypeScript (strict)** + **Tailwind CSS v4**, nằm trong thư mục `frontend/` của monorepo, chạy mặc định tại `http://localhost:3000`.

Hệ thống đầy đủ (backend) gồm:

| Phần | Công nghệ | Vị trí |
|---|---|---|
| Frontend | Next.js 16, React 19 | `frontend/` |
| Backend API | Go + Gin, JWT qua Keycloak | `backend/`, `/api/v1`, port `8000` |
| Cơ sở dữ liệu | MariaDB | `docker-compose.yml` |

**Hiện trạng thực tế của frontend:** chỉ có nền tảng nền (styling system, UI kit, provider React Query, trang chủ) được triển khai. Toàn bộ nhánh route nghiệp vụ — `(app)`, `(system-admin)`, `auth`, `api/vnpay`, cùng các thư mục component `account/`, `system-admin/`, `shared/`, `providers/` — tồn tại dưới dạng **file rỗng (stub, 0 byte)**, chưa có logic hay giao tiếp API nào. `next.config.ts` cũng rỗng (dùng cấu hình mặc định).

---

## 2. Entry points

Nên bắt đầu đọc từ đâu để hiểu từng luồng cụ thể (dựa theo cấu trúc thực tế hiện tại):

- **Routing tổng thể:** `src/app/` — cây route của App Router. Bắt đầu từ `src/app/layout.tsx` (root layout) và `src/app/page.tsx` (trang chủ). Các route nghiệp vụ nằm trong route group `(app)` (workspace/logistics), `(system-admin)`, `auth` — hiện là file rỗng.
- **Auth:** `src/lib/auth.ts` (rỗng), `src/app/auth/*` và `src/components/auth/*` (rỗng) — chưa có xác thực nào được triển khai.
- **Data fetching:** `src/app/layout.tsx` — cấu hình `QueryClientProvider` + `staleTime` (điểm nối dữ liệu); `src/components/providers/app-providers.tsx` (rỗng, dành cho việc gom providers). Chưa có hook/fetch thật.
- **UI kit setup:** `src/components/ui/*` (14 component đã triển khai), `src/lib/utils.ts` (`cn`), `src/app/globals.css` (design tokens), `components.json` (cấu hình shadcn/Base UI), dependencies trong `package.json`.
- **Cấu hình dự án:** `tsconfig.json` (path alias `@/*`), `next.config.ts` (rỗng), `postcss.config.mjs`, `eslint.config.mjs`.

---

## 3. Luồng dữ liệu chính (API → UI)

**Hiện tại chưa có luồng dữ liệu thật từ API tới UI.** Không tồn tại trong code:

- Không có lệnh `fetch`/`axios` gọi backend nào.
- Không có hook dữ liệu (`useQuery`, `useMutation`, ...) nào.
- Không có Route Handler hoạt động (`src/app/api/vnpay/route.ts` rỗng).
- Không có `middleware`, Server Action, hay đọc biến môi trường (`NEXT_PUBLIC_*`).

Dữ liệu duy nhất đang hiển thị là **hardcode** trong `src/app/page.tsx` (ví dụ số liệu "12 Kho", "1,420 SKU", "48 Tài xế") — không đi qua server hay API nào.

Các điểm nối hạ tầng đã được dựng sẵn, là nơi luồng API→UI sẽ đi qua khi có triển khai:

```
Backend Go (localhost:8000 /api/v1)
        │  (chưa được gọi từ frontend)
        ▼
QueryClientProvider (root layout, React Query)   ← dùng để đặt fetch/hooks
        ▼
[page.tsx — Server Component] → [Client Component (có "use client")]
        ▼
Thư viện UI (src/components/ui) → Toaster (thông báo lỗi/thành công)
```

- **QueryClientProvider** trong `src/app/layout.tsx` cấp `QueryClient` toàn cục: `staleTime = 5 phút`, `refetchOnWindowFocus = false`.
- **Toaster** (`src/components/ui/toast.tsx`) đã được mount ở root layout, sẵn sàng cho việc hiển thị trạng thái khi có mutation/hook.

Tham chiếu backend (cấu hình tại gốc repo, `.env.development` / `docker-compose.yml`): API `http://localhost:8000`, frontend `http://localhost:3000`, auth Keycloak realm `my_custom_realm` (development). Frontend chưa dùng các giá trị này.

---

## 4. Cấu trúc thư mục và ý nghĩa từng phần

```
frontend/
├── package.json            # Scripts: dev, build, start, lint; packageManager: yarn@1.22.22
├── yarn.lock               # Khóa dependencies (yarn)
├── next.config.ts          # Rỗng — dùng cấu hình Next.js mặc định
├── tsconfig.json           # Strict mode; alias "@/*" → "./src/*"
├── components.json         # Cấu hình shadcn: style "base-maia", icon hugeicons, rsc: true
├── eslint.config.mjs       # eslint-config-next core-web-vitals + typescript
├── postcss.config.mjs      # @tailwindcss/postcss
├── vitest.config.ts        # Rỗng — chưa có test setup
├── Dockerfile              # Rỗng — chưa có image build
├── AGENTS.md / CLAUDE.md   # Hướng dẫn agent (Next.js 16 khác biệt so với tài liệu cũ)
└── src/
    ├── app/                        # App Router
    │   ├── layout.tsx              # Root layout (Server Component): globals.css, font Inter,
    │   │                           #   QueryClientProvider + Toaster toàn app
    │   ├── page.tsx                # Trang chủ (Server Component) — UI tĩnh, dữ liệu hardcode
    │   ├── globals.css             # Tailwind v4 CSS-first + shadcn/tailwind.css + dark mode tokens
    │   ├── (app)/                  # Route group "app đã đăng nhập" — TOÀN BỘ FILE RỖNG
    │   │   ├── [workspace_id]/     #   Layout + logistics/{inbounds,inventory,products,tracking,trips}...
    │   │   ├── profile/            #   Hồ sơ người dùng (Shell, sections...) — rỗng
    │   │   └── workspaces/         #   Danh sách workspace — rỗng
    │   ├── (system-admin)/         # Route group "admin hệ thống" — TOÀN BỘ FILE RỖNG
    │   │   └── admin/              #   Dashboard, logistics, users, warehouses — rỗng
    │   ├── auth/                   # callback, signup, unauthorized, impersonation — rỗng
    │   ├── api/
    │   │   └── vnpay/route.ts      # Route Handler — RỖNG
    │   └── privacy|terms/page.tsx  # Trang tĩnh (rỗng)
    ├── components/
    │   ├── ui/                     # 14 component UI đã triển khai (xem mục 5)
    │   ├── account/                # Workspace shell, sidebar, header — file rỗng
    │   ├── auth/                   # SSOLoginButton, SignUpButton — rỗng
    │   ├── navigation/             # PageRestoreGuard — rỗng
    │   ├── providers/              # app-providers.tsx — RỖNG
    │   ├── shared/                 # Shell, form, modal, table... — rỗng
    │   ├── system-admin/           # Dashboard, tables, layout, packages... — rỗng
    │   ├── app-sidebar.tsx, nav-main.tsx, ... — rỗng
    │   └── data-table.tsx, chart-area-interactive.tsx — rỗng
    └── lib/
        ├── utils.ts                # cn() — hợp nhất className (clsx + tailwind-merge)
        └── auth.ts                 # RỖNG — chưa có logic xác thực phía client
```

### Thành phần đã triển khai (`src/components/ui/`)

| File | Nền tảng | Loại |
|---|---|---|
| `button.tsx` | `@base-ui/react/button` + CVA | Server-compatible (không có `"use client"`) |
| `badge.tsx` | `@base-ui/react/use-render` + CVA | Server-compatible |
| `card.tsx` | HTML thuần + Tailwind | Server-compatible |
| `input.tsx` | `@base-ui/react/input` | Server-compatible |
| `textarea.tsx` | HTML thuần | Server-compatible |
| `skeleton.tsx` | HTML thuần | Server-compatible |
| `checkbox.tsx` | `@base-ui/react/checkbox` | `"use client"` |
| `dialog.tsx` | `@base-ui/react/dialog` | `"use client"` |
| `dropdown-menu.tsx` | `@base-ui/react/menu` | `"use client"` |
| `select.tsx` | `@base-ui/react/select` | `"use client"` |
| `separator.tsx` | `@base-ui/react/separator` | `"use client"` |
| `table.tsx` | HTML thuần | `"use client"` |
| `toast.tsx` | `@base-ui/react/toast` (toast manager) | `"use client"` |
| `tooltip.tsx` | `@base-ui/react/tooltip` | `"use client"` |

---

## 5. Quyết định kiến trúc quan trọng

### 4.1. Server Component là mặc định; Client Component chỉ khi cần tương tác

- **Server Component:** root layout (`layout.tsx`) và trang chủ (`page.tsx`) đều là Server Component (không có chỉ thị `"use client"`).
- **Client Component:** các component UI có tương tác/trạng thái (checkbox, dialog, dropdown-menu, select, separator, table, toast, tooltip) được đánh dấu `"use client"`. Các component trình bày thuần (button, badge, card, input, textarea, skeleton) không có chỉ thị này — chúng chỉ bọc primitives (Base UI/HTML) và được Server Component render trực tiếp.
- Mô hình: **trang/layout ở server → render cây component → các phần tương tác là Client Component riêng biệt.**

### 4.2. Quản lý state: React Query, một QueryClient toàn cục

- **TanStack Query (`@tanstack/react-query`)** là thư viện quản lý state dữ liệu duy nhất được cài.
- `QueryClient` được tạo **inline ngay trong root layout** (`src/app/layout.tsx:8-15`) với:
  - `staleTime: 5 phút` — dữ liệu coi là "mới" trong 5 phút, tránh refetch lại khi điều hướng.
  - `refetchOnWindowFocus: false` — không tự refetch khi quay lại tab.
- Một instance duy nhất phủ toàn app (không tách client theo route). Chưa có cache persistence, chưa có mutation nào.
- Không dùng Redux/Zustand/Context khác (ngoài toast manager của Base UI). File `src/components/providers/app-providers.tsx` (nơi thường gom providers) **đang rỗng** — hiện root layout tự bọc `QueryClientProvider`.

### 4.3. UI kit: shadcn-style nhưng dựa trên Base UI

- Các component `ui/*` được sinh theo cấu hình shadcn (`components.json`: style `"base-maia"`, `rsc: true`), nhưng **nền tảng primitive là Base UI (`@base-ui/react`) thay vì Radix** — mỗi file nhập `{ X as XPrimitive } from "@base-ui/react/x"` rồi thêm `data-slot` + className.
- Variant được quản lý bằng **CVA** (`class-variance-authority`) + gộp className bằng `cn()` (`clsx` + `tailwind-merge`).
- Icon: **Hugeicons** (`@hugeicons/core-free-icons` + `@hugeicons/react`).
- Điểm khác biệt so với shadcn "kinh điển": style `base-maia`, borderRadius bo tròn lớn (`rounded-4xl`, `rounded-2xl`), toast tự dựng theo Base UI toast manager thay vì sonner/radix.

### 4.4. Styling: Tailwind CSS v4 (CSS-first), dark mode, oklch

- `globals.css` dùng cú pháp **Tailwind v4**: `@import "tailwindcss"`, `@theme inline` để map token → CSS variables.
- Dark mode bằng class: `@custom-variant dark (&:is(.dark *))` — chuyển bằng thêm class `.dark` vào `<html>`.
- Toàn bộ bảng màu dùng **oklch**, định nghĩa qua CSS variables (`--background`, `--foreground`, `--sidebar-*`, `--chart-*`, ...) cho cả `:root` và `.dark`.
- `components.json` khai báo shadcn base color `neutral`, `cssVariables: true`, `iconLibrary: hugeicons`, alias `@/components`, `@/ui`, `@/lib`, `@/hooks` (thư mục `hooks/` chưa tồn tại).

### 4.5. Tiêu chuẩn dự án và công cụ

- **TypeScript strict** (`tsconfig.json`) + path alias **`@/*` → `./src/*`** (component import bằng `@/components/...`).
- **Lint:** `eslint` (next/core-web-vitals + typescript). Script: `yarn lint`.
- **Package manager:** `yarn@1.22.22` (khai báo trong `package.json` `packageManager`, có `yarn.lock`).
- **Testing:** chưa có. `vitest.config.ts` rỗng, `package.json` không có script `test`, không có file test nào.
- **Env/config:** `next.config.ts` rỗng (không có proxy, image config, env config). Không có biến `NEXT_PUBLIC_*` nào trong frontend.
- **AGENTS.md** (frontend) cảnh báo: bản Next.js 16 trong repo có API/convention khác tài liệu cũ, phải đọc `node_modules/next/dist/docs/` trước khi viết code.

---

## 6. Kết nối backend

- **Base URL API:** **chưa có.** Không có biến môi trường `NEXT_PUBLIC_*` nào được định nghĩa hay sử dụng trong frontend (`next.config.ts` cũng rỗng). Backend thực tế chạy tại `http://localhost:8000` (`/api/v1`, theo `.env.development` / `docker-compose.yml` ở gốc repo) nhưng frontend chưa đọc cấu hình nào.
- **Format lỗi backend:** chưa có quy ước xử lý lỗi phía frontend. `.kilo/AGENTS.md` hiện **không có** mục "API convention" về định dạng lỗi. Backend có sentinel lỗi chung `internal/common/errors` (theo `README.md`) nhưng chưa có tài liệu quy ước nào dành cho frontend.
- **OpenAPI/Swagger codegen:** **chưa triển khai.** Không có file OpenAPI/Swagger, không có script codegen trong `package.json`.

---

## 7. Điểm cần lưu ý về hiện trạng

- Phần lớn cấu trúc route và component nghiệp vụ là **file rỗng** — chưa nên coi là "đã triển khai" khi đọc cây thư mục.
- Chưa có giao tiếp frontend ↔ backend Go, chưa có xác thực (login/callback/impersonation đều rỗng), chưa có bất kỳ quy ước data fetching nào thực thi.
- `Dockerfile` rỗng nên frontend chưa thể đóng gói container từ repo này.
