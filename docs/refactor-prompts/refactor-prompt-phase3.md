# Refactor Prompt — Smart Logistic Backend (Phase 3: Structure Consistency + Config/Secrets + Docs Housekeeping)

## Context

Tiếp theo Phase 1 (auth, testing, error handling) và Phase 2 (migration tooling, observability) — cả hai đã hoàn tất, verify qua CI xanh.

Backend Go tại `backend/`. Từ audit ban đầu (`ARCHITECTURE.md`), các vấn đề còn tồn đọng thuộc nhóm "code hygiene", không phải bug chức năng:

- `driver`, `order` dùng cấu trúc thư mục subdirectory: `entity/`, `dto/`, `handler/`, `service/`, `repository/`
- `inventory`, `tracking` dùng cấu trúc flat: tất cả file `.go` nằm thẳng ở root domain (`inventory/handler.go`, `inventory/service.go`...)
- File orphaned nghi ngờ không còn dùng: `pkg/utils/logger.go`, `internal/infrastructure/keycloak/client.go`, `internal/auth/handler/auth_handler.go` (chỉ có `package handlers`, không có code)
- Có domain chỉ tồn tại entity/DTO với struct tag `bson` (leftover từ MongoDB cũ), chưa migrate sang MariaDB — cần xác định domain nào (kiểm tra thực tế trong `internal/`, đừng giả định tên cụ thể vì có thể đã thay đổi từ lúc audit)
- `.env.production` hiện đang rỗng
- `README.md` mô tả kiến trúc cũ (Python/FastAPI/MongoDB), không khớp code Go hiện tại

**KHÔNG được sửa gì ngoài phạm vi liệt kê bên dưới. Không đụng business logic, không đụng những gì Phase 1/2 vừa làm (auth, test, error middleware, migration, observability).**

---

## Nhiệm vụ

### 1. Thống nhất cấu trúc thư mục domain

- Chuyển `inventory` và `tracking` sang cấu trúc subdirectory giống `driver`/`order`: tách file hiện có thành `entity/`, `dto/`, `handler/`, `service/`, `repository/` tương ứng.
- Đây là thay đổi thuần túy về tổ chức file — package name, import path nội bộ, và toàn bộ logic/test đã viết ở Phase 1 phải giữ nguyên hành vi. Chạy `go test ./...` ngay sau khi chuyển để đảm bảo không có gì bị vỡ do đổi cấu trúc.
- Cập nhật lại phần liên quan trong `ARCHITECTURE.md` nếu có mô tả cấu trúc thư mục cụ thể.

### 2. Dọn code chết / file orphaned

- Kiểm tra thực tế (không giả định) các file nghi ngờ orphaned từ audit: `pkg/utils/logger.go`, `internal/infrastructure/keycloak/client.go`, `internal/auth/handler/auth_handler.go`. Với mỗi file, chạy kiểm tra import thực tế (`grep`/tìm reference) trước khi xóa — chỉ xóa nếu xác nhận không còn nơi nào import/dùng tới. Nếu file nào vẫn đang được dùng, giữ nguyên, ghi rõ trong báo cáo lý do không xóa.
- Với các domain chỉ có entity/DTO dùng `bson` tag (leftover MongoDB): liệt kê domain thực tế tìm thấy trong `internal/`, hỏi tôi xác nhận domain nào giữ (chuyển sang tag SQL-compatible, chuẩn bị implement sau) và domain nào xóa hẳn — không tự quyết định xóa hàng loạt.

### 3. `.env.production` và `.gitignore`

- Điền `.env.production` với đầy đủ key cần thiết (tham khảo cấu trúc từ `.env.development`/`.env.example` hiện có), nhưng **dùng placeholder rõ ràng** cho giá trị nhạy cảm (ví dụ `CHANGE_ME_IN_PRODUCTION`), không điền secret thật vào file commit.
- Kiểm tra `.gitignore` đã chặn đúng các file `.env*` chứa secret thật chưa (nên giữ lại `.env.example` không bị ignore, nhưng chặn `.env`, `.env.production` nếu nó từng chứa giá trị thật trước đây — kiểm tra lịch sử git nếu cần, báo tôi nếu phát hiện secret thật từng bị commit).

### 4. Cập nhật `README.md`

- Viết lại phần mô tả kiến trúc cho khớp thực tế hiện tại: Go + Gin + MariaDB + Keycloak, không còn Python/FastAPI/MongoDB.
- Bao gồm: hướng dẫn chạy local (`docker-compose up -d`, chạy migration qua `cmd/migrate`), cấu trúc thư mục domain, các endpoint quan trọng (`/healthz`, `/readyz`, `/metrics`), cách chạy test (`go test ./...`, cần MariaDB cho integration test).
- Không cần viết dài dòng — README chỉ cần đủ để người mới join hiểu cách chạy project trong 5 phút đầu.

### 5. Thư mục stub rỗng (nếu có)

- Kiểm tra các thư mục như `ai_service/`, `data_pipeline/`, `agents/` (nếu tồn tại thực tế trong repo — xác nhận trước khi giả định). Nếu rỗng hoặc chỉ có placeholder, thêm ghi chú ngắn trong README nói rõ đây là chỗ dành cho roadmap tương lai, hoặc xóa nếu không có kế hoạch gần — hỏi tôi nếu không chắc hướng nào.

---

## Ràng buộc

- Không đổi route path, không đổi request/response DTO, không đổi migration SQL.
- Sau khi đổi cấu trúc thư mục (Task 1), bắt buộc `go build ./cmd/api/`, `go vet ./...`, `go test ./...` pass 100% — coi đây là gate trước khi làm tiếp Task 2-5.
- Không điền secret thật vào bất kỳ file nào sẽ commit.
- Nếu gặp domain/file không chắc chắn (đặc biệt Task 2 - xóa gì, Task 3 - domain bson nào giữ/xóa), dừng lại hỏi tôi.

---

## Output mong muốn

1. Diff/danh sách file thay đổi, file xóa, file tạo mới — đặc biệt rõ ràng cho Task 1 (đổi cấu trúc) vì đây là thay đổi lan rộng nhiều file.
2. `go test ./...` output xác nhận pass sau Task 1 và sau khi hoàn tất toàn bộ.
3. Danh sách file đã xóa ở Task 2, kèm bằng chứng đã kiểm tra không còn reference (ví dụ kết quả `grep`).
4. Domain bson tìm thấy thực tế + quyết định giữ/xóa (chờ tôi xác nhận nếu cần).
5. README.md mới (nội dung đầy đủ trong báo cáo hoặc link file).
6. Danh sách quyết định cần tôi xác nhận.
