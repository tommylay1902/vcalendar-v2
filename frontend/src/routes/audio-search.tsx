import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import useRecord from "@/hooks/useRecord";
import EventsSection from "@/components/events";
import { Events } from "@wailsio/runtime";
import { useEffect, useRef } from "react";

export const Route = createFileRoute("/audio-search")({
  component: RouteComponent,
});

function RouteComponent() {
  const { recording, transcript, final, toggleRecording, events } = useRecord();

  return (
    <div className="flex flex-col min-h-screen p-4 bg-gray-50">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Audio Search</h1>
        <p className="text-gray-600">
          Real-time audio processing and visualization
        </p>
      </div>

      {/* Transcript Display */}
      <div className="mb-6 bg-white p-4 rounded-lg shadow-sm border">
        <h2 className="text-lg font-semibold text-gray-700 mb-3">
          Live Transcript
        </h2>
        <div className="min-h-24 p-4 bg-gray-50 rounded border">
          {transcript.length > 0 ? (
            <p className="text-gray-800 leading-relaxed">
              {transcript.join(" ")}
            </p>
          ) : (
            <div className="text-gray-400 italic">
              {recording
                ? "Listening... speak now"
                : "Start recording to capture speech"}
            </div>
          )}
          {final && (
            <div className="mt-2 pt-2 border-t border-gray-200">
              <span className="text-sm font-medium text-gray-600">
                Final result:
              </span>
              <p className="text-gray-800 mt-1">{final}</p>
            </div>
          )}
        </div>
      </div>

      <div className="mb-6">
        <EventsSection events={events} />
      </div>

      <div className="mt-auto pt-6 border-t border-gray-200">
        <div className="flex flex-col items-center gap-4">
          <Button
            className={`px-8 py-3 text-lg font-medium transition-all ${recording ? "animate-pulse" : ""}`}
            onClick={toggleRecording}
            variant={recording ? "destructive" : "default"}
            size="lg"
          >
            {recording ? (
              <>
                <span className="mr-2">⏹️</span>
                Stop Recording
              </>
            ) : (
              <>
                <span className="mr-2">🎤</span>
                Start Recording
              </>
            )}
          </Button>

          {recording && (
            <div className="flex items-center gap-2 text-red-600">
              <div className="w-3 h-3 bg-red-600 rounded-full animate-pulse"></div>
              <span className="text-sm font-medium">
                Recording in progress...
              </span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
