type LogoProps = {
  className?: string;
  title?: string;
};

export function Logo({ className = "h-7 w-7", title = "PromptSieve" }: LogoProps) {
  return (
    <svg
      viewBox="0 0 32 32"
      fill="none"
      className={className}
      role="img"
      aria-label={title}
    >
      <title>{title}</title>
      <path
        d="M16 3.4 26.2 7.2v7.6c0 6.1-4.1 10.4-10.2 11.8C9.9 25.2 5.8 20.9 5.8 14.8V7.2Z"
        stroke="currentColor"
        strokeWidth="1.7"
        strokeLinejoin="round"
      />
      <path
        d="M10.4 12.4h11.2M11.5 16.3h9M12.8 20.1h6.4"
        stroke="currentColor"
        strokeWidth="1.7"
        strokeLinecap="round"
      />
    </svg>
  );
}
