# SKILL SPECIFICATION: FASTAPI TO GOLANG ARCHITECTURE CONVERTER

## 🎯 OBJECTIVE
Chuyển đổi toàn bộ mã nguồn Backend từ Python (FastAPI, Pydantic, Layered Architecture) sang Golang (Gin-Gonic, GORM/Mongo-Go-Driver, Standard Go Layout) nhằm tối ưu hóa hiệu năng, đảm bảo giữ vững logic nghiệp vụ và cấu trúc dữ liệu cũ để giao tiếp mượt mà với Frontend.

---

## 📂 TARGET DIRECTORY MAPPING (THEO DOMAIN-DRIVEN DESIGN)

Khi đọc một file Python ở đường dẫn cũ, Agent phải xác định chính xác file Go kết quả sẽ nằm ở đâu theo từng Domain cụ thể (orders, drivers, inventories...):

| Source File (Python / FastAPI) | Target File (Golang / Gin) | Package Name |
| :--- | :--- | :--- |
| `backend/app/schemas/[domain]/request.py` | `backend-go/internal/schemas/[domain]/request.go` | `package [domain]` |
| `backend/app/schemas/[domain]/response.py` | `backend-go/internal/schemas/[domain]/response.go` | `package [domain]` |
| `backend/app/models/` | `backend-go/internal/models/` | `package models` |
| `backend/app/repositories/` | `backend-go/internal/repositories/` | `package repositories` |
| `backend/app/services/` | `backend-go/internal/services/` | `package services` |
| `backend/app/routers/` | `backend-go/internal/handlers/` | `package handlers` |
| `backend/app/main.py` | `backend-go/cmd/api/main.go` | `package main` |

---

## 🛠️ TECHNICAL CONVERSION RULES (QUY TẮC KỸ THUẬT CỐT LÕI)

### 1. Quản lý Kiểu Dữ Liệu Tùy Chọn (Optional / Nullable)
- **Python:** `Optional[Type] = None` hoặc `Type | None = None`
- **Golang:** Bắt buộc chuyển thành **Con trỏ (Pointer)** `*Type`.
- **JSON Tags:** Các trường Optional phải đi kèm thuộc tính `,omitempty` trong tag `json` (Ví dụ: `json:"province,omitempty"`).

### 2. Kế Thừa Dữ Liệu (Inheritance vs Embedding)
- **Python:** Class con kế thừa class cha (Ví dụ: `class OrderResponse(OrderCreate)`).
- **Golang:** Không dùng kế thừa, sử dụng **Struct Embedding (Nhúng nặc danh)**. Thêm tag `json:",inline"` nếu cần làm phẳng (flatten) dữ liệu JSON đầu ra.

### 3. Xử lý Trùng Tên và Bí Danh (Alias cho MongoDB)
- Khóa chính `_id` từ MongoDB khi lên tầng Schema ứng với trường `id` của Pydantic:
  - Chuyển sang Go: `ID string `json:"_id" bson:"_id"`` để đảm bảo mapping chuẩn xác cả đầu ra API lẫn đầu vào Database.

### 4. Cơ Chế Validate Dữ Liệu (Validation & Hooks)
- Sử dụng công cụ Validator mặc định của Gin thông qua thẻ tag `binding`.
  - `Field(default=...)` hoặc `binding:"required"` cho các trường bắt buộc.
  - Cụm điều kiện tập hợp: Chuyển `description="ENUM"` sang tag `binding:"oneof=VAL1 VAL2 VAL3"`.
- **Validation Hooks (@model_validator):** Tách thành một hàm Method tường minh gắn với Struct (Ví dụ: `func (s *StructName) ExecuteValidators()`). Hàm này bổ khuyết hoặc chuẩn hóa dữ liệu sau khi bind JSON thành công.

### 5. Xử lý Thời Gian (Datetime)
- Toàn bộ các trường định dạng chuỗi thời gian (`created_at: Optional[str]`) khi chuyển sang Go nếu là mốc thời gian hệ thống, ưu tiên đổi thành kiểu dữ liệu `*time.Time` để tối ưu hóa việc lưu trữ và truy vấn.

### 6. Xử lý Đụng độ Package (Package Collision & Aliasing)
Vì các file Schema giờ đây được đặt trong thư mục theo domain (Ví dụ: `schemas/orders`), tên package sẽ là `package orders`. 
Nếu trong quá trình khởi tạo Service hoặc Handler mà Agent cần import cùng lúc nhiều package có tên giống nhau (Ví dụ: `models/orders` và `schemas/orders`), **Agent bắt buộc phải sử dụng Import Alias** để phân biệt.

**Ví dụ định dạng Import chuẩn của Agent:**
```go
import (
    dbModel "backend-go/internal/models"
    reqDto  "backend-go/internal/schemas/orders" // Gán Alias reqDto để tránh nhầm lẫn
)

func ProcessOrder(req reqDto.OrderCreate) error {
    // Logic xử lý...
}
## 🤖 AGENT EXECUTION WORKFLOW (LUỒNG THỰC THI CỦA AI AGENT)

Khi nhận được yêu cầu convert từ người dùng, Agent phải tuân thủ nghiêm ngặt 4 bước sau:

1. **Phân tích mã nguồn nguồn (Analyze Source):** Đọc kỹ file Python đầu vào, xác định các class, các tầng validator, các kiểu dữ liệu và alias (nếu có).
2. **Xác định tọa độ đích (Target Mapping):** Dựa vào bảng thư mục để đưa ra gợi ý đường dẫn file `.go` mới và tên `package` chính xác.
3. **Thực thi chuyển đổi (Execute Code Conversion):** Ánh xạ cú pháp Python/Pydantic sang Golang/Gin tương ứng. Ghi nhận chú thích (comment) rõ ràng ở các đoạn mã xử lý logic phức tạp (như xử lý con trỏ, ép kiểu, hàm tự viết lại).
4. **Đóng gói mã nguồn (Output Deliverable):** Chỉ trả về một block code Go hoàn chỉnh duy nhất, kèm theo giải thích ngắn gọn về các thay đổi quan trọng ở cuối phản hồi.

---

## 📝 EXAMPLES FOR REFERENCE (MẪU ĐỐI CHIẾU CHUẨN)

### Mẫu 1: Trích xuất Request Schema
- **FastAPI Input:**
  ```python
  class Item(BaseModel):
      product_id: Optional[str] = None
      quantity: int = 1