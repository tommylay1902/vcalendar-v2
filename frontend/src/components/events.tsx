import { WailsEvent } from "node_modules/@wailsio/runtime/types/events";
import { useEffect, useState } from "react";
import { DateTime } from "luxon";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

interface EventsSectionProps {
  events: WailsEvent<"vcalendar-v2:send-events"> | null;
}

const EventsSection = ({ events }: EventsSectionProps) => {
  const [date, setDate] = useState<string>("");
  const convertDate = (dateString: string | undefined) => {
    if (!dateString) {
      setDate(DateTime.now().toFormat("LLL dd yyyy"));
      return;
    }
    const luxonDateTime = DateTime.fromISO(dateString);
    setDate(luxonDateTime.toFormat("LLL dd yyyy"));
  };

  const getHours = (dateString: string) => {
    return DateTime.fromISO(dateString).toFormat("HH:mm");
  };

  useEffect(() => {
    convertDate(events?.data?.Date);
  }, [events?.data?.Date]);

  if (!events?.data?.events || events.data.events.length === 0) {
    return null;
  }

  return (
    <div className="w-full">
      <div className="text-lg font-medium mb-2 text-center">
        Events for {date}
      </div>
      <div className="flex justify-center">
        <Carousel
          opts={{
            align: "center",
          }}
          className="w-full max-w-lg" // Reduced from max-w-xl to max-w-xs
        >
          <CarouselContent className="-ml-1">
            {events.data.events.map((event, index) => (
              <CarouselItem key={index} className="pl-1 basis-2/3 md:basis-1/2">
                <div className="p-0.5">
                  <Card className="border shadow-sm">
                    <CardHeader className="p-2 font-extrabold">
                      <div className="text-sm text-center font-extrabold">
                        Time: {getHours(event.start.dateTime)}
                      </div>
                    </CardHeader>
                    <CardContent className="p-2 flex items-center justify-center min-h-[160px] min-w-[200px]">
                      <div className="text-sm text-center line-clamp-2">
                        {event.summary}
                      </div>
                    </CardContent>
                  </Card>
                </div>
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselPrevious className="h-6 w-6" /> {/* Smaller buttons */}
          <CarouselNext className="h-6 w-6" />
        </Carousel>
      </div>
    </div>
  );
};
export default EventsSection;
