"use client";
import Link from "next/link";
import { Dribbble, Menu, Plus } from "lucide-react";
import Image from "next/image";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import YourTeam from "@/components/YourTeamDashView";
import Home from "@/components/HomeDashView";
import LineupGeneration from "@/components/LineupGenerationDashView";

import { ModeToggle } from "@/components/ui/toggle-mode";
import { Pacifico } from "next/font/google";
import { Separator } from "@/components/ui/separator";

const pacifico = Pacifico({
  weight: "400",
  subsets: ["latin"],
});

export default function Dashboard() {
  const [selectedPage, setSelectedPage] = useState("home");

  return (
    <div className="grid min-h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[200px_1fr]">
      <div className="hidden border-r bg-muted/40 md:block">
        <div className="flex h-full max-h-screen flex-col gap-2 items-center">
          <div className="flex h-14 items-center border-b px-4 md:h-[120px] lg:h-[120px]">
            <Image src="/logo.png" alt="Logo" width={100} height={80} />
          </div>
          <div className="flex-1">
            <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
              <div
                className={`flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary cursor-pointer ${
                  selectedPage === "home" ? "text-primary" : ""
                }`}
                onClick={() => setSelectedPage("home")}
              >
                <Plus className="h-4 w-4" />
                Home
              </div>
              <div
                className={`flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary cursor-pointer ${
                  selectedPage === "your-team" ? "text-primary" : ""
                }`}
                onClick={() => setSelectedPage("your-team")}
              >
                <Plus className="h-4 w-4" />
                Your Team
              </div>
              <div
                className={`flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all hover:text-primary cursor-pointer ${
                  selectedPage === "streaming-optimization"
                    ? "text-primary"
                    : ""
                }`}
                onClick={() => setSelectedPage("streaming-optimization")}
              >
                <Plus className="h-4 w-4" />
                Lineup Generation
              </div>
              <Link
                href="#"
                className="flex items-center gap-3 rounded-lg px-3 py-2 text-muted-foreground transition-all"
              >
                <Plus className="h-4 w-4" />
                More Coming Soon
              </Link>
            </nav>
          </div>
        </div>
      </div>
      <div className="flex flex-col">
        <header className="flex h-14 items-center gap-4 border-b bg-muted/40 px-4 md:h-[120px] lg:h-[120px] lg:px-6">
          <Sheet>
            <SheetTrigger asChild>
              <Button
                variant="outline"
                size="icon"
                className="shrink-0 md:hidden"
              >
                <Menu className="h-5 w-5" />
                <span className="sr-only">Toggle navigation menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex flex-col">
              <nav className="grid gap-2 text-lg font-medium">
                <Link
                  href="#"
                  className="flex items-center gap-2 text-lg font-semibold"
                >
                  <Dribbble className="h-6 w-6" />
                  <span className="">Court Visionaries</span>
                </Link>
                <Link
                  href="#"
                  className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                >
                  <Plus className="h-5 w-5" />
                  Your Team
                </Link>
                <Link
                  href="#"
                  className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                >
                  <Plus className="h-5 w-5" />
                  Lineup Generation
                </Link>
                <Link
                  href="#"
                  className="mx-[-0.65rem] flex items-center gap-4 rounded-xl px-3 py-2 text-muted-foreground hover:text-foreground"
                >
                  <Plus className="h-5 w-5" />
                  More Coming Soon
                </Link>
              </nav>
            </SheetContent>
          </Sheet>
          <nav className="w-full">
            <ul className="flex items-center justify-between">
              <hr className="w-1/2 border-primary flex-shrink-1 flex-grow-1"></hr>
              <li
                className={`text-5xl font-bold flex-shrink-0 ${pacifico.className}`}
              >
                Court Visionaries
              </li>
              <hr className="w-1/2 border-primary flex-shrink-1 flex-grow-1"></hr>
              <li className="flex flex-col ring-primary">
                <Separator orientation="vertical" />
              </li>
              <li className="flex-shrink-0 mr-10">
                <ModeToggle />
              </li>
            </ul>
          </nav>
        </header>
        <main className="flex flex-1 flex-col gap-4 p-4 lg:gap-6 lg:p-6">
          {selectedPage === "your-team" && <YourTeam />}
          {selectedPage === "home" && <Home />}
          {selectedPage === "streaming-optimization" && <LineupGeneration />}
        </main>
      </div>
    </div>
  );
}
