"use client";

import { persistTheme } from "@/lib/theme";
import { MoonIcon, SunIcon } from "./icons";
import { btnGhost } from "./ui";

export function ThemeToggle() {
  const toggle = () => {
    const next = document.documentElement.classList.contains("dark") ? "light" : "dark";
    document.documentElement.classList.toggle("dark", next === "dark");
    persistTheme(next);
  };

  return (
    <button type="button" onClick={toggle} className={btnGhost} aria-label="Toggle color theme">
      <SunIcon className="hidden h-4 w-4 dark:block" />
      <MoonIcon className="block h-4 w-4 dark:hidden" />
    </button>
  );
}
