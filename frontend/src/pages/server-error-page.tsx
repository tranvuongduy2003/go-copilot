import { Button } from '@/components/ui/button';
import { ROUTES } from '@/constants';
import { AlertTriangle, Home, RefreshCw } from 'lucide-react';
import { Link } from 'react-router-dom';

interface ServerErrorPageProps {
  errorId?: string;
  onRetry?: () => void;
}

export function ServerErrorPage({ errorId, onRetry }: ServerErrorPageProps) {
  const handleRetry = () => {
    if (onRetry) {
      onRetry();
    } else {
      window.location.reload();
    }
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-background px-4">
      <div className="text-center">
        <div className="mb-6 flex justify-center">
          <div className="rounded-full bg-destructive/10 p-4">
            <AlertTriangle className="h-16 w-16 text-destructive" />
          </div>
        </div>

        <h1 className="text-7xl font-bold text-muted-foreground/20">500</h1>
        <h2 className="mt-4 text-2xl font-semibold">Something went wrong</h2>
        <p className="mx-auto mt-2 max-w-md text-muted-foreground">
          We're sorry, but something unexpected happened on our end. Our team has been notified and
          is working to fix the issue.
        </p>

        {errorId && (
          <p className="mt-4 text-sm text-muted-foreground">
            Error Reference: <code className="rounded bg-muted px-2 py-1 font-mono">{errorId}</code>
          </p>
        )}

        <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row sm:justify-center">
          <Button onClick={handleRetry} className="w-full sm:w-auto">
            <RefreshCw className="mr-2 h-4 w-4" />
            Try Again
          </Button>
          <Button variant="outline" asChild className="w-full sm:w-auto">
            <Link to={ROUTES.DASHBOARD}>
              <Home className="mr-2 h-4 w-4" />
              Go to Dashboard
            </Link>
          </Button>
        </div>

        <div className="mt-8 border-t pt-6">
          <p className="text-sm text-muted-foreground">
            If the problem persists, please{' '}
            <a
              href="mailto:support@example.com"
              className="text-primary underline-offset-4 hover:underline"
            >
              contact support
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}
