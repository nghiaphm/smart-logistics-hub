"use client"

import { createContext, useContext, useSyncExternalStore } from "react"
import type { ReactNode } from "react"
import { useRouter } from "next/navigation"
import { clearTokens, decodeJwt, getAccessToken, isAccessTokenExpired, subscribeTokenChanges } from "@/lib/auth"

type AuthUser = {
  sub?: string
  preferred_username?: string
  realm_access?: { roles?: string[] }
  [key: string]: unknown
}

type AuthContextValue = {
  isAuthenticated: boolean
  user: AuthUser | null
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

let snapshotCache: { token: string; user: AuthUser | null } | null = null

function getAuthSnapshot(): AuthUser | null {
  const token = getAccessToken() ?? ""
  if (snapshotCache && snapshotCache.token === token) {
    return snapshotCache.user
  }
  let user: AuthUser | null = null
  if (token && !isAccessTokenExpired()) {
    user = decodeJwt<AuthUser>(token)
  }
  snapshotCache = { token, user }
  return user
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter()
  const user = useSyncExternalStore(subscribeTokenChanges, getAuthSnapshot, () => null)

  const isAuthenticated = user !== null

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        user,
        logout: () => {
          clearTokens()
          router.push("/")
        },
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return context
}
