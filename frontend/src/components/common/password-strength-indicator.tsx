import { cn } from '@/lib/utils';
import { useMemo } from 'react';

interface PasswordStrengthIndicatorProps {
  password: string;
  className?: string;
}

interface StrengthConfig {
  label: string;
  color: string;
  bgColor: string;
}

function calculatePasswordStrength(password: string): number {
  if (!password) return 0;

  let strength = 0;

  if (password.length >= 8) strength += 1;
  if (password.length >= 12) strength += 1;
  if (/[a-z]/.test(password)) strength += 1;
  if (/[A-Z]/.test(password)) strength += 1;
  if (/[0-9]/.test(password)) strength += 1;
  if (/[^a-zA-Z0-9]/.test(password)) strength += 1;

  return Math.min(Math.floor((strength / 6) * 4), 4);
}

const strengthConfigs: Record<number, StrengthConfig> = {
  0: { label: 'Too weak', color: 'text-destructive', bgColor: 'bg-destructive' },
  1: { label: 'Weak', color: 'text-orange-500', bgColor: 'bg-orange-500' },
  2: { label: 'Fair', color: 'text-yellow-500', bgColor: 'bg-yellow-500' },
  3: { label: 'Good', color: 'text-blue-500', bgColor: 'bg-blue-500' },
  4: { label: 'Strong', color: 'text-green-500', bgColor: 'bg-green-500' },
};

export function PasswordStrengthIndicator({ password, className }: PasswordStrengthIndicatorProps) {
  const strength = useMemo(() => calculatePasswordStrength(password), [password]);
  const config = strengthConfigs[strength];

  if (!password) return null;

  return (
    <div className={cn('space-y-1', className)}>
      <div className="flex gap-1">
        {[1, 2, 3, 4].map((level) => (
          <div
            key={level}
            className={cn(
              'h-1 flex-1 rounded-full transition-colors',
              level <= strength ? config.bgColor : 'bg-muted'
            )}
          />
        ))}
      </div>
      <p className={cn('text-xs', config.color)}>Password strength: {config.label}</p>
    </div>
  );
}
