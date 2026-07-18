export function Footer() {
  return (
    <footer className="border-t border-border">
      <div className="mx-auto flex max-w-5xl flex-col gap-2 px-4 py-6 text-xs text-muted-foreground sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8">
        <p>PromptSieve · Open-source privacy and security scanner for AI prompts.</p>
        <p>v0.1.0 · Self-hosted · No prompt retention</p>
      </div>
    </footer>
  );
}
