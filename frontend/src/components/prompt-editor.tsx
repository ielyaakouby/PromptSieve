"use client";

import type { KeyboardEvent } from "react";
import { LockIcon, ScanIcon, SpinnerIcon, TrashIcon } from "./icons";
import { btnGhost, btnPrimary, cardClass } from "./ui";

type PromptEditorProps = {
  text: string;
  loading: boolean;
  onChange: (value: string) => void;
  onScan: () => void;
  onClear: () => void;
};

export function PromptEditor({ text, loading, onChange, onScan, onClear }: PromptEditorProps) {
  const onKeyDown = (event: KeyboardEvent<HTMLTextAreaElement>) => {
    if ((event.metaKey || event.ctrlKey) && event.key === "Enter") {
      event.preventDefault();
      if (text.trim().length > 0 && !loading) {
        onScan();
      }
    }
  };

  return (
    <section className={`${cardClass} p-4 sm:p-5`} aria-labelledby="prompt-editor-title">
      <div className="mb-3 flex items-center justify-between gap-3">
        <h2 id="prompt-editor-title" className="text-sm font-medium text-foreground">
          Paste your prompt
        </h2>
        <p className="font-mono text-xs text-muted-foreground">
          {text.length.toLocaleString()} chars
        </p>
      </div>
      <label htmlFor="prompt-input" className="sr-only">
        Prompt to scan
      </label>
      <textarea
        id="prompt-input"
        value={text}
        onChange={(event) => onChange(event.target.value)}
        onKeyDown={onKeyDown}
        placeholder="Paste the prompt you want to check before sending it to an AI model..."
        spellCheck={false}
        className="min-h-48 w-full resize-y rounded-lg border border-border bg-background px-3 py-3 font-mono text-sm leading-6 text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      />
      <div className="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p className="flex items-start gap-2 text-xs text-muted-foreground">
          <LockIcon className="mt-0.5 h-3.5 w-3.5 shrink-0" />
          <span>
            Analyzed by your PromptSieve server. Prompts are not stored.
          </span>
        </p>
        <div className="flex items-center justify-end gap-2">
          {text.length > 0 && (
            <button type="button" onClick={onClear} className={btnGhost}>
              <TrashIcon className="h-4 w-4" />
              Clear
            </button>
          )}
          <button
            type="button"
            onClick={onScan}
            disabled={loading || text.trim().length === 0}
            className={btnPrimary}
          >
            {loading ? <SpinnerIcon className="h-4 w-4" /> : <ScanIcon className="h-4 w-4" />}
            {loading ? "Scanning..." : "Scan Prompt"}
          </button>
        </div>
      </div>
    </section>
  );
}
