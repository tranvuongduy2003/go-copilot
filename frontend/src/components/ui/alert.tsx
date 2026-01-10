import * as React from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { AlertCircle, CheckCircle, Info, AlertTriangle, X } from 'lucide-react';

const alertVariants = cva(
  'relative w-full rounded-lg border px-4 py-3 text-sm grid gap-1 [&>svg]:size-4 [&>svg]:translate-y-0.5 [&>svg+div]:translate-y-[-2px] [&>svg~*]:pl-7 [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-3.5',
  {
    variants: {
      variant: {
        default: 'bg-background text-foreground',
        info: 'border-blue-200 bg-blue-50 text-blue-900 [&>svg]:text-blue-600 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-100',
        success:
          'border-green-200 bg-green-50 text-green-900 [&>svg]:text-green-600 dark:border-green-800 dark:bg-green-950 dark:text-green-100',
        warning:
          'border-yellow-200 bg-yellow-50 text-yellow-900 [&>svg]:text-yellow-600 dark:border-yellow-800 dark:bg-yellow-950 dark:text-yellow-100',
        destructive:
          'border-red-200 bg-red-50 text-red-900 [&>svg]:text-red-600 dark:border-red-800 dark:bg-red-950 dark:text-red-100',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
);

const alertIconMap = {
  default: null,
  info: Info,
  success: CheckCircle,
  warning: AlertTriangle,
  destructive: AlertCircle,
};

interface AlertProps extends React.ComponentProps<'div'>, VariantProps<typeof alertVariants> {
  dismissible?: boolean;
  onDismiss?: () => void;
}

function Alert({
  className,
  variant = 'default',
  dismissible,
  onDismiss,
  children,
  ...props
}: AlertProps) {
  const Icon = alertIconMap[variant || 'default'];

  return (
    <div
      data-slot="alert"
      role="alert"
      className={cn(alertVariants({ variant }), className)}
      {...props}
    >
      {Icon && <Icon />}
      {children}
      {dismissible && (
        <button
          type="button"
          onClick={onDismiss}
          className="absolute right-2 top-2 rounded-md p-1 opacity-70 hover:opacity-100 focus:outline-none focus:ring-1 focus:ring-ring"
        >
          <X className="size-4" />
          <span className="sr-only">Dismiss</span>
        </button>
      )}
    </div>
  );
}

function AlertTitle({ className, ...props }: React.ComponentProps<'h5'>) {
  return (
    <h5
      data-slot="alert-title"
      className={cn('mb-1 font-medium leading-none tracking-tight', className)}
      {...props}
    />
  );
}

function AlertDescription({ className, ...props }: React.ComponentProps<'div'>) {
  return (
    <div
      data-slot="alert-description"
      className={cn('text-sm [&_p]:leading-relaxed', className)}
      {...props}
    />
  );
}

export { Alert, AlertTitle, AlertDescription };
