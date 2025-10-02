import { useEffect, useState, type MouseEvent } from "react";
import Login from "./components/Login";
import ArtworkUploader from "./components/ArtworkUploader";
import Gallery from "./components/Gallery";

type HelloResponse = {
  message: string;
};

function App() {
  const [msg, setMsg] = useState<string>("Loading...");
  const [upload, setUpload] = useState<boolean>(false);

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

  const handleClick = (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setUpload(!upload)
  }

  return (
    <>
      <div style={{ fontFamily: "system-ui, sans-serif", padding: "1rem" }}>
        {/* imagine bubbly, baloon, but stylized sharp modern all cap isometric letters floating here, three.js? */}
        <h1>GO ART!</h1>
        <h2>{msg}</h2>
        <p>go-art starter with Go + Vite + React + TypeScript 🎨</p>


        {upload ? <div><button onClick={handleClick}>X Close Uploader</button> <ArtworkUploader /> </div> : <button onClick={handleClick}>^ Open Image Uploader</button>}
        <Gallery />

        <Login />
      </div>
    </>
  );
}

export default App;
