export interface LabelOption {
  label: string;
  value: string;
  variant?: "default" | "secondary" | "destructive" | "outline";
}

export const USER_STATUS_MAP: LabelOption[] = [
  { label: "启用", value: "active", variant: "default" },
  { label: "禁用", value: "inactive", variant: "destructive" },
];

export function getLabelByValue(
  options: LabelOption[],
  value: string,
): string {
  return options.find((o) => o.value === value)?.label ?? value;
}
