import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import type { ReactNode } from "react";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const title = "PromptSieve — Sensitive Data Scanner for AI Prompts";
const description =
  "PromptSieve is an open-source tool that detects sensitive information, secrets and personal data in prompts before they are sent to AI systems.";

export const metadata: Metadata = {
  title,
  description,
  applicationName: "PromptSieve",
  keywords: [
    "prompt security",
    "secret scanning",
    "PII",
    "LLM",
    "privacy",
    "open source",
  ],
  authors: [{ name: "PromptSieve" }],
  openGraph: {
    title,
    description,
    siteName: "PromptSieve",
    type: "website",
    locale: "en_US",
  },
  twitter: {
    card: "summary",
    title,
    description,
  },
};

const themeScript = `(function(){try{var k="promptsieve-theme";var legacy="securepaste-theme";var t=localStorage.getItem(k)||localStorage.getItem(legacy);if(t==="dark"||(t!=="light"&&window.matchMedia("(prefers-color-scheme: dark)").matches)){document.documentElement.classList.add("dark");}}catch(e){}})();`;

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <head>
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
      </head>
      <body className="min-h-full bg-background font-sans text-foreground">{children}</body>
    </html>
  );
}
