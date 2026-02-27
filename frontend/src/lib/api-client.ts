import axios from "axios";
import { toast } from "sonner";
import { useAuthStore } from "@/stores/auth-store";
import type { ApiResponse } from "@/types/api";
import { ApiError } from "@/types/api";

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "/api/v1",
  timeout: 15000,
});

apiClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse;
    if (body.code !== 0) {
      toast.error(body.message);
      return Promise.reject(new ApiError(body.code, body.message));
    }
    return response;
  },
  (error: unknown) => {
    if (axios.isAxiosError(error)) {
      if (error.response?.status === 401) {
        useAuthStore.getState().logout();
        window.location.href = "/login";
      }
      const msg =
        (error.response?.data as ApiResponse | undefined)?.message ??
        "网络异常，请稍后重试";
      toast.error(msg);
      return Promise.reject(new ApiError(error.response?.status ?? 500, msg));
    }
    return Promise.reject(error);
  },
);

export default apiClient;
