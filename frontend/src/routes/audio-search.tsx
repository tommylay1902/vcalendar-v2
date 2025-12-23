import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import useRecord from "@/hooks/useRecord";
import EventsSection from "@/components/events";

export const Route = createFileRoute("/audio-search")({
  component: RouteComponent,
});

function RouteComponent() {
  const { recording, transcript, final, toggleRecording, events } = useRecord();
  return (
    <div className="flex flex-col">
      {recording && final === null ? (
        <div> {transcript.join(" ")}</div>
      ) : (
        final != null && <EventsSection events={events} />
      )}
      <div className="flex justify-center mt-20">
        <Button className="cursor-pointer" onClick={toggleRecording}>
          {!recording ? "Start Recording" : "Stop Recording"}
        </Button>
      </div>
    </div>
  );
}
