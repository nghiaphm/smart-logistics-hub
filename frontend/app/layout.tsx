import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
// Import Provider vừa tạo từ thư mục components
import AuthProvider from "../components/ui/SessionProvider";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Smart Logistics Hub | AI-Powered Delivery",
  description: "Hệ thống quản lý và điều phối vận tải thông minh",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="vi">
      <body className={inter.className}>
        {/* Bọc AuthProvider bên ngoài children để cung cấp context cho toàn bộ app */}
        <AuthProvider>
          <div className="min-h-screen bg-slate-50">
            {children}
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}