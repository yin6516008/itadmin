import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";

interface PageHeaderProps {
  title: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export function PageHeader({ title, action }: PageHeaderProps) {
  return (
    <div className="mb-4 flex items-center justify-between">
      <h2 className="text-lg font-semibold tracking-tight">{title}</h2>
      {action && (
        <Button size="sm" onClick={action.onClick}>
          <Plus className="mr-1 h-4 w-4" />
          {action.label}
        </Button>
      )}
    </div>
  );
}
