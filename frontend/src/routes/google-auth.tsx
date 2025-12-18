import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Events } from "@wailsio/runtime";
import { AuthCodeToken } from "bindings/vcalendar-v2/model";
export const Route = createFileRoute("/google-auth")({
  component: RouteComponent,
});

function RouteComponent() {
  const [token, setToken] = useState("");
  const [authUrl, setAuthUrl] = useState("");

  useEffect(() => {
    console.log("hello", authUrl);
    Events.On("vcalendar-v2:auth-url", (event) => {
      setAuthUrl(event.data);
      console.log("HELLOOO", authUrl);
    });
  }, [authUrl]);

  const sendToken = () => {
    const AuthCodeToken = { Token: token };
    Events.Emit("vcalendar-v2:auth-code-token", AuthCodeToken);
  };

  return (
    <div className="flex flex-col m-3">
      <h1 className="text-center font-bold mb-1">
        {" "}
        Click the link below and then input the token you recieve after
        authenticating
      </h1>
      <Link to={authUrl}>Click Here</Link>
      <Input
        id="token"
        placeholder="Google Token"
        value={token}
        onChange={(e) => setToken(e.target.value)}
      />
      <Button onClick={sendToken}>Submit</Button>
    </div>
  );
}
