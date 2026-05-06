# 🚚 Smart Logistics Hub (SPX Clone)

Dự án mô phỏng hệ thống vận hành kho bãi và giao vận dựa trên mô hình **Shopee Express**, tích hợp AI để tự động hóa quy trình và Big Data để phân tích hiệu suất.

## 🏗️ Kiến trúc hệ thống
Hệ thống được xây dựng theo kiến trúc Microservices hướng tới sự linh hoạt và khả năng mở rộng cao.

- **Frontend:** Next.js (Mô phỏng chính xác giao diện quản lý đơn hàng của SPX).
- **Backend:** FastAPI (Python) - API First Design.
- **Database:** MongoDB (Lưu trữ đơn hàng, vận đơn), PostgreSQL (Data Warehouse).
- **AI Module:** YOLOv8 + CRNN (Nhận diện biển số xe tải ra vào kho).
- **Data Engineering:** - **Apache Spark:** Xử lý và chuyển đổi dữ liệu (ETL).
    - **Apache Airflow:** Điều phối và lập lịch các job dữ liệu.
- **Security:** Keycloak (Centralized Identity Management).

## 🚀 Tính năng hiện tại
- [x] **Quản lý đơn hàng:** Schema chuẩn SPX với thông tin người gửi, người nhận, gói hàng và thanh toán.
- [x] **Tracking System:** Hệ thống lưu vết log đơn hàng theo thời gian thực.
- [x] **API First:** Toàn bộ logic nghiệp vụ được đóng gói qua RESTful API.
- [ ] **AI Gate:** Nhận diện biển số xe tại cổng kho (Đang phát triển).
- [ ] **Data Warehouse:** Chuyển đổi dữ liệu sang OLAP để phân tích doanh thu (Đang phát triển).

## 🛠️ Cấu trúc thư mục
```text
.
├── backend/            # FastAPI Source Code
│   ├── app/
│   │   ├── models/     # Pydantic Schemas (Order, Trip, Warehouse...)
│   │   ├── api/        # API Endpoints
│   │   └── core/       # Cấu hình hệ thống
│   └── requirements.txt
├── frontend/           # Next.js Source Code (React)
├── ai_module/          # YOLOv8 & Plate Recognition
├── data_pipeline/      # Spark Jobs & Airflow DAGs
└── .gitignore
