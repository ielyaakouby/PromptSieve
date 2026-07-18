export const THEME_STORAGE_KEY = "promptsieve-theme";
export const LEGACY_THEME_STORAGE_KEY = "securepaste-theme";

export type Theme = "light" | "dark";

export function readStoredTheme(): Theme | null {
  if (typeof window === "undefined") return null;
  const current = window.localStorage.getItem(THEME_STORAGE_KEY);
  if (current === "light" || current === "dark") return current;

  const legacy = window.localStorage.getItem(LEGACY_THEME_STORAGE_KEY);
  if (legacy === "light" || legacy === "dark") {
    window.localStorage.setItem(THEME_STORAGE_KEY, legacy);
    window.localStorage.removeItem(LEGACY_THEME_STORAGE_KEY);
    return legacy;
  }
  return null;
}

export function preferredTheme(): Theme {
  const stored = readStoredTheme();
  if (stored) return stored;
  if (typeof window !== "undefined" && window.matchMedia("(prefers-color-scheme: dark)").matches) {
    return "dark";
  }
  return "light";
}

export function applyTheme(theme: Theme): void {
  document.documentElement.classList.toggle("dark", theme === "dark");
}

export function persistTheme(theme: Theme): void {
  window.localStorage.setItem(THEME_STORAGE_KEY, theme);
}
