import type { ScanResponse } from "./types";

export function getApiUrl(): string {
  const configured =
    process.env.NEXT_PUBLIC_API_URL ||
    process.env.NEXT_PUBLIC_PROMPTSIEVE_API_URL ||
    process.env.NEXT_PUBLIC_SECUREPASTE_API_URL;

  return (configured || "http://127.0.0.1:8080").replace(/\/$/, "");
}

export async function scanPrompt(text: string): Promise<ScanResponse> {
  const res = await fetch(`${getApiUrl()}/api/v1/scan`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text }),
  });

  if (!res.ok) {
    throw new Error(`Scan request failed (${res.status})`);
  }

  return (await res.json()) as ScanResponse;
}
