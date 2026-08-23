// src/app/layout.tsx
import "@/app/globals.css";
import { Inter } from "next/font/google";
import { AppProviders } from "@/components/providers/app-providers";

const inter = Inter({ subsets: ["latin"] });

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="vi">
      <body className={`${inter.className} bg-neutral-50/50 antialiased dark:bg-neutral-950`}>
        <AppProviders>{children}</AppProviders>
      </body>
    </html>
  );
}