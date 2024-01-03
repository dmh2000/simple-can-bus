import { useState, useEffect } from "react";
import axios from "axios";

import "./App.css";

interface CanDevice {
  DioSet: number;
  DioOut: number;
  DacSet: number;
  AdcOut: number;
}

function App() {
  const [canDevice, setCanDevice] = useState<CanDevice>({
    DioSet: 0,
    DioOut: 0,
    DacSet: 0,
    AdcOut: 0,
  });

  const [Dio, setDio] = useState<number>(0);
  const [updateDio, setUpdateDio] = useState<boolean>(false); // [dio, setDio] = useState<number>(0);
  const [Dac, setDac] = useState<number>(0);
  const [updateDac, setUpdateDac] = useState<boolean>(false); // [dac, setDac] = useState<number>(0);
  const [update, setUpdate] = useState<boolean>(false);

  const headers = {
    "Content-Type": "application/json",
  };

  // update at 1 hz
  useEffect(() => {
    const interval = setInterval(() => {
      axios({
        // Endpoint to send files
        url: "http://localhost:6001/can/3",
        method: "GET",
        headers: headers,
      })
        .then((response) => {
          console.log("get", response.data);
          return response.data;
        })
        .then((data) => {
          console.log("data", data);
          setCanDevice(data);
        })
        .catch((error) => {
          console.log("error", error);
        });
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  // update when changed
  useEffect(() => {
    axios
      .put("http://localhost:6001/can/1", { Dio: Dio }, { headers })
      .then((response) => {
        console.log("put", response.data);
        return response.data;
      })
      .catch((error) => {
        console.log("error", error);
      });
    let data: CanDevice = canDevice;
    data.DioSet = Dio;
    setCanDevice(data);
  }, [updateDio]);

  // update when changed
  useEffect(() => {
    axios
      .put("http://localhost:6001/can/2", { Dac: Dac }, { headers })
      .then((response) => {
        console.log("put", response.data);
        return response.data;
      })
      .catch((error) => {
        console.log("error", error);
      });
    let data: CanDevice = canDevice;
    data.DacSet = Dac;
    setCanDevice(data);
  }, [updateDac]);

  const onSubmit = (event: any) => {
    event.preventDefault();
    console.log("onSubmit", Dio, Dac);
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
        <span className="hdg">Dio Set</span> <span>{canDevice.DioSet} </span>
      </div>
      <div className="aleft">
        <span className="hdg">Dio Out</span> <span> {canDevice.DioOut}</span>
      </div>
      <div className="aleft">
        <span className="hdg">Dac Set</span> <span> {canDevice.DacSet}</span>
      </div>
      <div className="aleft">
        <span className="hdg">Adc Out</span> <span> {canDevice.AdcOut}</span>
      </div>

      <div></div>
      <form onSubmit={onSubmit}>
        <div>
          <input
            type="number"
            value={Dio}
            onChange={onDioChange}
            min={0}
            max={65535}
          />
        </div>
        <div>
          <input
            type="number"
            value={Dac}
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
