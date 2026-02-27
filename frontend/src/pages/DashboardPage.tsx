import { PageHeader } from "@/components/shared/PageHeader";
import { Users, UserCheck, UserX } from "lucide-react";

export function DashboardPage() {
  return (
    <>
      <PageHeader title="工作台" />
      <div className="grid gap-4 md:grid-cols-3">
        <StatCard icon={Users} label="用户总数" value={0} />
        <StatCard
          icon={UserCheck}
          label="已启用"
          value={0}
          className="text-green-600"
        />
        <StatCard
          icon={UserX}
          label="已禁用"
          value={0}
          className="text-red-600"
        />
      </div>
    </>
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
  className,
}: {
  icon: React.ElementType;
  label: string;
  value: number;
  className?: string;
}) {
  return (
    <div className="rounded-lg border bg-white p-6">
      <div className="flex items-center gap-3">
        <Icon className={`h-5 w-5 text-muted-foreground ${className ?? ""}`} />
        <span className="text-sm text-muted-foreground">{label}</span>
      </div>
      <p className={`mt-2 text-2xl font-bold ${className ?? ""}`}>{value}</p>
    </div>
  );
}
