import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import { AudioService } from "bindings/vcalendar-v2/service";
import { useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
export const Route = createFileRoute("/audio-search")({
  component: RouteComponent,
});

function RouteComponent() {
  const [recording, setRecording] = useState<boolean>(false);
  const [transcript, setTranscript] = useState<string[]>([]);
  const [final, setFinal] = useState<string | null>(null);
  const toggleRecording = () => {
    if (!recording) {
      AudioService.StartRecord();
    } else {
      AudioService.StopRecord();
    }
    setRecording((prev: boolean) => !prev);
  };

  useEffect(() => {
    Events.On("vcalendar-v2:send-transcription", (event) => {
      if (event.data.IsFinal) {
        setFinal("final: " + event.data.Message);

        setTranscript([]);
      } else {
        setFinal(null);

        setTranscript((prev) => {
          // Ensure event.data is a string
          if (
            typeof event.data.Message !== "string" ||
            !event.data.Message.trim()
          ) {
            return prev;
          }

          const incomingWords = event.data.Message.trim().split(/\s+/);
          const lastWord = incomingWords[incomingWords.length - 1];

          // Only proceed if we actually got a word
          if (!lastWord) {
            return prev;
          }

          if (incomingWords.length <= prev.length) {
            const prevLastWord = prev[prev.length - 1];

            if (prevLastWord !== lastWord) {
              const newTranscript = [...prev];
              newTranscript[newTranscript.length - 1] = lastWord;
              return newTranscript;
            }

            return prev;
          }

          return [...prev, lastWord];
        });
      }
    });
  }, []);

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
