import { Loading } from '@/components/ui/spinner';

export function PageLoader() {
  return (
    <div className="flex h-96 items-center justify-center">
      <Loading text="Loading..." />
    </div>
  );
}
