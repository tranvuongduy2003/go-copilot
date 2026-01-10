import { APP_NAME } from '@/constants';
import { Outlet } from 'react-router-dom';

export function AuthLayout() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-muted/50 p-4">
      <div className="mb-8 text-center">
        <h1 className="text-2xl font-bold text-primary">{APP_NAME}</h1>
        <p className="text-sm text-muted-foreground">Welcome back</p>
      </div>
      <Outlet />
    </div>
  );
}
