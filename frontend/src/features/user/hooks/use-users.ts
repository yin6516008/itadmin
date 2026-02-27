import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import apiClient from "@/lib/api-client";
import type { ApiResponse, PageResult } from "@/types/api";
import type {
  User,
  UserSearchParams,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateUserStatusRequest,
} from "../types";

const QUERY_KEY = "users";

export function useUsers(params: UserSearchParams) {
  return useQuery({
    queryKey: [QUERY_KEY, params],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<PageResult<User>>>(
        "/users",
        { params },
      );
      return data.data;
    },
  });
}

export function useCreateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (req: CreateUserRequest) => {
      const { data } = await apiClient.post<ApiResponse<User>>("/users", req);
      return data.data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: [QUERY_KEY] }),
  });
}

export function useUpdateUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, ...req }: UpdateUserRequest & { id: string }) => {
      const { data } = await apiClient.put<ApiResponse<User>>(
        `/users/${id}`,
        req,
      );
      return data.data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: [QUERY_KEY] }),
  });
}

export function useUpdateUserStatus() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({
      id,
      ...req
    }: UpdateUserStatusRequest & { id: string }) => {
      const { data } = await apiClient.put<ApiResponse<User>>(
        `/users/${id}/status`,
        req,
      );
      return data.data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: [QUERY_KEY] }),
  });
}

export function useDeleteUser() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(`/users/${id}`);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: [QUERY_KEY] }),
  });
}
