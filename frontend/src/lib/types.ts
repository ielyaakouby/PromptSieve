export interface Finding {
  id: string;
  detector: string;
  category: string;
  type: string;
  severity: string;
  confidence: number;
  start: number;
  end: number;
  maskedValue: string;
  replacementToken?: string;
  line?: number;
  description: string;
  recommendation: string;
}

export interface ScanResponse {
  status: string;
  score: number;
  findings: Finding[];
  sanitizedText: string;
}

export type ScanStatus = "SAFE" | "REVIEW" | "SENSITIVE" | "CRITICAL";
