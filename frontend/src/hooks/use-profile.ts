"use client"

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"

import { apiClient } from "@/lib/api-client"
import type { components } from "@/types/api"

type Profile = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.ProfileResponse"]
type CreateProfileRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.CreateProfileRequest"]
type UpdateProfileRequest = components["schemas"]["my-web-app_com_smart-logistic-hub_internal_profile_dto.UpdateProfileRequest"]

export function useProfile() {
  return useQuery({
    queryKey: ["profile"],
    queryFn: () => apiClient<Profile>("/profile"),
  })
}

export function useCreateProfile(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateProfileRequest) =>
      apiClient<Profile>("/profile", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["profile"] })
      onSuccess?.()
    },
  })
}

export function useUpdateProfile(onSuccess?: () => void) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: UpdateProfileRequest) =>
      apiClient<Profile>("/profile", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["profile"] })
      onSuccess?.()
    },
  })
}
