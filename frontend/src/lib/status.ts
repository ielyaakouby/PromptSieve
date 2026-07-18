import type { Finding, ScanStatus } from "./types";

export const CATEGORY_LABELS: Record<string, string> = {
  personal_information: "Personal Information",
  credential: "Credentials",
  api_key: "API Keys",
  source_control_token: "Source Control Tokens",
  cloud_credential: "Cloud Credentials",
  auth_token: "Authentication Tokens",
  financial_information: "Financial Information",
  cryptographic_material: "Private Keys",
};

export const TYPE_LABELS: Record<string, string> = {
  email: "Email",
  phone: "Phone Number",
  date_of_birth: "Date of Birth",
  address: "Address",
  person_name: "Person Name",
  password: "Password",
  database_password: "Database Password",
  generic_secret: "Secret",
  openai_api_key: "OpenAI API Key",
  github_token: "GitHub Token",
  aws_access_key: "AWS Access Key",
  aws_secret_access_key: "AWS Secret Access Key",
  jwt: "JWT",
  bearer_token: "Bearer Token",
  private_key: "Private Key",
  iban: "IBAN",
  bic: "BIC / SWIFT",
  payment_card: "Payment Card Number",
  cvv: "CVV / CVC",
  card_expiry: "Card Expiration",
};

export function categoryLabel(category: string): string {
  return CATEGORY_LABELS[category] || titleCase(category.replace(/_/g, " "));
}

export function typeLabel(type: string): string {
  return TYPE_LABELS[type] || titleCase(type.replace(/_/g, " "));
}

export function titleCase(value: string): string {
  return value.replace(/\b\w/g, (char) => char.toUpperCase());
}

export function lineNumber(text: string, index: number): number {
  if (index <= 0) return 1;
  return text.slice(0, index).split("\n").length;
}

export function findingLine(text: string, finding: Finding): number {
  if (typeof finding.line === "number" && finding.line > 0) {
    return finding.line;
  }
  return lineNumber(text, finding.start);
}

export function uniqueCategories(findings: Finding[]): string[] {
  return [...new Set(findings.map((finding) => finding.category))];
}

export function groupFindings(findings: Finding[]): Array<[string, Finding[]]> {
  const groups = new Map<string, Finding[]>();
  for (const finding of findings) {
    const key = finding.category || "other";
    const list = groups.get(key);
    if (list) {
      list.push(finding);
    } else {
      groups.set(key, [finding]);
    }
  }
  return [...groups.entries()];
}

interface StatusMeta {
  id: ScanStatus;
  label: string;
  headline: string;
  tone: "safe" | "review" | "sensitive" | "critical";
}

const STATUS_META: Record<ScanStatus, StatusMeta> = {
  SAFE: {
    id: "SAFE",
    label: "Safe",
    headline: "No sensitive information detected",
    tone: "safe",
  },
  REVIEW: {
    id: "REVIEW",
    label: "Review recommended",
    headline: "Review recommended",
    tone: "review",
  },
  SENSITIVE: {
    id: "SENSITIVE",
    label: "Sensitive data detected",
    headline: "Sensitive data detected",
    tone: "sensitive",
  },
  CRITICAL: {
    id: "CRITICAL",
    label: "High-risk data detected",
    headline: "High-risk data detected",
    tone: "critical",
  },
};

export function statusMeta(status: string): StatusMeta {
  if (status in STATUS_META) {
    return STATUS_META[status as ScanStatus];
  }
  return {
    id: "REVIEW",
    label: status,
    headline: status,
    tone: "review",
  };
}

export function severityTone(severity: string): "safe" | "review" | "sensitive" | "critical" | "info" {
  const value = severity.toUpperCase();
  if (value === "CRITICAL") return "critical";
  if (value === "HIGH") return "sensitive";
  if (value === "MEDIUM") return "review";
  if (value === "LOW") return "info";
  return "info";
}
