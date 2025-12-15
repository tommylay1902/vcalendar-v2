import { useState, useEffect } from "react";
import { AudioService } from "../bindings/changeme/service";

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
    <div className="container">
      {!recording ? (
        <button onClick={startRecording} className="w">
          Start Recording
        </button>
      ) : (
        <button onClick={stopRecording}>Stop Recording</button>
      )}
    </div>
  );
}

export default App;
