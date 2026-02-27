import { LayoutDashboard, Users } from "lucide-react";
import type { LucideIcon } from "lucide-react";

export interface NavItem {
  title: string;
  url: string;
  icon?: LucideIcon;
}

export const navMain: NavItem[] = [
  { title: "工作台", url: "/", icon: LayoutDashboard },
  { title: "用户管理", url: "/users", icon: Users },
];
