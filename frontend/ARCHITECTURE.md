# Frontend Architecture — Smart Logistics Hub

Tài liệu mô tả kiến trúc **hiện tại** của frontend (`frontend/`), dựa trên code thực tế trong repo. Không bao gồm đề xuất thay đổi.

---

## 1. Tổng quan hệ thống

Frontend là một ứng dụng **Next.js 16 (App Router)** + **React 19** + **TypeScript (strict)** + **Tailwind CSS v4**, nằm trong thư mục `frontend/` của monorepo, chạy mặc định tại `http://localhost:3000`.

Hệ thống đầy đủ (backend) gồm:

| Phần | Công nghệ | Vị trí |
| --- | --- | --- |
| Frontend | Next.js 16, React 19 | `frontend/` |
| Backend API | Go + Gin, JWT qua Keycloak | `backend/`, `/api/v1`, port `8000` |
| Cơ sở dữ liệu | MariaDB | `docker-compose.yml` |

**Hiện trạng thực tế của frontend:** chỉ có nền tảng nền (styling system, UI kit, provider React Query, trang chủ) được triển khai. Toàn bộ nhánh route nghiệp vụ — `(app)`, `(system-admin)`, `auth`, `api/vnpay`, cùng các thư mục component `account/`, `system-admin/`, `shared/`, `providers/` — tồn tại dưới dạng **file rỗng (stub, 0 byte)**, chưa có logic hay giao tiếp API nào. `next.config.ts` cũng rỗng (dùng cấu hình mặc định).

---

## 2. Entry points

Nên bắt đầu đọc từ đâu để hiểu từng luồng cụ thể (dựa theo cấu trúc thực tế hiện tại):

- **Routing tổng thể:** `src/app/` — cây route của App Router. Bắt đầu từ `src/app/layout.tsx` (root layout) và `src/app/page.tsx` (trang chủ). Các route nghiệp vụ nằm trong route group `(app)` (workspace/logistics), `(system-admin)`, `auth` — hiện là file rỗng.
- **Auth (đã triển khai Giai đoạn 2):** bắt đầu từ `src/components/auth/SSOLoginButton.tsx` (nút "Đăng nhập", mount trên `src/app/page.tsx`) → `src/lib/auth.ts` (`createAuthorizationUrl` build URL login Keycloak — realm `smart-logistics`, client `frontend-web`) → Keycloak redirect về `src/app/auth/callback/` (`exchangeCodeForTokens` đổi code lấy token, `setTokens` lưu localStorage + cookie) → redirect vào (app)/. `src/proxy.ts` (Next.js 16 proxy) chặn (app)/ + (system-admin)/ khi chưa có cookie hợp lệ. `src/contexts/user.context.tsx` vẫn rỗng.
- **Data fetching:** `src/components/providers/app-providers.tsx` (client provider: tạo `QueryClient` + mount `Toaster`) — điểm nối dữ liệu; `src/contexts/query-provider.tsx` (rỗng); `src/hooks/**` — toàn bộ hook là stub rỗng, chưa có hook/fetch thật.
- **Cấu hình & trạng thái toàn cục:** `src/config/` (roles.ts — rỗng), `src/contexts/` (user, workspace, warehouse, warehouse-mode, lang, toast, finance... — toàn bộ rỗng), `src/hooks/` (useTheme, usePermission, use-mobile, finance, admin, logistic... — toàn bộ rỗng).
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
AppProviders (client provider — QueryClient + Toaster)   ← dùng để đặt fetch/hooks
        ▼
[page.tsx — Server Component] → [Client Component (có "use client")]
        ▼
Thư viện UI (src/components/ui) → Toaster (thông báo lỗi/thành công)
```

- **AppProviders** (`src/components/providers/app-providers.tsx`, client component) tạo `QueryClient` toàn cục qua `useState`: `staleTime = 5 phút`, `refetchOnWindowFocus = false`, và mount `Toaster`. Root layout (`src/app/layout.tsx`, Server Component) chỉ render `<AppProviders>{children}</AppProviders>`.
- **Toaster** (`src/components/ui/toast.tsx`) được mount bên trong `AppProviders`, sẵn sàng cho việc hiển thị trạng thái khi có mutation/hook.
- **Contexts & hooks:** `src/contexts/*` (user, workspace, warehouse, lang, toast, finance, query-provider...) và `src/hooks/**` (useTheme, usePermission, finance, admin, logistic...) đều là stub rỗng — chưa có hook/context nào đưa dữ liệu API vào UI.

Tham chiếu backend (cấu hình tại gốc repo, `.env.development` / `docker-compose.yml`): API `http://localhost:8000`, frontend `http://localhost:3000`, auth Keycloak realm `smart-logistics` (development). Frontend chưa dùng các giá trị này.

---

## 4. Cấu trúc thư mục và ý nghĩa từng phần

```
frontend/
├── package.json            # Scripts: dev, build, start, lint, generate:api; packageManager: yarn@1.22.22
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
    │   │                           #   chỉ render <AppProviders>
    │   ├── page.tsx                # Trang chủ (Server Component) — UI tĩnh, có nút đăng nhập SSOLoginButton
    │   ├── globals.css             # Tailwind v4 CSS-first + shadcn/tailwind.css + dark mode tokens
    │   ├── (app)/                  # Route group "app đã đăng nhập" — layout [workspace_id] ĐÃ TRIỂN KHAI
    │   │   ├── [workspace_id]/     #   Layout shell (sidebar + header); logistics/page.tsx demo shared components
    │   │   ├── profile/            #   Hồ sơ người dùng (Shell, sections...) — rỗng
    │   │   └── workspaces/         #   Danh sách workspace — rỗng
    │   ├── (system-admin)/         # Route group "admin hệ thống" — TOÀN BỘ FILE RỖNG
    │   │   └── admin/              #   Dashboard, logistics, users, warehouses — rỗng
    │   ├── auth/                   # callback ĐÃ TRIỂN KHAI (nhận redirect + lưu token); signup, unauthorized, impersonation — rỗng
    │   ├── api/
    │   │   └── vnpay/route.ts      # Route Handler — RỖNG
    │   └── privacy|terms/page.tsx  # Trang tĩnh (rỗng)
    ├── components/
    │   ├── ui/                     # 14 component UI đã triển khai (xem mục 5)
    │   ├── account/                # Workspace shell, sidebar, header — file rỗng
    │   ├── auth/                   # SSOLoginButton ĐÃ TRIỂN KHAI (redirect Keycloak); SignUpButton — rỗng
    │   ├── navigation/             # PageRestoreGuard — rỗng
    │   ├── providers/              # app-providers.tsx — ĐÃ TRIỂN KHAI (QueryClient + AuthProvider + Toaster)
    │   ├── shared/                 # ĐÃ TRIỂN KHAI: AppShell, form/Form, modal/AppModalShell+Actions, DataTable; còn stub khác rỗng
    │   ├── system-admin/           # Dashboard, tables, layout, packages... — rỗng
    │   ├── app-sidebar.tsx, nav-main.tsx, app-header.tsx — ĐÃ TRIỂN KHAI (shell layout)
    │   └── data-table.tsx, chart-area-interactive.tsx — rỗng
    ├── config/                     # roles.ts, roles.test.ts — toàn bộ RỖNG (stub)
    ├── contexts/                   # auth.context.tsx ĐÃ TRIỂN KHAI (AuthProvider + useAuth);
    │                               #   user, workspace, warehouse, warehouse-mode, lang, toast,
    │                               #   finance-*, query-provider... — toàn bộ RỖNG (stub)
    ├── hooks/
    │   ├── useTheme.ts, usePermission.ts, use-mobile.ts, use-is-breakpoint.ts,
    │   │   useTableRowSelection.ts, useFinance*, useClientListPageShell,
    │   │   useWarehouseToast, useVisibilityAwareInterval... — RỖNG (stub)
    │   ├── admin/                  # useAdminViewMode.ts — RỖNG
    │   └── logistic/               # useInboundProcess, useInventoryStock, useOutboundDispatch,
    │                               #   useTripTracking, useYoloCameraStream — RỖNG
    ├── lib/
    │   ├── utils.ts                # cn() — hợp nhất className (clsx + tailwind-merge)
    │   ├── auth.ts                 # ĐÃ TRIỂN KHAI — token JWT Keycloak (lưu/đọc/refresh/exchange, OIDC state)
    │   └── api-client.ts           # ĐÃ TRIỂN KHAI — apiClient<T>(fetch, base URL, Authorization header, ApiError)
    ├── proxy.ts                    # ĐÃ TRIỂN KHAI — Next.js 16 proxy bảo vệ route (app)/ + (system-admin)/
    └── types/                      # api.ts (re-export ApiError + type sinh) + api-generated.ts (codegen OpenAPI)
```

### Thành phần đã triển khai (`src/components/ui/`)

| File | Nền tảng | Loại |
| --- | --- | --- |
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

### 5.1. Server Component là mặc định; Client Component chỉ khi cần tương tác

- **Server Component:** root layout (`layout.tsx`) và trang chủ (`page.tsx`) đều là Server Component (không có chỉ thị `"use client"`).
- **Client Component:** các component UI có tương tác/trạng thái (checkbox, dialog, dropdown-menu, select, separator, table, toast, tooltip) được đánh dấu `"use client"`. Các component trình bày thuần (button, badge, card, input, textarea, skeleton) không có chỉ thị này — chúng chỉ bọc primitives (Base UI/HTML) và được Server Component render trực tiếp.
- Mô hình: **trang/layout ở server → render cây component → các phần tương tác là Client Component riêng biệt.**

### 5.2. Quản lý state: React Query, một QueryClient toàn cục

- **TanStack Query (`@tanstack/react-query`)** là thư viện quản lý state dữ liệu duy nhất được cài.
- `QueryClient` được tạo **bên trong client provider** `src/components/providers/app-providers.tsx` (dùng `useState(() => new QueryClient(...))`) với:
  - `staleTime: 5 phút` — dữ liệu coi là "mới" trong 5 phút, tránh refetch lại khi điều hướng.
  - `refetchOnWindowFocus: false` — không tự refetch khi quay lại tab.
- Root layout (`src/app/layout.tsx`) là Server Component, chỉ render `<AppProviders>`. QueryClient không được tạo ở server để tránh vượt ranh giới Server→Client.
- Một instance duy nhất phủ toàn app (không tách client theo route). Chưa có cache persistence, chưa có mutation nào.
- Không dùng Redux/Zustand. Thư mục `src/contexts/` có sẵn các file context stub (user, workspace, warehouse, lang, toast, finance, query-provider...) nhưng **toàn bộ rỗng** — chưa có Context nào được triển khai ngoài provider ở `app-providers.tsx`.

### 5.3. UI kit: shadcn-style nhưng dựa trên Base UI

- Các component `ui/*` được sinh theo cấu hình shadcn (`components.json`: style `"base-maia"`, `rsc: true`), nhưng **nền tảng primitive là Base UI (`@base-ui/react`) thay vì Radix** — mỗi file nhập `{ X as XPrimitive } from "@base-ui/react/x"` rồi thêm `data-slot` + className.
- Variant được quản lý bằng **CVA** (`class-variance-authority`) + gộp className bằng `cn()` (`clsx` + `tailwind-merge`).
- Icon: **Hugeicons** (`@hugeicons/core-free-icons` + `@hugeicons/react`).
- Điểm khác biệt so với shadcn "kinh điển": style `base-maia`, borderRadius bo tròn lớn (`rounded-4xl`, `rounded-2xl`), toast tự dựng theo Base UI toast manager thay vì sonner/radix.

### 5.4. Styling: Tailwind CSS v4 (CSS-first), dark mode, oklch

- `globals.css` dùng cú pháp **Tailwind v4**: `@import "tailwindcss"`, `@theme inline` để map token → CSS variables.
- Dark mode bằng class: `@custom-variant dark (&:is(.dark *))` — chuyển bằng thêm class `.dark` vào `<html>`.
- Toàn bộ bảng màu dùng **oklch**, định nghĩa qua CSS variables (`--background`, `--foreground`, `--sidebar-*`, `--chart-*`, ...) cho cả `:root` và `.dark`.
- `components.json` khai báo shadcn base color `neutral`, `cssVariables: true`, `iconLibrary: hugeicons`, alias `@/components`, `@/ui`, `@/lib`, `@/hooks` (thư mục `hooks/` tồn tại nhưng toàn bộ file rỗng).

### 5.5. Tiêu chuẩn dự án và công cụ

- **TypeScript strict** (`tsconfig.json`) + path alias **`@/*` → `./src/*`** (component import bằng `@/components/...`).
- **Lint:** `eslint` (next/core-web-vitals + typescript). Script: `yarn lint`.
- **Package manager:** `yarn@1.22.22` (khai báo trong `package.json` `packageManager`, có `yarn.lock`).
- **Testing:** chưa có. `vitest.config.ts` rỗng, `package.json` không có script `test`. File `src/config/roles.test.ts` tồn tại nhưng là stub rỗng.
- **Env/config:** `next.config.ts` rỗng (không có proxy, image config, env config). Không có biến `NEXT_PUBLIC_*` nào trong frontend.
- **AGENTS.md** (frontend) cảnh báo: bản Next.js 16 trong repo có API/convention khác tài liệu cũ, phải đọc `node_modules/next/dist/docs/` trước khi viết code.

### 5.6. Xác thực: Keycloak OIDC, token ở client, proxy bảo vệ route

- **Keycloak OIDC authorization code flow**, **public client** `frontend-web` (không cần client secret); realm `smart-logistics`, server `http://localhost:8180` (cấu hình qua `NEXT_PUBLIC_KEYCLOAK_URL/REALM/CLIENT_ID` trong `.env.development`).
- Token JWT lưu **localStorage** (`access_token`, `refresh_token`) — client SPA không có secret. `lib/auth.ts` quản lý lưu/đọc/refresh; `ensureFreshAccessToken()` tự refresh qua token endpoint Keycloak khi token sắp hết hạn (buffer 60s) — `api-client.ts` gọi trước mỗi request để tự đính `Authorization: Bearer <token>`.
- Proxy/server không đọc được localStorage → `setTokens`/`clearTokens` đồng thời set/xoá cookie `slh_access_token` (không HttpOnly — chỉ phục vụ chặn route UX, xác thực thật vẫn qua Authorization header ở backend). Cookie được cập nhật mỗi lần refresh.
- **Chống CSRF (OIDC state):** `SSOLoginButton` tạo `state` (`crypto.randomUUID`) lưu `sessionStorage`; callback validate qua `consumeOAuthState` (dùng 1 lần).
- **Bảo vệ route:** Next.js 16 đổi tên `middleware` → **`proxy.ts`** (Node runtime). Matcher bảo vệ (app)/, (system-admin)/ và `[workspace_id]`; giữ public `/`, `/auth/*`, `/terms`, `/privacy`, `/api/*`; kiểm tra cookie `slh_access_token` + `exp` của JWT — thiếu/hết hạn → redirect `/`.

---

## 6. Kết nối backend

- **Base URL API:** đã cấu hình qua biến `NEXT_PUBLIC_API_URL=http://localhost:8000/api/v1` trong `frontend/.env.development` (Next.js tự load file này khi chạy `next dev`). Chưa có `.env.production` — khi build production cần cung cấp `NEXT_PUBLIC_API_URL` tại thời điểm build, nếu không client bundle sẽ inline giá trị `undefined`.
- **API client:** `frontend/src/lib/api-client.ts` cung cấp `apiClient<T>(path, options)` dùng `fetch`, base URL lấy từ `NEXT_PUBLIC_API_URL`, tự đính header `Authorization: Bearer <token>` nếu có token (đọc từ `localStorage` key `access_token`), ném `ApiError` khi response không ok, và ném lỗi rõ ràng nếu thiếu `NEXT_PUBLIC_API_URL`. Hàm trả về body JSON đã parse (không bọc trong `{ data: ... }`).
- **Type dùng chung:** `frontend/src/types/api.ts` re-export toàn bộ type sinh tự động
  (`src/types/api-generated.ts` — `components`, `paths`) và giữ class runtime `ApiError`
  (`status`, `code`, `message`). Không còn `ApiResponse<T>` viết tay — dùng schema
  `PaginatedResponse` generate theo từng module.
- **Format lỗi backend:** backend trả lỗi theo `{ error: { code, message } }` với `code` là HTTP status (nguồn: ErrorHandler trong `backend/internal/infrastructure/middleware/error_handler.go`; auth middleware cũng đã chuẩn hoá về format này). `apiClient` parse đúng format này, fallback `code` = HTTP status nếu body lỗi không đúng format.
- **Đã kiểm chứng:** pipeline `env → api-client → hiển thị` đã test thành công trong dev (gọi `GET /api/v1/warehouses`, render dữ liệu thật từ backend). Luồng đăng nhập lấy token chưa implement — chờ Giai đoạn 2 (auth Keycloak).
- **OpenAPI/Swagger codegen:** đã triển khai. Backend sinh spec `backend/docs/swagger.json`
  (swagger 2.0, từ annotation swag). Frontend cài dev deps `openapi-typescript` +
  `swagger2openapi`, script `yarn generate:api` (`frontend/scripts/generate-api.mjs`) đọc spec,
  convert sang OpenAPI 3.0 rồi generate ra `src/types/api-generated.ts` (24 paths / 53 schemas).

---

## 7. Điểm cần lưu ý về hiện trạng

- Phần lớn cấu trúc route, component, cũng như các thư mục mới `config/`, `contexts/`, `hooks/` (gồm `hooks/admin`, `hooks/logistic`) và `lib/auth.ts` là **file rỗng (stub)** — chưa nên coi là "đã triển khai" khi đọc cây thư mục.
- Chưa có giao tiếp frontend ↔ backend Go, chưa có xác thực (login/callback/impersonation đều rỗng), chưa có bất kỳ quy ước data fetching nào thực thi.
- `Dockerfile` rỗng nên frontend chưa thể đóng gói container từ repo này.
