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

  const [dio, setDio] = useState<number>(0);
  const [updateDio, setUpdateDio] = useState<boolean>(false); // [dio, setDio] = useState<number>(0);
  const [dac, setDac] = useState<number>(0);
  const [updateDac, setUpdateDac] = useState<boolean>(false); // [dac, setDac] = useState<number>(0);
  const [update, setUpdate] = useState<boolean>(false);

  const headers = {
    "Content-Type": "application/json",
  };
  useEffect(() => {
    axios({
      // Endpoint to send files
      url: "http://localhost:6001/can/3",
      method: "GET",
      headers: headers,
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
  }, [update]);

  useEffect(() => {
    axios
      .put("http://localhost:6001/can/1", { dio }, { headers })
      .then((response) => {
        console.log("response", response.data);
        return response.data;
      })
      .catch((error) => {
        console.log("error", error);
      });
  }, [updateDio]);

  useEffect(() => {
    axios
      .put("http://localhost:6001/can/2", { dac }, { headers })
      .then((response) => {
        console.log("response", response.data);
        return response.data;
      })
      .catch((error) => {
        console.log("error", error);
      });
  }, [updateDac]);

  const onSubmit = (event: any) => {
    event.preventDefault();
    console.log("onSubmit", dio, dac);
    setUpdate(!update);
    setUpdateDio(!updateDio);
    setUpdateDac(!updateDac);
  };

  const onDioChange = (event: any) => {
    console.log(event.target.value);
    setDio(event.target.value as number);
  };

  const onDacChange = (event: any) => {
    console.log(event.target.value);
    setDac(event.target.value as number);
  };

  return (
    <div className="App">
      <div className="aleft">
        <span className="hdg">Dio In</span> <span>{canDevice.DioIn} </span>
      </div>
      <div className="aleft">
        <span className="hdg">Dio Out</span> <span> {canDevice.DioOut}</span>
      </div>
      <div className="aleft">
        <span className="hdg">Dac In</span> <span> {canDevice.DacIn}</span>
      </div>
      <div className="aleft">
        <span className="hdg">Adc Out</span> <span> {canDevice.AdcOut}</span>
      </div>

      <div></div>
      <form onSubmit={onSubmit}>
        <div>
          <input
            type="number"
            value={dio}
            onChange={onDioChange}
            min={0}
            max={65535}
          />
        </div>
        <div>
          <input
            type="number"
            value={dac}
            onChange={onDacChange}
            min={0}
            max={10000}
          />
        </div>
        <div>
          <button type="submit">Send</button>
        </div>
      </form>
    </div>
  );
}

export default App;
