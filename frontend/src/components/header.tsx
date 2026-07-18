import Link from "next/link";
import { Logo } from "./logo";
import { ThemeToggle } from "./theme-toggle";

export function Header() {
  return (
    <header className="border-b border-border bg-background/90 backdrop-blur-sm">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center gap-2.5 text-foreground">
            <Logo className="h-6 w-6 text-primary" />
            <span className="hidden text-sm font-semibold tracking-tight sm:inline">
              PromptSieve
            </span>
          </Link>
          <nav aria-label="Primary" className="flex items-center">
            <span className="rounded-md bg-muted px-2 py-1 text-xs font-medium text-foreground">
              Scanner
            </span>
          </nav>
        </div>
        <div className="flex items-center gap-1">
          <ThemeToggle />
        </div>
      </div>
    </header>
  );
}
