import { Button } from '@/components/ui/button';
import { APP_NAME } from '@/constants';
import { useUIStore } from '@/stores/ui-store';
import { Menu, Moon, Sun } from 'lucide-react';
import { Link } from 'react-router-dom';
import { UserMenu } from './user-menu';

export function Header() {
  const { theme, setTheme, toggleSidebar } = useUIStore();

  return (
    <header className="sticky top-0 z-40 border-b bg-background">
      <div className="flex h-14 items-center gap-4 px-4 lg:px-6">
        <Button variant="ghost" size="icon" className="lg:hidden" onClick={toggleSidebar}>
          <Menu className="size-5" />
          <span className="sr-only">Toggle sidebar</span>
        </Button>

        <Link to="/" className="flex items-center gap-2 font-semibold">
          <span className="text-primary">{APP_NAME}</span>
        </Link>

        <div className="flex-1" />

        <Button
          variant="ghost"
          size="icon"
          onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
        >
          {theme === 'dark' ? <Sun className="size-5" /> : <Moon className="size-5" />}
          <span className="sr-only">Toggle theme</span>
        </Button>

        <UserMenu />
      </div>
    </header>
  );
}
