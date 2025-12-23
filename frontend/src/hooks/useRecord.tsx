import { AudioService } from "bindings/vcalendar-v2/service";
import { useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import { WailsEvent } from "node_modules/@wailsio/runtime/types/events";

function useRecord() {
  const [recording, setRecording] = useState<boolean>(false);
  const [transcript, setTranscript] = useState<string[]>([]);
  const [final, setFinal] = useState<string | null>(null);
  const [events, setEvents] =
    useState<WailsEvent<"vcalendar-v2:send-events"> | null>(null);
  const toggleRecording = () => {
    if (!recording) {
      AudioService.StartRecord();
    } else {
      AudioService.StopRecord();
    }
    setRecording((prev: boolean) => !prev);
  };

  useEffect(() => {
    Events.On("vcalendar-v2:send-events", (event) => {
      setEvents(event);
    });

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
  return {
    recording,
    transcript,
    final,
    toggleRecording,
    events,
  };
}
export default useRecord;
