import { Footer } from "@/components/footer";
import { Header } from "@/components/header";
import { Scanner } from "@/components/scanner";

export default function Home() {
  return (
    <div className="flex min-h-full flex-col">
      <Header />
      <Scanner />
      <Footer />
    </div>
  );
}
