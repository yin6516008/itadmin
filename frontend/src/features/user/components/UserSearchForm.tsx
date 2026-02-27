import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Search, RotateCcw } from "lucide-react";
import { USER_STATUS_MAP } from "@/lib/constants";

interface UserSearchFormProps {
  onSearch: (values: { keyword?: string; status?: string }) => void;
}

export function UserSearchForm({ onSearch }: UserSearchFormProps) {
  const [keyword, setKeyword] = useState("");
  const [status, setStatus] = useState<string>("");

  function handleSearch() {
    onSearch({
      keyword: keyword || undefined,
      status: status || undefined,
    });
  }

  function handleReset() {
    setKeyword("");
    setStatus("");
    onSearch({});
  }

  return (
    <div className="mb-4 flex items-center gap-3">
      <Input
        placeholder="搜索姓名 / 邮箱"
        value={keyword}
        onChange={(e) => setKeyword(e.target.value)}
        className="w-52"
        onKeyDown={(e) => e.key === "Enter" && handleSearch()}
      />
      <Select value={status} onValueChange={setStatus}>
        <SelectTrigger className="w-28">
          <SelectValue placeholder="全部状态" />
        </SelectTrigger>
        <SelectContent>
          {USER_STATUS_MAP.map((o) => (
            <SelectItem key={o.value} value={o.value}>
              {o.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Button size="sm" onClick={handleSearch}>
        <Search />
        搜索
      </Button>
      <Button variant="outline" size="sm" onClick={handleReset}>
        <RotateCcw />
        重置
      </Button>
    </div>
  );
}
