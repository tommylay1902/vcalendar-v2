import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import { AudioService } from "bindings/vcalendar-v2/service";
import { useEffect, useState } from "react";
import { CalendarEvents } from "bindings/vcalendar-v2/model";
import { Events } from "@wailsio/runtime";
import useRecord from "@/hooks/useRecord";
export const Route = createFileRoute("/audio-search")({
  component: RouteComponent,
});

function RouteComponent() {
  const { recording, transcript, final, toggleRecording } = useRecord();
  return (
    <div className="flex flex-col">
      {recording && final === null ? (
        <div> {transcript.join(" ")}</div>
      ) : (
        <div>{final}</div>
      )}
      <div className="flex justify-center mt-20">
        <Button className="cursor-pointer" onClick={toggleRecording}>
          {!recording ? "Start Recording" : "Stop Recording"}
        </Button>
      </div>
    </div>
  );
}
