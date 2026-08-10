# Refactor Prompt — Smart Logistic Backend (Phase 1: Security + Testing + Error Handling)

## Context

Bạn đang refactor backend Go của project **Smart Logistic Hub** tại `backend/`.

- Framework: Gin v1.12
- DB: MariaDB 11 (`database/sql` + `go-sql-driver/mysql`)
- Auth: Keycloak (JWT RS256 qua `internal/infrastructure/keycloak/verifier.go`)
- Kiến trúc: modular monolith, layering `Handler → Service → Repository → MariaDB`
- 4 domain đã implement đầy đủ: `internal/driver`, `internal/order` (subdirectory style: entity/dto/handler/service/repository), `internal/inventory`, `internal/tracking` (flat style: tất cả file .go ở root domain)
- Entry point: `backend/cmd/api/main.go`
- Middleware có sẵn: `internal/infrastructure/middleware/auth.go` (AuthMiddleware), `rbac.go` (RequireRole), `cors.go`
- Error type dùng chung: `internal/common/errors/errors.go` (`APIError`, sentinel errors: `ErrBadRequest`, `ErrUnauthorized`, `ErrForbidden`, `ErrNotFound`, `ErrConflict`, `ErrInternal`)

**KHÔNG được sửa gì ngoài phạm vi liệt kê bên dưới. Không refactor cấu trúc thư mục (subdirectory vs flat) trong lần này — đó là task riêng.**

---

## Nhiệm vụ

### 1. Fix lỗ hổng auth trên Inventory & Tracking

- Áp `AuthMiddleware` (đã có sẵn, dùng interface `JWTVerifier` trong `middleware/auth.go`) cho toàn bộ route của `internal/inventory` và `internal/tracking`, theo đúng cách `order.RegisterRoutes` và `driver.RegisterRoutes` đang nhận `authMw` làm tham số.
- Cập nhật signature `RegisterRoutes` của 2 domain này để nhận thêm `authMw gin.HandlerFunc`, và cập nhật lời gọi tương ứng trong `main.go`.
- **Quyết định RBAC**: mặc định tất cả route Inventory/Tracking yêu cầu JWT hợp lệ (any authenticated user), riêng `DELETE` yêu cầu role `admin` — theo đúng pattern đã áp dụng cho `DELETE /api/v1/orders/:id`. Dùng `RequireRole("admin")` có sẵn trong `rbac.go`.
- Không đổi route path, không đổi request/response DTO.

### 2. Viết unit test cho Service layer (4 domain FULL)

- Viết test cho `service.go` (hoặc `service/service.go`) của: `driver`, `order`, `inventory`, `tracking`.
- Dùng interface mocking cho Repository layer (tạo mock repository implement cùng interface, không gọi DB thật trong unit test).
- Ưu tiên cover các case:
  - Business logic quan trọng: tạo order (trừ/validate tồn kho nếu có liên quan), cập nhật trạng thái driver, cập nhật/trừ/hoàn `available_qty`/`reserved_qty` trong inventory, ghi tracking event.
  - Input không hợp lệ → trả đúng sentinel error (`ErrBadRequest`, `ErrNotFound`, v.v.)
  - Case not-found → map đúng `ErrNotFound`.
- Đặt file test cạnh file service tương ứng (`service_test.go`), theo convention Go chuẩn (`_test.go`, package `_test` hoặc cùng package tùy theo domain đang dùng gì — kiểm tra domain khác trong repo nếu có để đồng bộ).

### 3. Viết integration test cho Repository layer

- Dùng MariaDB thật (CI đã có sẵn service container MariaDB — dùng chung DSN qua env test).
- Test CRUD cơ bản cho `driver`, `order`, `inventory`, `tracking` repository: insert → get → update → delete, và case `sql.ErrNoRows` → phải trả `ErrNotFound` (kiểm tra đúng theo mục 10 trong `ARCHITECTURE.md`: repository return `apierrors.ErrNotFound` cho `sql.ErrNoRows`).
- Dùng transaction rollback sau mỗi test hoặc truncate table trước mỗi test để đảm bảo test độc lập, không phụ thuộc thứ tự chạy.

### 4. Centralized error-handling middleware

- Tạo middleware Gin mới trong `internal/infrastructure/middleware/error_handler.go`, thay thế logic `resolveError()` hiện đang lặp lại per-handler.
- Middleware nhận error từ `c.Errors` (dùng `c.Error(err)` trong handler thay vì gọi `resolveError()` trực tiếp), extract `APIError.StatusCode` nếu có, fallback 500 cho error không xác định.
- Format response lỗi thống nhất, ví dụ:
  ```json
  { "error": { "code": <status>, "message": "<msg>" } }
  ```
- Đảm bảo error 500 KHÔNG leak message gốc (SQL error, internal detail) ra client — chỉ log server-side (dùng `slog` đã có), trả message generic ("Internal server error") cho client.
- Cập nhật lại 4 domain FULL để dùng `c.Error(err)` + middleware này thay vì gọi `resolveError()` trực tiếp trong từng handler.

---

## Ràng buộc

- Giữ nguyên toàn bộ route path, request/response DTO hiện có — không breaking change cho API contract.
- Không thêm dependency mới ngoài thư viện test chuẩn (`testing`, `net/http/httptest`, có thể dùng `stretchr/testify` cho assertion nếu go.mod đã có sẵn hoặc thêm mới — báo trước nếu cần thêm).
- Code phải pass `go build ./cmd/api/`, `go vet ./...`, `go fmt ./...`, và `go test ./...` (test thật, không skip).
- Không sửa migration SQL, không đổi schema DB.
- Sau khi xong, liệt kê rõ: file nào đã sửa, file nào tạo mới, và bất kỳ quyết định nào cần tôi xác nhận (ví dụ nếu 1 route cần role khác `admin`).

---

## Output mong muốn

1. Diff/danh sách file thay đổi.
2. `go test ./...` output cho thấy test pass.
3. Danh sách case chưa cover được (nếu có) kèm lý do.
