import { Button } from '@/components/ui/button';
import { ROUTES } from '@/constants';
import { cn } from '@/lib/utils';
import { useAuthStore } from '@/stores/auth-store';
import { useUIStore } from '@/stores/ui-store';
import { FileText, LayoutDashboard, Settings, Shield, Users, X } from 'lucide-react';
import { NavLink } from 'react-router-dom';

interface NavigationItem {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  permission?: string;
}

const navigationItems: NavigationItem[] = [
  { name: 'Dashboard', href: ROUTES.DASHBOARD, icon: LayoutDashboard },
  { name: 'Users', href: ROUTES.USERS, icon: Users, permission: 'users:read' },
  { name: 'Roles', href: ROUTES.ROLES, icon: Shield, permission: 'roles:read' },
  { name: 'Audit Logs', href: ROUTES.AUDIT_LOGS, icon: FileText, permission: 'audit:read' },
  { name: 'Settings', href: ROUTES.SETTINGS, icon: Settings },
];

export function Sidebar() {
  const { sidebarMobileOpen, setMobileSidebarOpen } = useUIStore();
  const { hasPermission } = useAuthStore();

  const filteredNavigation = navigationItems.filter(
    (item) => !item.permission || hasPermission(item.permission)
  );

  const closeSidebar = () => setMobileSidebarOpen(false);

  return (
    <>
      {sidebarMobileOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={closeSidebar}
          aria-hidden="true"
        />
      )}

      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r bg-background transition-transform duration-300 lg:static lg:translate-x-0',
          sidebarMobileOpen ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex h-14 items-center justify-between border-b px-4 lg:hidden">
          <span className="font-semibold">Menu</span>
          <Button variant="ghost" size="icon" onClick={closeSidebar}>
            <X className="size-5" />
            <span className="sr-only">Close sidebar</span>
          </Button>
        </div>

        <nav className="flex-1 space-y-1 p-4">
          {filteredNavigation.map((item) => (
            <NavLink
              key={item.name}
              to={item.href}
              onClick={closeSidebar}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                  isActive
                    ? 'bg-primary/10 text-primary'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                )
              }
            >
              <item.icon className="size-5" />
              {item.name}
            </NavLink>
          ))}
        </nav>
      </aside>
    </>
  );
}
