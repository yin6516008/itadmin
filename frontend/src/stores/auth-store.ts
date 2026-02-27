import { create } from "zustand";

interface AuthUser {
  id: string;
  name: string;
  avatar: string;
  permissions: string[];
}

interface AuthState {
  token: string | null;
  user: AuthUser | null;
  login: (token: string, user: AuthUser) => void;
  logout: () => void;
  hasPermission: (perm: string) => boolean;
}

const storedToken = localStorage.getItem("token");
const storedUser = localStorage.getItem("user");

export const useAuthStore = create<AuthState>((set, get) => ({
  token: storedToken,
  user: storedUser ? (JSON.parse(storedUser) as AuthUser) : null,

  login: (token: string, user: AuthUser) => {
    localStorage.setItem("token", token);
    localStorage.setItem("user", JSON.stringify(user));
    set({ token, user });
  },

  logout: () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    set({ token: null, user: null });
  },

  hasPermission: (perm: string) => {
    const { user } = get();
    if (!user) return false;
    if (user.permissions.includes("*")) return true;
    return user.permissions.includes(perm);
  },
}));
