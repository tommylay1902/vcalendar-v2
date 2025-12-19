import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import { AudioService } from "bindings/vcalendar-v2/service";
import { useState } from "react";

export const Route = createFileRoute("/audio-search")({
  component: RouteComponent,
});

function RouteComponent() {
  const [recording, setRecording] = useState<boolean>(false);

  const toggleRecording = () => {
    if (!recording) {
      AudioService.StartRecord();
    } else {
      AudioService.StopRecord();
    }
    setRecording((prev: boolean) => !prev);
  };

  return (
    <div>
      <Button onClick={toggleRecording}>
        {!recording ? "Start Recording" : "Stop Recording"}
      </Button>
    </div>
  );
}
