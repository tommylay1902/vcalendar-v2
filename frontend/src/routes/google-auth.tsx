import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/google-auth")({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex flex-col m-3">
      <h1 className="text-center font-bold mb-1">
        {" "}
        Provide google token for authentication
      </h1>

      <Input id="token" placeholder="Google Token" />
    </div>
  );
}
