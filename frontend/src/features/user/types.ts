import type { PageParams } from "@/types/api";

export interface User {
  id: string;
  name: string;
  email: string;
  phone: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  phone?: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  phone?: string;
}

export interface UpdateUserStatusRequest {
  status: "active" | "inactive";
}

export interface UserSearchParams extends PageParams {
  keyword?: string;
  status?: string;
}
