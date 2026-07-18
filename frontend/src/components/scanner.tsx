"use client";

import { useState } from "react";
import { scanPrompt } from "@/lib/api";
import type { ScanResponse } from "@/lib/types";
import { ChevronIcon } from "./icons";
import { PromptEditor } from "./prompt-editor";
import { ResultsPanel } from "./results-panel";

const WORKFLOW = ["Paste", "Scan", "Detect", "Review", "Sanitize", "Copy"] as const;

export function Scanner() {
  const [text, setText] = useState("");
  const [result, setResult] = useState<ScanResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const scan = async () => {
    setLoading(true);
    setError("");
    try {
      const data = await scanPrompt(text);
      setResult(data);
    } catch (err: unknown) {
      setResult(null);
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    } finally {
      setLoading(false);
    }
  };

  const reset = () => {
    setText("");
    setResult(null);
    setError("");
  };

  const onChange = (value: string) => {
    setText(value);
    if (result || error) {
      setResult(null);
      setError("");
    }
  };

  return (
    <main className="mx-auto flex w-full max-w-5xl flex-1 flex-col gap-8 px-4 py-8 sm:px-6 sm:py-10 lg:px-8">
      <section className="space-y-3">
        <p className="text-xs font-medium uppercase tracking-[0.14em] text-muted-foreground">
          Prompt security
        </p>
        <h1 className="max-w-2xl text-[1.75rem] font-semibold tracking-tight text-foreground sm:text-[2rem]">
          Protect your data before you prompt.
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground sm:text-[15px]">
          PromptSieve scans AI prompts for sensitive information before they are sent to an AI
          model.
        </p>
        <ol className="flex flex-wrap items-center gap-1 pt-1 text-xs text-muted-foreground">
          {WORKFLOW.map((step, index) => (
            <li key={step} className="flex items-center gap-1">
              {index > 0 && <ChevronIcon className="h-3 w-3" />}
              <span>{step}</span>
            </li>
          ))}
        </ol>
      </section>

      <PromptEditor
        text={text}
        loading={loading}
        onChange={onChange}
        onScan={scan}
        onClear={reset}
      />

      <ResultsPanel result={result} originalText={text} error={error} onReset={reset} />
    </main>
  );
}
