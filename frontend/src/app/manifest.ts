import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "PromptSieve",
    short_name: "PromptSieve",
    description:
      "PromptSieve is an open-source tool that detects sensitive information, secrets and personal data in prompts before they are sent to AI systems.",
    start_url: "/",
    display: "standalone",
    background_color: "#09090b",
    theme_color: "#0f766e",
    icons: [
      {
        src: "/icon.svg",
        sizes: "any",
        type: "image/svg+xml",
      },
    ],
  };
}
