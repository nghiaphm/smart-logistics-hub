import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

import { decodeJwt } from "@/lib/auth";

const AUTH_COOKIE_NAME = "slh_access_token";
const ADMIN_ROLE = "admin";
const ADMIN_PREFIX = "/admin";
const ADMIN_ROUTES = ["/admin", "/admin/users", "/admin/warehouses"];

type JwtClaims = {
  exp?: number;
  realm_access?: { roles?: string[] };
  resource_access?: Record<string, { roles?: string[] }>;
};

function getToken(request: NextRequest): string | null {
  return request.cookies.get(AUTH_COOKIE_NAME)?.value ?? null;
}

function isAuthenticated(request: NextRequest): boolean {
  const token = getToken(request);
  if (!token) {
    return false;
  }
  const payload = decodeJwt<JwtClaims>(token);
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
  const payload = decodeJwt<JwtClaims>(token);
  if (!payload) {
    return false;
  }
  const roles = new Set<string>();
  if (payload.realm_access?.roles) {
    for (const role of payload.realm_access.roles) {
      roles.add(role);
    }
  }
  if (payload.resource_access) {
    for (const client of Object.values(payload.resource_access)) {
      for (const role of client?.roles ?? []) {
        roles.add(role);
      }
    }
  }
  return roles.has(ADMIN_ROLE);
}

function isAdminRoute(pathname: string): boolean {
  return ADMIN_ROUTES.some((route) =>
    route === "/admin"
      ? pathname === route
      : pathname === route || pathname.startsWith(`${route}/`)
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
