import { WailsEvent } from "node_modules/@wailsio/runtime/types/events";
import { useEffect } from "react";
interface EventsSectionProps {
  events: WailsEvent<"vcalendar-v2:send-events"> | null;
}
const EventsSection = ({ events }: EventsSectionProps) => {
  return (
    <div>
      {events?.data.events.map((e, index) => (
        <div key={index}>{e.summary}</div>
      ))}
    </div>
  );
};
export default EventsSection;
