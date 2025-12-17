import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import { AudioService } from "bindings/changeme/service";
import { useState } from "react";

export const Route = createFileRoute("/audio-search")({
  component: RouteComponent,
});

function RouteComponent() {
  const [recording, setRecording] = useState<boolean>(false);

  const toggleRecording = () => {
    setRecording((prev: boolean) => !prev);
    if (recording) {
      AudioService.StartRecord();
    } else {
      AudioService.StopRecord();
    }
  };

  return (
    <div>
      <Button onClick={toggleRecording}>
        {!recording ? "Start Recording" : "Stop Recording"}
      </Button>
    </div>
  );
}
