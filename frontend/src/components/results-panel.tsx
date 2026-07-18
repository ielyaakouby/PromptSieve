"use client";

import { useState } from "react";
import type { Finding, ScanResponse } from "@/lib/types";
import {
  categoryLabel,
  findingLine,
  groupFindings,
  severityTone,
  statusMeta,
  typeLabel,
  uniqueCategories,
} from "@/lib/status";
import { AlertIcon, CheckIcon, CopyIcon, ResetIcon, ShieldCheckIcon } from "./icons";
import { btnOutline, btnPrimary, btnSecondary, cardClass } from "./ui";

type ResultsPanelProps = {
  result: ScanResponse | null;
  originalText: string;
  error: string;
  onReset: () => void;
};

function toneClasses(tone: "safe" | "review" | "sensitive" | "critical" | "info") {
  switch (tone) {
    case "safe":
      return "border-success/30 bg-success/10 text-success";
    case "review":
      return "border-warning/30 bg-warning/10 text-warning";
    case "sensitive":
      return "border-warning/40 bg-warning/15 text-warning";
    case "critical":
      return "border-danger/30 bg-danger/10 text-danger";
    case "info":
      return "border-info/30 bg-info/10 text-info";
  }
}

export function ResultsPanel({ result, originalText, error, onReset }: ResultsPanelProps) {
  const [copied, setCopied] = useState<"safe" | "original" | null>(null);
  const [showErrorDetails, setShowErrorDetails] = useState(false);

  const copy = async (value: string, kind: "safe" | "original") => {
    await navigator.clipboard.writeText(value);
    setCopied(kind);
    window.setTimeout(() => setCopied(null), 1600);
  };

  if (error) {
    return (
      <section className={`${cardClass} p-5`} aria-labelledby="scan-error-title">
        <div className="flex items-start gap-3">
          <AlertIcon className="mt-0.5 h-5 w-5 text-danger" />
          <div className="min-w-0 flex-1">
            <h2 id="scan-error-title" className="text-base font-semibold text-foreground">
              Scan failed
            </h2>
            <p className="mt-1 text-sm text-muted-foreground">
              PromptSieve could not complete the scan. Please try again.
            </p>
            <button
              type="button"
              className="mt-3 text-xs font-medium text-muted-foreground underline-offset-2 hover:underline"
              onClick={() => setShowErrorDetails((open) => !open)}
              aria-expanded={showErrorDetails}
            >
              {showErrorDetails ? "Hide details" : "Show details"}
            </button>
            {showErrorDetails && (
              <pre className="mt-2 overflow-x-auto rounded-lg bg-muted px-3 py-2 font-mono text-xs text-muted-foreground">
                {error}
              </pre>
            )}
          </div>
        </div>
      </section>
    );
  }

  if (!result) {
    return (
      <section className={`${cardClass} px-5 py-8`} aria-labelledby="empty-scan-title">
        <div className="mx-auto max-w-md text-center">
          <ShieldCheckIcon className="mx-auto h-5 w-5 text-muted-foreground" />
          <h2 id="empty-scan-title" className="mt-3 text-sm font-semibold text-foreground">
            No scan yet
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Paste a prompt above to check it for sensitive information.
          </p>
        </div>
      </section>
    );
  }

  const findings = result.findings || [];
  const meta = statusMeta(result.status);

  if (findings.length === 0) {
    return (
      <section className={`${cardClass} p-5`} aria-labelledby="safe-scan-title">
        <div className="flex items-start gap-3">
          <ShieldCheckIcon className="mt-0.5 h-5 w-5 text-success" />
          <div>
            <p className={`inline-flex rounded-md border px-2 py-0.5 text-xs font-medium ${toneClasses("safe")}`}>
              {meta.label}
            </p>
            <h2 id="safe-scan-title" className="mt-2 text-base font-semibold text-foreground">
              No sensitive information detected
            </h2>
            <p className="mt-1 text-sm text-muted-foreground">
              PromptSieve did not identify any sensitive patterns using the currently enabled
              detection rules.
            </p>
          </div>
        </div>
      </section>
    );
  }

  const categories = uniqueCategories(findings);

  return (
    <section className="space-y-4" aria-labelledby="scan-results-title">
      <div className="flex items-end justify-between gap-3">
        <h2 id="scan-results-title" className="text-lg font-semibold tracking-tight text-foreground">
          Scan results
        </h2>
        <button type="button" onClick={onReset} className={btnSecondary}>
          <ResetIcon className="h-4 w-4" />
          Reset
        </button>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <Metric label="Risk score" value={`${result.score} / 100`} hint={meta.label} tone={meta.tone} />
        <Metric
          label="Sensitive items"
          value={String(findings.length)}
          hint={findings.length === 1 ? "finding" : "findings"}
        />
        <Metric
          label="Categories"
          value={String(categories.length)}
          hint={categories.map(categoryLabel).join(", ")}
        />
      </div>

      <div className={`${cardClass} overflow-hidden`}>
        <div className="border-b border-border px-5 py-3">
          <h3 className="text-sm font-medium text-foreground">Detection results</h3>
        </div>
        <div className="divide-y divide-border">
          {groupFindings(findings).map(([category, items]) => (
            <div key={category} className="px-5 py-4">
              <p className="mb-3 text-xs font-medium uppercase tracking-wide text-muted-foreground">
                {categoryLabel(category)}
              </p>
              <ul className="space-y-3">
                {items.map((finding, index) => (
                  <FindingRow
                    key={`${finding.id}-${index}`}
                    finding={finding}
                    line={findingLine(originalText, finding)}
                  />
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>

      <div className={`${cardClass} p-5`}>
        <h3 className="text-sm font-medium text-foreground">Sanitized prompt</h3>
        <p className="mt-1 text-xs text-muted-foreground">
          Secrets are replaced with redaction tokens. Copy the safe version by default.
        </p>
        <textarea
          readOnly
          value={result.sanitizedText}
          aria-label="Sanitized prompt"
          className="mt-3 min-h-32 w-full resize-y rounded-lg border border-border bg-background px-3 py-3 font-mono text-sm leading-6 text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
        <div className="mt-4 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button
            type="button"
            onClick={() => copy(originalText, "original")}
            className={btnOutline}
            title="Copies the original unredacted prompt"
          >
            {copied === "original" ? <CheckIcon className="h-4 w-4" /> : <CopyIcon className="h-4 w-4" />}
            {copied === "original" ? "Copied original" : "Copy Original"}
          </button>
          <button
            type="button"
            onClick={() => copy(result.sanitizedText, "safe")}
            className={btnPrimary}
          >
            {copied === "safe" ? <CheckIcon className="h-4 w-4" /> : <CopyIcon className="h-4 w-4" />}
            {copied === "safe" ? "Copied safe prompt" : "Copy Safe Prompt"}
          </button>
        </div>
      </div>
    </section>
  );
}

function Metric({
  label,
  value,
  hint,
  tone,
}: {
  label: string;
  value: string;
  hint?: string;
  tone?: "safe" | "review" | "sensitive" | "critical" | "info";
}) {
  return (
    <div className={`${cardClass} p-4`}>
      <p className="text-xs font-medium text-muted-foreground">{label}</p>
      <p className="mt-1 text-xl font-semibold tracking-tight text-foreground">{value}</p>
      {hint && (
        <p className={`mt-2 inline-flex rounded-md border px-2 py-0.5 text-xs font-medium ${tone ? toneClasses(tone) : "border-border text-muted-foreground"}`}>
          {hint}
        </p>
      )}
    </div>
  );
}

function FindingRow({ finding, line }: { finding: Finding; line: number }) {
  const tone = severityTone(finding.severity);
  return (
    <li className="rounded-lg border border-border bg-background px-3 py-3">
      <div className="flex flex-wrap items-start justify-between gap-2">
        <div>
          <p className="text-sm font-medium text-foreground">{typeLabel(finding.type)}</p>
          <p className="mt-0.5 text-xs text-muted-foreground">{finding.description}</p>
        </div>
        <span className={`rounded-md border px-2 py-0.5 text-[11px] font-semibold tracking-wide ${toneClasses(tone)}`}>
          {finding.severity}
        </span>
      </div>
      <p className="mt-2 truncate font-mono text-xs text-muted-foreground">
        {finding.maskedValue}
      </p>
      <p className="mt-2 text-xs text-muted-foreground">
        Line {line}
        {finding.recommendation ? ` · ${finding.recommendation}` : ""}
      </p>
    </li>
  );
}
