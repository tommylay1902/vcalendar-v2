import { createFileRoute, Link } from "@tanstack/react-router";
import { Events } from "@wailsio/runtime";

export const Route = createFileRoute("/")({
  component: HomeComponent,
});

function HomeComponent() {
  Events.On("vcalendar-v2:token-needed", (event) => {
    console.log(event.data);
  });
  return (
    <div className="p-2">
      <div>
        <Link to={"/audio-search"}>Audio Search</Link>
      </div>
      <h3>Welcome Home!</h3>
    </div>
  );
}
