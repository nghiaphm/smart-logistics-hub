import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { decodeJwt } from "@/lib/auth";
import { isSystemAdmin, type JwtUser } from "@/lib/permissions";

const AUTH_COOKIE_NAME = "slh_access_token";
const ADMIN_PREFIX = "/admin";
const ADMIN_ROUTES = ["/admin", "/admin/users", "/admin/warehouses"];
const ADMIN_WORKSPACE_LOGISTICS_RE = /^\/admin\/logistics($|\/)/;

function getToken(request: NextRequest): string | null {
  return request.cookies.get(AUTH_COOKIE_NAME)?.value ?? null;
}

function isAuthenticated(request: NextRequest): boolean {
  const token = getToken(request);
  if (!token) {
    return false;
  }
  const payload = decodeJwt<JwtUser>(token);
  if (!payload || typeof payload.exp !== "number") {
    return false;
  }
  return Date.now() < payload.exp * 1000;
}

function hasAdminRole(request: NextRequest): boolean {
  const token = getToken(request);
  if (!token) {
    return false;
  }
  return isSystemAdmin(decodeJwt<JwtUser>(token));
}

function isAdminRoute(pathname: string): boolean {
  return (
    ADMIN_ROUTES.some((route) =>
      route === "/admin"
        ? pathname === route
        : pathname === route || pathname.startsWith(`${route}/`)
    ) || ADMIN_WORKSPACE_LOGISTICS_RE.test(pathname)
  );
}

export function proxy(request: NextRequest) {
  const pathname = request.nextUrl.pathname.replace(/\/+$/, "") || "/";
  if (pathname === "/") {
    return NextResponse.next();
  }
  if (!isAuthenticated(request)) {
    return NextResponse.redirect(new URL("/", request.url));
  }
  if (pathname.startsWith(ADMIN_PREFIX)) {
    if (!hasAdminRole(request)) {
      return NextResponse.redirect(new URL("/workspaces", request.url));
    }
    if (!isAdminRoute(pathname)) {
      return new NextResponse(null, { status: 404 });
    }
  }
  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!auth|privacy|terms|api|_next|favicon.ico|robots.txt|sitemap.xml).*)"],
};
