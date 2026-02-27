import { createBrowserRouter } from "react-router-dom";
import { AdminLayout } from "@/components/layouts/AdminLayout";
import { AuthGuard } from "@/components/shared/AuthGuard";
import { LoginPage } from "@/pages/LoginPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { UserListPage } from "@/pages/UserListPage";
import { NotFoundPage } from "@/pages/NotFoundPage";

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/",
    element: (
      <AuthGuard>
        <AdminLayout />
      </AuthGuard>
    ),
    children: [
      { index: true, element: <DashboardPage /> },
      { path: "users", element: <UserListPage /> },
    ],
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
]);
