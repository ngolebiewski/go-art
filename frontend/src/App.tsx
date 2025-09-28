import { useEffect, useState } from "react";

type HelloResponse = {
  message: string;
};

function App() {
  const [msg, setMsg] = useState<string>("Loading...");

  useEffect(() => {
    const getHi = async () => {
      try {
        const response = await fetch("/api/hello");
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const data: HelloResponse = await response.json();
        setMsg(data.message);
      } catch (error: any) {
          setMsg("Error: " + error.message);
      }
    };

    // Artificial delay to see the Loading State
    setTimeout(() => getHi(), 1000)
  }, []);

  return (
    <body>
      <div style={{ fontFamily: "system-ui, sans-serif", padding: "1rem" }}>
        {/* imagine bubbly, baloon, but stylized sharp modern all cap isometric letters floating here, three.js? */}
        <h1>GO ART!</h1>

        <h2>{msg}</h2>
        <p>go-art starter with Go + Vite + React + TypeScript 🎨</p>
      </div>
    </body>
  );
}

export default App;
