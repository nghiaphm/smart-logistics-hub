import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { decodeJwt } from "@/lib/auth";

const AUTH_COOKIE_NAME = "slh_access_token";

function isAuthenticated(request: NextRequest): boolean {
  const token = request.cookies.get(AUTH_COOKIE_NAME)?.value;
  if (!token) {
    return false;
  }
  const payload = decodeJwt<{ exp?: number }>(token);
  if (!payload || typeof payload.exp !== "number") {
    return false;
  }
  return Date.now() < payload.exp * 1000;
}

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  if (pathname === "/") {
    return NextResponse.next();
  }
  if (!isAuthenticated(request)) {
    return NextResponse.redirect(new URL("/", request.url));
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!auth|privacy|terms|api|_next|favicon.ico|robots.txt|sitemap.xml).*)"],
};
