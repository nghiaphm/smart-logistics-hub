# 🚚 Smart Logistics Hub (SPX Clone)

Smart Logistics Hub là dự án mô phỏng hệ thống vận hành kho bãi và giao vận hiện đại dựa trên mô hình của Shopee Express (SPX).
Hệ thống được xây dựng theo định hướng thực tế với kiến trúc Backend, AI và Data Engineering nhằm mô phỏng các bài toán logistics quy mô lớn như:

* Quản lý vận đơn và điều phối giao hàng
* Warehouse Management System (WMS)
* Transportation Management System (TMS)
* AI nhận diện biển số xe tải tại cổng kho
* Data Pipeline & Analytics phục vụ phân tích vận hành

---

# 🏗️ System Architecture

Hệ thống được xây dựng theo mô hình **Modular Monolith kết hợp Microservices** nhằm đảm bảo:

* Dễ mở rộng
* Tách biệt domain
* Tối ưu hiệu năng cho AI processing
* Hỗ trợ Data Engineering pipeline

## ⚙️ Tech Stack

### Frontend

* **Next.js**
* Mô phỏng giao diện quản lý vận đơn và vận hành kho của SPX

### Backend

* **FastAPI (Python)**
* API First Design
* RESTful API
* Pydantic Validation

### Database

* **MongoDB**

  * OLTP Database
  * Lưu trữ dữ liệu vận hành:

    * Orders
    * Trips
    * Tracking Logs
    * Inventory
    * AI Events

* **PostgreSQL**

  * Data Warehouse (OLAP)
  * Phân tích dữ liệu doanh thu và hiệu suất vận hành

### AI Module

* **YOLOv8 + CRNN**
* Nhận diện biển số xe tải tại cổng kho
* Vehicle Matching & Tracking

### Data Engineering

* **Apache Spark**

  * ETL Pipeline
  * Data Transformation

* **Apache Airflow**

  * Workflow Scheduling
  * Pipeline Automation

### Security

* **Keycloak**
* Centralized Identity & Access Management

---

# 🧩 Core Business Flow

## 📦 Logistics Flow

Order Created
→ Assign Driver
→ Create Trip
→ Tracking Logs
→ Delivery Completed

---

## 🏭 Warehouse Flow

Inbound
→ Inventory Update
→ Picking
→ Dispatch
→ Outbound

---

## 🤖 AI Gate Flow

Camera Detection
→ Plate Recognition
→ Driver Matching
→ AI Event Created
→ Tracking Updated

---

# 🗄️ Database Design

## Core MongoDB Collections

* `orders`
* `drivers`
* `trips`
* `tracking_logs`
* `inventory`
* `warehouses`
* `inbounds`
* `billing`
* `ai_events`

---

# 🚀 Current Features

* [x] Order Management System
* [x] Tracking Log System
* [x] Warehouse & Inventory Schema
* [x] Transportation & Trip Management
* [x] MongoDB Schema Design
* [x] RESTful API Structure
* [x] API First Backend Architecture

### 🚧 In Progress

* [ ] AI Gate Automation
* [ ] Vehicle Plate Recognition
* [ ] Spark ETL Pipeline
* [ ] Airflow DAG Automation
* [ ] PostgreSQL Data Warehouse
* [ ] Analytics Dashboard

---

# 📌 Roadmap

## v0.1 - Core Logistics Backend

* [x] Orders
* [x] Drivers
* [x] Trips
* [x] Tracking Logs
* [x] Warehouses
* [x] MongoDB Integration

## v0.2 - Warehouse Operations

* [ ] Inventory Management
* [ ] Inbound Operations
* [ ] Inventory Movements

## v0.3 - AI Smart Gate

* [ ] YOLOv8 Integration
* [ ] OCR License Plate Recognition
* [ ] Driver & Vehicle Matching

## v0.4 - Data Engineering Platform

* [ ] Spark ETL
* [ ] Airflow Scheduling
* [ ] PostgreSQL OLAP
* [ ] Analytics Dashboard

---

# 📄 Why MongoDB?

MongoDB được sử dụng cho hệ thống vận hành (OLTP) vì:

* Flexible schema cho logistics domain
* High write throughput
* Hỗ trợ tracking logs & GPS data
* Tối ưu cho nested document structures
* Dễ mở rộng cho event-based systems

---

# 🛠️ Project Structure

```text
.
├── backend/                  # FastAPI Backend
│   ├── app/
│   │   ├── api/              # API Endpoints
│   │   ├── models/           # Pydantic Schemas
│   │   ├── services/         # Business Logic
│   │   ├── repositories/     # Database Layer
│   │   └── core/             # Config & Security
│   └── requirements.txt
│
├── frontend/                 # Next.js Frontend
│
├── ai_module/                # YOLOv8 & OCR Service
│
├── data_pipeline/            # Spark Jobs & Airflow DAGs
│
├── docs/
│   ├── architecture/
│   ├── database/
│   ├── business-flow/
│   └── api/
│
└── .gitignore
```

---

# ⚙️ Getting Started

## Backend

```bash
cd backend
pip install -r requirements.txt
uvicorn app.main:app --reload
```

---

## Frontend

```bash
cd frontend
npm install
npm run dev
```

---

# 🎯 Project Goal

Mục tiêu của dự án là xây dựng mô hình Smart Logistics Hub hiện đại với:

* Warehouse Management System (WMS)
* Transportation Management System (TMS)
* AI-based Logistics Monitoring
* Big Data Analytics Pipeline

Dự án tập trung vào:

* System Design
* Backend Engineering
* AI Integration
* Data Engineering
* Real-world Logistics Workflow

---

# 📚 Future Improvements

* Redis Caching
* Kafka Event Streaming
* Route Optimization
* Real-time GPS Tracking
* AI-based Warehouse Monitoring
* BI Dashboard & Reporting
* Microservices Migration

---
