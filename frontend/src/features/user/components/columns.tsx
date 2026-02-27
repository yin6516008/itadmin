import type { ColumnDef } from "@tanstack/react-table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { USER_STATUS_MAP } from "@/lib/constants";
import { format } from "date-fns";
import type { User } from "../types";

export function getUserColumns(opts: {
  onEdit: (user: User) => void;
  onDelete: (user: User) => void;
  onToggleStatus: (user: User) => void;
}): ColumnDef<User, unknown>[] {
  return [
    { accessorKey: "name", header: "姓名" },
    { accessorKey: "email", header: "邮箱" },
    { accessorKey: "phone", header: "手机号" },
    {
      accessorKey: "status",
      header: "状态",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        const opt = USER_STATUS_MAP.find((o) => o.value === status);
        return (
          <Badge variant={opt?.variant ?? "default"}>
            {opt?.label ?? status}
          </Badge>
        );
      },
    },
    {
      accessorKey: "created_at",
      header: "创建时间",
      cell: ({ row }) => {
        const val = row.getValue("created_at") as string;
        return val ? format(new Date(val), "yyyy-MM-dd HH:mm:ss") : "-";
      },
    },
    {
      id: "actions",
      header: "操作",
      cell: ({ row }) => (
        <div className="flex gap-1">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => opts.onEdit(row.original)}
          >
            编辑
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => opts.onToggleStatus(row.original)}
          >
            {row.original.status === "active" ? "禁用" : "启用"}
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="text-destructive hover:text-destructive"
            onClick={() => opts.onDelete(row.original)}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];
}
