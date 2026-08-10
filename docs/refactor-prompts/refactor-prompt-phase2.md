# Refactor Prompt — Smart Logistic Backend (Phase 2: Migration Tooling + Observability)

## Context

Tiếp theo Phase 1 (đã hoàn tất: auth cho Inventory/Tracking, service unit tests, repository integration tests, centralized error middleware — đã verify bằng `go test ./...` pass trên MariaDB 11 thật qua Docker).

Backend Go tại `backend/`:
- Framework: Gin v1.12, DB: MariaDB 11 (`database/sql` + `go-sql-driver/mysql`)
- Migration SQL hiện có tại `backend/migrations/*.sql`, hiện đang apply thủ công (không có tool quản lý version)
- Logging: `slog` (structured), đã có `internal/infrastructure/logger`
- Middleware: `internal/infrastructure/middleware/` (auth.go, rbac.go, cors.go, error_handler.go — vừa thêm ở Phase 1)
- Health/readiness endpoint đã có sẵn (kiểm tra path thật trong `cmd/api/main.go` trước khi giả định)
- `docker-compose.yml` hiện có 4 service: `mariadb`, `backend`, `postgres-keycloak`, `keycloak`

**KHÔNG được sửa gì ngoài phạm vi liệt kê bên dưới. Không đụng tới business logic domain (driver/order/inventory/tracking), không đổi lại những gì Phase 1 vừa làm.**

---

## Nhiệm vụ

### 1. Tích hợp migration tool (golang-migrate)

- Thêm dependency `github.com/golang-migrate/migrate/v4` vào `go.mod` (driver mysql: `github.com/golang-migrate/migrate/v4/database/mysql`).
- Đổi tên file migration hiện có trong `backend/migrations/` sang đúng convention của golang-migrate nếu chưa đúng format: `{version}_{description}.up.sql` và `{version}_{description}.down.sql` (version dạng số tăng dần, ví dụ `000001_initial_schema.up.sql`). Kiểm tra file hiện tại đã đúng format chưa trước khi đổi tên — nếu đã đúng thì giữ nguyên.
- Viết 1 lệnh CLI nhỏ (`cmd/migrate/main.go` hoặc tương tự) hoặc tích hợp migration chạy tự động khi `cmd/api/main.go` khởi động (tùy chọn nào ít rủi ro hơn với hệ thống hiện tại — nếu chọn auto-run lúc start, phải có cách tắt được qua biến env cho môi trường không muốn auto-migrate, ví dụ `AUTO_MIGRATE=false`).
- Cập nhật `docker-compose.yml`: thêm bước chạy migration cho service `backend` (qua entrypoint script hoặc lệnh riêng trước khi start server) HOẶC để `cmd/migrate` chạy như 1 bước riêng trong `depends_on`/init container — chọn cách nào khớp với thiết kế hiện tại, giải thích lý do chọn.
- Cập nhật CI workflow: thêm bước chạy migration trước khi chạy integration test (nếu CI hiện chưa làm việc này mà đang dựa vào schema có sẵn từ đâu đó — kiểm tra lại `.github/workflows/*.yml` hiện tại trước khi sửa).

### 2. Request ID middleware (kích hoạt code có sẵn)

- Trong `internal/infrastructure/middleware/` (hoặc nơi đang định nghĩa `WithRequestID`), kiểm tra implementation hiện tại — nếu đã viết đúng logic nhưng chưa được đăng ký vào middleware chain trong `cmd/api/main.go`, chỉ cần wire vào `r.Use(...)` theo đúng thứ tự hợp lý (trước hoặc sau `ErrorHandler()` — request ID nên có sớm nhất để mọi middleware/log sau đó đều gắn được ID).
- Đảm bảo request ID được đưa vào `slog` context, để mọi dòng log trong 1 request đều có cùng field `request_id`.
- Response header trả về `X-Request-ID` cho client, để debug production dễ hơn khi người dùng report lỗi kèm request ID.

### 3. Metrics cơ bản (Prometheus)

- Thêm `github.com/prometheus/client_golang` vào `go.mod`.
- Expose endpoint `/metrics` (dùng `promhttp.Handler()`), đăng ký route riêng không qua `AuthMiddleware` (endpoint nội bộ cho hệ thống giám sát, không phải API công khai) nhưng cân nhắc: nếu lo ngại bảo mật, có thể tách sang cổng nội bộ riêng — hỏi tôi nếu không chắc hướng nào phù hợp hơn với hệ thống hiện tại.
- Metrics tối thiểu cần có: tổng số HTTP request (counter, gắn label method + path + status code), latency mỗi request (histogram), tổng số lỗi 5xx (có thể derive từ status code label, không cần counter riêng).
- Implement qua 1 middleware Gin mới (`internal/infrastructure/middleware/metrics.go`), áp dụng toàn cục.

### 4. Liveness check tách biệt readiness

- Kiểm tra endpoint health hiện tại trong `cmd/api/main.go` — nếu hiện chỉ có 1 endpoint gộp chung (vừa check DB vừa là liveness), tách thành 2:
  - `/healthz` (liveness) — chỉ xác nhận process đang chạy, KHÔNG check DB, luôn trả 200 nếu server còn sống.
  - `/readyz` (readiness) — check DB connection thật (ping), trả 503 nếu DB không kết nối được. Đây có thể là endpoint hiện tại đã có, chỉ cần đổi path/rename nếu chưa đúng convention này.

---

## Ràng buộc

- Không đổi route path của API nghiệp vụ hiện có (driver/order/inventory/tracking).
- Không thêm dependency ngoài `golang-migrate/migrate/v4` và `prometheus/client_golang` — nếu cần thêm gì khác, báo trước.
- Code phải pass `go build ./cmd/api/`, `go vet ./...`, `go fmt ./...`, `go test ./...` (bao gồm test Phase 1 vẫn phải pass, không được phá vỡ).
- Migration tool phải hoạt động đúng cả khi chạy `docker-compose up -d` (dùng MariaDB thật) — verify bằng cách chạy full stack, không chỉ build thành công.
- Nếu gặp quyết định không rõ (ví dụ: migration tự động lúc start server có rủi ro gì với dữ liệu production không), dừng lại hỏi tôi thay vì tự quyết.

---

## Output mong muốn

1. Diff/danh sách file thay đổi, file tạo mới.
2. `go test ./...` output cho thấy toàn bộ test (Phase 1 + mới) pass.
3. Xác nhận `docker-compose up -d` chạy full stack thành công, migration tự apply đúng (hoặc lệnh CLI migration chạy đúng nếu chọn hướng thủ công).
4. Ví dụ curl thực tế gọi `/metrics`, `/healthz`, `/readyz` và response mẫu.
5. Cập nhật lại `ARCHITECTURE.md` — chỉ phần bị ảnh hưởng (mục 7.1 Migration, mục 11 Observability, mục 15 nếu có vấn đề nào được fix).
6. Danh sách quyết định cần tôi xác nhận (nếu có).
