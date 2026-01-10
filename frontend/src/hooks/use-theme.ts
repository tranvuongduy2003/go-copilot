import { useUIStore } from '@/stores/ui-store';

export function useTheme() {
  const { theme, setTheme } = useUIStore();
  return { theme, setTheme };
}
