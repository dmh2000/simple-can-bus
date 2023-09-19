import { useState, useEffect } from "react";
import axios from "axios";

import "./App.css";

interface CanDevice {
  DioIn: number;
  DioOut: number;
  DacIn: number;
  AdcOut: number;
}

function App() {
  const [canDevice, setCanDevice] = useState<CanDevice>({
    DioIn: 0,
    DioOut: 0,
    DacIn: 0,
    AdcOut: 0,
  });

  useEffect(() => {
    console.log("useEffect");
    axios({
      // Endpoint to send files
      url: "http://localhost:6001/can/3",
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((response) => {
        console.log("response", response.data);
        return response.data;
      })
      .then((data) => {
        console.log("data", data);
        setCanDevice(data);
      })
      .catch((error) => {
        console.log("error", error);
      });
  }, []);

  return (
    <div className="App">
      <div>Dio In {canDevice.DioIn} </div>
      <div>Dio Out {canDevice.DioOut}</div>
      <div>Dac In {canDevice.DacIn}</div>
      <div>Adc Out {canDevice.AdcOut}</div>
    </div>
  );
}

export default App;
