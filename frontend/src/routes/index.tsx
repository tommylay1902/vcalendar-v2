import { createFileRoute, Link } from "@tanstack/react-router";
import { Events } from "@wailsio/runtime";
import { useEffect } from "react";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  return (
    <div className="p-2">
      <div>
        <Link to={"/audio-search"}>Audio Search</Link>
      </div>
      <h3>Welcome Home!</h3>
    </div>
  );
}
