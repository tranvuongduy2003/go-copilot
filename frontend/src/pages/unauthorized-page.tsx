import { Button } from '@/components/ui/button';
import { ROUTES } from '@/constants';
import { ShieldX } from 'lucide-react';
import { Link } from 'react-router-dom';

interface UnauthorizedPageProps {
  requiredPermission?: string;
}

export function UnauthorizedPage({ requiredPermission }: UnauthorizedPageProps) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background px-4">
      <div className="text-center">
        <div className="mb-6 flex justify-center">
          <div className="rounded-full bg-destructive/10 p-4">
            <ShieldX className="h-16 w-16 text-destructive" />
          </div>
        </div>

        <h1 className="text-7xl font-bold text-muted-foreground/20">403</h1>
        <h2 className="mt-4 text-2xl font-semibold">Access Denied</h2>
        <p className="mx-auto mt-2 max-w-md text-muted-foreground">
          You don't have permission to access this page. Please contact your administrator if you
          believe this is a mistake.
        </p>

        {requiredPermission && (
          <p className="mt-4 text-sm text-muted-foreground">
            Required permission:{' '}
            <code className="rounded bg-muted px-2 py-1 font-mono">{requiredPermission}</code>
          </p>
        )}

        <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
          <Button asChild className="w-full sm:w-auto">
            <Link to={ROUTES.DASHBOARD}>Go to Dashboard</Link>
          </Button>
          <Button
            variant="outline"
            onClick={() => window.history.back()}
            className="w-full sm:w-auto"
          >
            Go Back
          </Button>
        </div>
      </div>
    </div>
  );
}
