import { useState, useEffect } from "react";
import { AudioService } from "../bindings/changeme/service";
import { Button } from "./components/ui/button";
function App() {
  const [recording, setRecording] = useState(false);
  const startRecording = () => {
    AudioService.StartRecord();
    setRecording(true);
  };

  const stopRecording = () => {
    AudioService.StopRecord();
    setRecording(false);
  };
  return (
    <div className="flex min-h-svh flex-col items-center justify-center">
      {!recording ? (
        <Button onClick={startRecording}>Start Recording</Button>
      ) : (
        <Button onClick={stopRecording}>Stop Recording</Button>
      )}
    </div>
  );
}

export default App;
