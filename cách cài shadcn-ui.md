Bước 1: Khởi tạo shadcn/uiMở Terminal / PowerShell và di chuyển vào thư mục frontend:PowerShellcd frontend
npx shadcn@latest init
CLI sẽ hỏi một số câu hỏi cấu hình, bạn chọn như sau:Which style would you like to use? $\rightarrow$ Default (hoặc New York tùy sở thích)Which color would you like to use as the base color? $\rightarrow$ Zinc hoặc SlateWould you like to use CSS variables for theming? $\rightarrow$ yes(Quá trình này sẽ tự động cập nhật tailwind.config.ts, src/app/globals.css và tạo file tiện ích src/lib/utils.ts).Bước 2: Cài đặt các Component thiết yếu cho Dashboard LogisticsChạy lệnh thêm các UI component phổ biến dùng cho quản lý kho, bảng dữ liệu và form nhập liệu:PowerShell# Các thành phần giao diện cơ bản & Layout
npx shadcn@latest add button card badge separator

# Thành phần Form & Nhập liệu
npx shadcn@latest add form input select textarea checkbox

# Thành phần Bảng & Dữ liệu
npx shadcn@latest add table

# Thành phần Tương tác & Popup
npx shadcn@latest add dialog dropdown-menu toast
Các file component vừa tải sẽ tự động nằm trong thư mục src/components/ui/ (ví dụ: src/components/ui/button.tsx, src/components/ui/table.tsx...).Bước 3: Cài đặt Icon (Lucide React)shadcn/ui sử dụng thư viện icon lucide-react:PowerShellyarn add lucide-react