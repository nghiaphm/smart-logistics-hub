cách tạo dự án nextjs

**Bước 1: Dọn dẹp thư mục frontend cũ**

Mở PowerShell tại thư mục gốc của dự án (smart-logistic-project), xóa thư mục frontend cũ:



PowerShell

Remove-Item -Path ".\\frontend" -Recurse -Force -ErrorAction SilentlyContinue

**Bước 2: Khởi tạo dự án Next.js mới bằng Yarn**

Chạy lệnh tạo dự án với các cờ (flags) thiết lập sẵn chuẩn Production:



PowerShell

yarn create next-app frontend --typescript --tailwind --eslint --app --src-dir --import-alias "@/\*" --use-yarn

Khi được hỏi về việc tùy chỉnh thêm, bạn chỉ cần để mặc định.



**Bước 3: Xóa thư mục ẩn .git con (nếu có sinh ra)**

Kiểm tra và xóa bỏ file .git bên trong frontend để tránh lỗi submodule:



PowerShell

Remove-Item -Path ".\\frontend\\.git" -Recurse -Force -ErrorAction SilentlyContinue

**Bước 4: Cài đặt các thư viện thiết yếu cho hệ thống Logistics**

Di chuyển vào thư mục frontend và cài đặt các package chuẩn để làm việc với API, Auth và UI:



PowerShell

cd frontend



\# UI Icons \& Tiện ích CSS

yarn add lucide-react clsx tailwind-merge



\# Quản lý State \& Gọi API (Khuyên dùng Axios hoặc TanStack Query)

yarn add axios @tanstack/react-query



\# Form \& Validate dữ liệu

yarn add react-hook-form zod @hookform/resolvers

**Bước 5: Thiết lập cấu trúc thư mục chuẩn (Feature-based Architecture)**

**Một cấu trúc rõ ràng trong frontend/src/ giúp quản lý 10 domain (products, warehouses, orders, drivers...) dễ dàng và tương thích 1:1 với Backend Go:**



Plaintext

frontend/src/

├── app/                      # Next.js App Router (Routing \& Pages)

│   ├── (auth)/               # Route nhóm auth (login, callback)

│   ├── (dashboard)/          # Route nhóm giao diện quản trị chính

│   │   ├── products/         # /products (List, Create, Detail)

│   │   ├── warehouses/       # /warehouses

│   │   ├── inventory/        # /inventory

│   │   ├── orders/           # /orders

│   │   ├── drivers/          # /drivers

│   │   ├── layout.tsx        # Layout chung (Sidebar, Header, Topbar)

│   │   └── page.tsx          # Dashboard overview

│   ├── layout.tsx            # Root layout

│   └── globals.css

├── components/               # UI components dùng chung

│   ├── ui/                   # Button, Input, Table, Modal, Badge...

│   └── layout/               # Sidebar, Navbar, UserNav...

├── features/                 # Logic nghiệp vụ theo từng Domain

│   ├── product/

│   │   ├── api/              # Hàm gọi API (getProducts, createProduct...)

│   │   ├── components/       # ProductTable, ProductFormModal...

│   │   └── types/            # TypeScript interfaces (ProductDTO...)

│   ├── driver/

│   └── order/

├── lib/                      # Cấu hình thư viện dùng chung

│   ├── api-client.ts         # Axios instance có gắn Interceptor tự gán Bearer Token

│   └── utils.ts

└── types/                    # Types hệ thống (API Response, User Role...)

**Bước 6: Cấu hình biến môi trường kết nối Backend \& Keycloak**

Tạo file frontend/.env.local để lưu cấu hình môi trường phát triển:



Code snippet

\# Backend Go API URL

NEXT\_PUBLIC\_API\_BASE\_URL=http://localhost:8000/api/v1



\# Keycloak Configuration

NEXT\_PUBLIC\_KEYCLOAK\_URL=http://localhost:8080

NEXT\_PUBLIC\_KEYCLOAK\_REALM=smart-logistics

NEXT\_PUBLIC\_KEYCLOAK\_CLIENT\_ID=frontend-web

Tạo thêm file frontend/.env.example (file này sẽ commit lên Git để làm mẫu):



Code snippet

NEXT\_PUBLIC\_API\_BASE\_URL=http://localhost:8000/api/v1

NEXT\_PUBLIC\_KEYCLOAK\_URL=http://localhost:8080

NEXT\_PUBLIC\_KEYCLOAK\_REALM=smart-logistics

NEXT\_PUBLIC\_KEYCLOAK\_CLIENT\_ID=frontend-web

**Bước 7: Khởi chạy và kiểm tra**

Chạy server phát triển của Next.js:



PowerShell

yarn dev

Mở trình duyệt truy cập http://localhost:3000. Khi màn hình mặc định của Next.js hiển thị thành công, bạn đã có một nền tảng mã nguồn sạch sẽ, không còn phụ thuộc vào code AI tự sinh trước đó.



ERROR
Received malformed response from registry
Bước 1: Dọn sạch cache và cấu hình lại Registry chuẩn của Yarn
Chạy 2 lệnh sau trong PowerShell để reset cache và trỏ trực tiếp về registry chính thức của npm:

PowerShell
yarn cache clean
yarn config set registry https://registry.npmjs.org/
Bước 2: Dọn dẹp thư mục lỗi trước khi tạo lại
Nếu lệnh trước đó đã tạo dở dang một phần thư mục frontend, hãy xóa sạch để tránh xung đột:

PowerShell
Remove-Item -Path ".\frontend" -Recurse -Force -ErrorAction SilentlyContinue
Bước 3: Khởi tạo lại dự án Next.js bằng npx (Ủy quyền cài đặt qua Yarn)
Sử dụng npx để tải bản template mới nhất của Next.js, kèm cờ --use-yarn để Next.js sử dụng Yarn làm trình quản lý gói:

PowerShell
npx create-next-app@latest frontend --typescript --tailwind --eslint --app --src-dir --import-alias "@/*" --use-yarn