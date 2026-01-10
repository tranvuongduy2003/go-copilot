import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { UserFilter } from '@/types/user';
import { Search, X } from 'lucide-react';

interface UserFiltersProps {
  filters: UserFilter;
  onFiltersChange: (filters: UserFilter) => void;
  onReset: () => void;
}

export function UserFilters({ filters, onFiltersChange, onReset }: UserFiltersProps) {
  const hasActiveFilters = filters.search || filters.status || filters.roleId;

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onFiltersChange({ ...filters, search: event.target.value, page: 1 });
  };

  const handleStatusChange = (value: string) => {
    const status = value === 'all' ? undefined : (value as UserFilter['status']);
    onFiltersChange({ ...filters, status, page: 1 });
  };

  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-1 items-center gap-2">
        <div className="relative flex-1 sm:max-w-sm">
          <Search className="absolute left-2.5 top-2.5 size-4 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search users..."
            value={filters.search ?? ''}
            onChange={handleSearchChange}
            className="pl-8"
          />
        </div>
        <Select value={filters.status ?? 'all'} onValueChange={handleStatusChange}>
          <SelectTrigger className="w-36">
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="active">Active</SelectItem>
            <SelectItem value="inactive">Inactive</SelectItem>
            <SelectItem value="pending">Pending</SelectItem>
            <SelectItem value="suspended">Suspended</SelectItem>
          </SelectContent>
        </Select>
      </div>
      {hasActiveFilters && (
        <Button variant="ghost" size="sm" onClick={onReset}>
          <X className="mr-2 size-4" />
          Clear filters
        </Button>
      )}
    </div>
  );
}
