import { Header } from "@/components/Header";
import { Hero } from "@/components/Hero";
import { ArchitectureSection, VisionSection } from "@/components/Architecture";
import { FeaturesSection } from "@/components/Features";
import { LifecycleSection } from "@/components/Lifecycle";
import { MemorySection } from "@/components/Memory";
import { ResearchSection } from "@/components/Research";
import { RoadmapSection } from "@/components/Roadmap";
import { Footer } from "@/components/Footer";

export default function HomePage() {
  return (
    <>
      <Header />
      <main>
        <Hero />
        <VisionSection />
        <ArchitectureSection />
        <FeaturesSection />
        <LifecycleSection />
        <MemorySection />
        <ResearchSection />
        <RoadmapSection />
      </main>
      <Footer />
    </>
  );
}
