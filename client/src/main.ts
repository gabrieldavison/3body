// main.ts
import "./style.css";
import { setupEditor } from "./editor";
import { visualizeMemory } from "./visualizeMemory";
import { colors } from "./design";
import { MemoryState } from "./types";

// Set up the event source for real-time updates

let currentMemoryState: MemoryState = { objects: [] };

function setupEventSource() {
  const eventSource = new EventSource("http://localhost:8080/memory-stream");

  eventSource.onmessage = (event) => {
    try {
      currentMemoryState = JSON.parse(event.data);
      console.log("Received memory state:", currentMemoryState);
    } catch (error) {
      console.error(
        "Error parsing memory state:",
        error,
        "Raw data:",
        event.data
      );
    }
  };

  eventSource.onerror = (error) => {
    console.error("EventSource failed:", error);
    // Attempt to reconnect after a delay
    setTimeout(() => {
      console.log("Attempting to reconnect...");
      eventSource.close();
      setupEventSource();
    }, 5000);
  };

  eventSource.onopen = () => {
    console.log("SSE connection established");
  };

  return eventSource;
}

// Set up the initial event source
let eventSource = setupEventSource();

// Handle incoming memory state updates
eventSource.onmessage = (event) => {
  try {
    currentMemoryState = JSON.parse(event.data);
    // console.log(currentMemoryState);
  } catch (error) {
    console.error("Error parsing memory state:", error);
  }
};

const forthOutput = document.getElementById("forth-output") as HTMLDivElement;

// Set up the editor with evaluation handler
setupEditor("forth-input", async (code) => {
  try {
    const response = await fetch("http://localhost:8080/evaluate", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ input: code }),
    });

    const result = await response.json();

    if (result.error) {
      forthOutput.innerHTML += `<div>> ${code}</div><div style="color: ${colors.red}">Error: ${result.error}</div>`;
    } else if (result.output && result.output.length > 0) {
      forthOutput.innerHTML += `<div>> ${code}</div><div>${result.output.join(
        " "
      )}</div>`;
    } else {
      forthOutput.innerHTML += `<div>> ${code}</div><div>OK</div>`;
    }
  } catch (error) {
    forthOutput.innerHTML += `<div>> ${code}</div><div style="color: ${
      colors.red
    }">Error: ${(error as Error).message}</div>`;
  }
  forthOutput.scrollTop = forthOutput.scrollHeight;
});

// Animation loop for continuous rendering
function animate() {
  console.log(currentMemoryState);
  visualizeMemory(currentMemoryState, "mem-container", {
    width: 350,
    height: 350,
    labelSize: 20,
  });
  requestAnimationFrame(animate);
}

// Start the animation loop
requestAnimationFrame(animate);

// Error handling for SSE connection
eventSource.onerror = (error) => {
  console.error("EventSource failed:", error);
};
