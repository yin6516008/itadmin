import { useMemo, useState } from "react";
import { toast } from "sonner";
import { PageHeader } from "@/components/shared/PageHeader";
import { DataTable } from "@/components/shared/DataTable";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { UserSearchForm } from "@/features/user/components/UserSearchForm";
import { UserFormDialog } from "@/features/user/components/UserFormDialog";
import { getUserColumns } from "@/features/user/components/columns";
import {
  useUsers,
  useDeleteUser,
  useUpdateUserStatus,
} from "@/features/user/hooks/use-users";
import type { User, UserSearchParams } from "@/features/user/types";

export function UserListPage() {
  const [searchParams, setSearchParams] = useState<UserSearchParams>({
    page: 1,
    size: 20,
  });
  const { data, isLoading } = useUsers(searchParams);

  const [formOpen, setFormOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const [deleteTarget, setDeleteTarget] = useState<User | null>(null);
  const [statusTarget, setStatusTarget] = useState<User | null>(null);

  const deleteUser = useDeleteUser();
  const updateStatus = useUpdateUserStatus();

  function handleSearch(values: { keyword?: string; status?: string }) {
    setSearchParams((prev) => ({ ...prev, ...values, page: 1 }));
  }

  function handleCreate() {
    setEditingUser(null);
    setFormOpen(true);
  }

  function handleEdit(user: User) {
    setEditingUser(user);
    setFormOpen(true);
  }

  async function handleConfirmDelete() {
    if (!deleteTarget) return;
    try {
      await deleteUser.mutateAsync(deleteTarget.id);
      toast.success("删除成功");
    } catch {
      // handled by interceptor
    }
    setDeleteTarget(null);
  }

  async function handleConfirmToggleStatus() {
    if (!statusTarget) return;
    const nextStatus =
      statusTarget.status === "active" ? "inactive" : "active";
    const label = nextStatus === "active" ? "启用" : "禁用";
    try {
      await updateStatus.mutateAsync({ id: statusTarget.id, status: nextStatus });
      toast.success(`${label}成功`);
    } catch {
      // handled by interceptor
    }
    setStatusTarget(null);
  }

  const columns = useMemo(
    () =>
      getUserColumns({
        onEdit: handleEdit,
        onDelete: setDeleteTarget,
        onToggleStatus: setStatusTarget,
      }),
    [],
  );

  return (
    <>
      <PageHeader
        title="用户管理"
        action={{ label: "新建用户", onClick: handleCreate }}
      />
      <UserSearchForm onSearch={handleSearch} />
      <DataTable
        columns={columns}
        data={data?.list ?? []}
        total={data?.total}
        page={searchParams.page}
        pageSize={searchParams.size}
        loading={isLoading}
        onPageChange={(page) =>
          setSearchParams((prev) => ({ ...prev, page }))
        }
      />

      <UserFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        user={editingUser}
      />

      {/* 删除确认 */}
      <AlertDialog
        open={!!deleteTarget}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>确认删除</AlertDialogTitle>
            <AlertDialogDescription>
              确定要删除用户「{deleteTarget?.name}」吗？此操作不可撤销。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleConfirmDelete}
              className="bg-destructive text-white hover:bg-destructive/90"
            >
              删除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 状态切换确认 */}
      <AlertDialog
        open={!!statusTarget}
        onOpenChange={(open) => !open && setStatusTarget(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              确认{statusTarget?.status === "active" ? "禁用" : "启用"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              确定要
              {statusTarget?.status === "active" ? "禁用" : "启用"}
              用户「{statusTarget?.name}」吗？
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmToggleStatus}>
              确定
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
