// import Image from "next/image";

// export default function Home() {
//   return (
//     <div className="flex flex-col flex-1 items-center justify-center bg-zinc-50 font-sans dark:bg-black">
//       <main className="flex flex-1 w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
//         <Image
//           className="dark:invert h-5 w-[100px]"
//           src="/next.svg"
//           alt="Next.js logo"
//           width={100}
//           height={20}
//           priority
//         />
//         <div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
//           <h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight text-black dark:text-zinc-50">
//             To get started, edit the{" "}
//             <code className="rounded bg-black/[.06] px-1.5 py-0.5 font-mono text-[0.9em] dark:bg-white/[.08]">
//               page.tsx
//             </code>{" "}
//             file.
//           </h1>
//           <p className="max-w-md text-lg leading-8 text-zinc-600 dark:text-zinc-400">
//             Looking for a starting point or more instructions? Head over to{" "}
//             <a
//               href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
//               className="font-medium text-zinc-950 dark:text-zinc-50"
//             >
//               Templates
//             </a>{" "}
//             or the{" "}
//             <a
//               href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
//               className="font-medium text-zinc-950 dark:text-zinc-50"
//             >
//               Learning
//             </a>{" "}
//             center.
//           </p>
//         </div>
//         <div className="flex flex-col gap-4 text-base font-medium sm:flex-row">
//           <a
//             className="flex h-12 w-full items-center justify-center gap-2 rounded-full bg-foreground px-5 text-background transition-colors hover:bg-[#383838] dark:hover:bg-[#ccc] md:w-[158px]"
//             href="https://vercel.com/new?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
//             target="_blank"
//             rel="noopener noreferrer"
//           >
//             <Image
//               className="dark:invert h-[14px] w-4"
//               src="/vercel.svg"
//               alt="Vercel logomark"
//               width={16}
//               height={14}
//             />
//             Deploy Now
//           </a>
//           <a
//             className="flex h-12 w-full items-center justify-center rounded-full border border-solid border-black/[.08] px-5 transition-colors hover:border-transparent hover:bg-black/[.04] dark:border-white/[.145] dark:hover:bg-[#1a1a1a] md:w-[158px]"
//             href="https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
//             target="_blank"
//             rel="noopener noreferrer"
//           >
//             Documentation
//           </a>
//         </div>
//       </main>
//     </div>
//   );
// }

// src/app/page.tsx
// src/app/page.tsx
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

import { SSOLoginButton } from "@/components/auth/SSOLoginButton";

import { HugeiconsIcon } from "@hugeicons/react";
import { 
  DeliveryTruck01Icon, 
  PackageIcon, 
  Store01Icon, 
  ArrowRight01Icon 
} from "@hugeicons/core-free-icons";

export default function HomePage() {
  return (
    <main className="min-h-screen bg-neutral-50/50 p-8 dark:bg-neutral-950">
      <div className="mx-auto max-w-5xl space-y-8">

        {/* Header Section */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">Smart Logistics Hub</h1>
            <p className="text-sm text-muted-foreground mt-1">
              Hệ thống điều hành và giám sát kho vận tập trung
            </p>
          </div>
          <div className="flex items-center gap-3">
            <Badge variant="outline" className="gap-1.5 py-1 px-3">
              <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              Hệ thống: Sẵn sàng
            </Badge>
            <div className="flex items-center gap-2">
              <SSOLoginButton />
              <SSOLoginButton isRegister />
            </div>
          </div>
        </div>

        {/* Dashboard Quick Cards */}
        <div className="grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Kho bãi</CardTitle>
              <HugeiconsIcon icon={Store01Icon} className="h-5 w-5 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">12 Kho</div>
              <p className="text-xs text-muted-foreground mt-1">Đang hoạt động trên toàn quốc</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Tồn kho & Mặt hàng</CardTitle>
              <HugeiconsIcon icon={PackageIcon} className="h-5 w-5 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">1,420 SKU</div>
              <p className="text-xs text-muted-foreground mt-1">Được theo dõi tự động</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Đội xe & Vận tải</CardTitle>
              <HugeiconsIcon icon={DeliveryTruck01Icon} className="h-5 w-5 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">48 Tài xế</div>
              <p className="text-xs text-muted-foreground mt-1">Đang thực hiện chuyến đi</p>
            </CardContent>
          </Card>
        </div>

        {/* Action Panel */}
        <Card>
          <CardHeader>
            <CardTitle>Bắt đầu phiên làm việc</CardTitle>
            <CardDescription>
              Chọn phân hệ quản lý để thực hiện các thao tác nhập xuất, điều phối chuyến xe
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-wrap gap-3">
            <Button className="gap-2">
              <HugeiconsIcon icon={PackageIcon} className="h-4 w-4" /> Quản lý sản phẩm
            </Button>
            <Button variant="secondary" className="gap-2">
              <HugeiconsIcon icon={DeliveryTruck01Icon} className="h-4 w-4" /> Điều phối tài xế
            </Button>
            <Button variant="outline" className="gap-2">
              Xem báo cáo AI <HugeiconsIcon icon={ArrowRight01Icon} className="h-4 w-4" />
            </Button>
          </CardContent>
        </Card>

      </div>
    </main>
  );
}