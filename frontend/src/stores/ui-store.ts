import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';

type Theme = 'light' | 'dark' | 'system';

interface UIState {
  sidebarCollapsed: boolean;
  sidebarMobileOpen: boolean;
  theme: Theme;
  commandPaletteOpen: boolean;
  activeModal: string | null;
  globalLoading: boolean;
}

interface UIActions {
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  toggleMobileSidebar: () => void;
  setMobileSidebarOpen: (open: boolean) => void;
  setTheme: (theme: Theme) => void;
  toggleCommandPalette: () => void;
  setCommandPaletteOpen: (open: boolean) => void;
  openModal: (modalId: string) => void;
  closeModal: () => void;
  setGlobalLoading: (loading: boolean) => void;
  reset: () => void;
}

type UIStore = UIState & UIActions;

const initialState: UIState = {
  sidebarCollapsed: false,
  sidebarMobileOpen: false,
  theme: 'system',
  commandPaletteOpen: false,
  activeModal: null,
  globalLoading: false,
};

export const useUIStore = create<UIStore>()(
  persist(
    immer((set) => ({
      ...initialState,

      toggleSidebar: () => {
        set((state) => {
          state.sidebarCollapsed = !state.sidebarCollapsed;
        });
      },

      setSidebarCollapsed: (collapsed: boolean) => {
        set((state) => {
          state.sidebarCollapsed = collapsed;
        });
      },

      toggleMobileSidebar: () => {
        set((state) => {
          state.sidebarMobileOpen = !state.sidebarMobileOpen;
        });
      },

      setMobileSidebarOpen: (open: boolean) => {
        set((state) => {
          state.sidebarMobileOpen = open;
        });
      },

      setTheme: (theme: Theme) => {
        set((state) => {
          state.theme = theme;
        });

        const root = window.document.documentElement;
        root.classList.remove('light', 'dark');

        if (theme === 'system') {
          const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches
            ? 'dark'
            : 'light';
          root.classList.add(systemTheme);
        } else {
          root.classList.add(theme);
        }
      },

      toggleCommandPalette: () => {
        set((state) => {
          state.commandPaletteOpen = !state.commandPaletteOpen;
        });
      },

      setCommandPaletteOpen: (open: boolean) => {
        set((state) => {
          state.commandPaletteOpen = open;
        });
      },

      openModal: (modalId: string) => {
        set((state) => {
          state.activeModal = modalId;
        });
      },

      closeModal: () => {
        set((state) => {
          state.activeModal = null;
        });
      },

      setGlobalLoading: (loading: boolean) => {
        set((state) => {
          state.globalLoading = loading;
        });
      },

      reset: () => {
        set(() => ({
          ...initialState,
        }));
      },
    })),
    {
      name: 'ui_preferences',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        sidebarCollapsed: state.sidebarCollapsed,
        theme: state.theme,
      }),
    }
  )
);

export const selectSidebarCollapsed = (state: UIStore) => state.sidebarCollapsed;
export const selectTheme = (state: UIStore) => state.theme;
export const selectGlobalLoading = (state: UIStore) => state.globalLoading;
export const selectActiveModal = (state: UIStore) => state.activeModal;
