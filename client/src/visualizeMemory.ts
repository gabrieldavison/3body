import { colors } from "./design";
import { MemoryObject, MemoryState, VisualizationConfig } from "./types";

// Visual constants
const CONSTANTS = {
  CELL_GAP: 4,
  NODE_COLOR: colors.blue,
  HEAD_COLOR: colors.organge,
  EMPTY_COLOR: colors.yellow,
  LINK_COLOR: colors.black,
  LINK_OUTLINE_COLOR: "#FFFFFF",
  CURRENT_NODE_COLOR: colors.red,
  BACKGROUND_COLOR: colors.background,
  LABEL_COLOR: colors.black,
  LINK_WIDTH: 2,
  LINK_OUTLINE_WIDTH: 4,
  HOVER_BORDER_COLOR: colors.black,
  HOVER_BORDER_WIDTH: 2,
  EMPTY_CELL_BORDER: colors.black,
  LABEL_FONT: "bold {size}px 'Courier New', monospace",
};

// Keep track of hover state outside render loop
let hoveredCell: { row: number; col: number } | null = null;
let currentCanvas: HTMLCanvasElement | null = null;

export function visualizeMemory(
  memory: MemoryState,
  containerId: string,
  config: VisualizationConfig
) {
  const container = document.getElementById(containerId);
  if (!container) return;

  // Get existing canvas or create new one
  let canvas = container.querySelector("canvas");
  if (!canvas) {
    canvas = document.createElement("canvas");
    container.appendChild(canvas);

    // Set up event listeners when creating a new canvas
    const getCellFromEvent = (event: MouseEvent) => {
      const rect = canvas!.getBoundingClientRect();
      const dpr = window.devicePixelRatio || 1;
      const labelSize = config.labelSize ?? 25;

      const scaleX = canvas!.width / rect.width;
      const scaleY = canvas!.height / rect.height;

      const x = (event.clientX - rect.left) * scaleX;
      const y = (event.clientY - rect.top) * scaleY;

      const adjustedX = x / dpr - labelSize;
      const adjustedY = y / dpr - labelSize;

      if (adjustedX < 0 || adjustedY < 0) return null;

      const rows = 20; // Fixed grid size
      const cols = 20;

      const availableWidth = config.width - (cols - 1) * CONSTANTS.CELL_GAP;
      const availableHeight = config.height - (rows - 1) * CONSTANTS.CELL_GAP;
      const cellWidth = availableWidth / cols;
      const cellHeight = availableHeight / rows;
      const cellSize = Math.min(cellWidth, cellHeight);

      const col = Math.floor(adjustedX / (cellSize + CONSTANTS.CELL_GAP));
      const row = Math.floor(adjustedY / (cellSize + CONSTANTS.CELL_GAP));

      const maxX = adjustedX - col * (cellSize + CONSTANTS.CELL_GAP);
      const maxY = adjustedY - row * (cellSize + CONSTANTS.CELL_GAP);

      if (maxX > cellSize || maxY > cellSize) return null;
      if (row >= rows || col >= cols) return null;

      return { row, col };
    };

    canvas.addEventListener("mousemove", (event: MouseEvent) => {
      hoveredCell = getCellFromEvent(event);
    });

    canvas.addEventListener("mouseleave", () => {
      hoveredCell = null;
    });

    canvas.addEventListener("click", (event: MouseEvent) => {
      const cell = getCellFromEvent(event);
      if (cell) {
        navigator.clipboard
          .writeText(`${cell.row} ${cell.col} `)
          .catch(console.error);
      }
    });
  }

  currentCanvas = canvas;

  // Setup canvas with DPR
  const dpr = window.devicePixelRatio || 1;
  const width = config.width + (config.labelSize ?? 25);
  const height = config.height + (config.labelSize ?? 25);

  canvas.style.width = width + "px";
  canvas.style.height = height + "px";
  canvas.width = width * dpr;
  canvas.height = height * dpr;

  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  ctx.scale(dpr, dpr);

  const rows = 20; // Fixed grid size
  const cols = 20;
  const labelSize = config.labelSize ?? 25;

  // Calculate cell size based on available space
  const availableWidth = config.width - (cols - 1) * CONSTANTS.CELL_GAP;
  const availableHeight = config.height - (rows - 1) * CONSTANTS.CELL_GAP;
  const cellWidth = availableWidth / cols;
  const cellHeight = availableHeight / rows;
  const cellSize = Math.min(cellWidth, cellHeight);

  // Clear the canvas with background color
  ctx.fillStyle = CONSTANTS.BACKGROUND_COLOR;
  ctx.fillRect(0, 0, width, height);

  // Create a map of occupied cells for quick lookup
  const occupiedCells = new Set(
    memory.objects.map((obj: MemoryObject) => `${obj.x},${obj.y}`)
  );

  // Draw empty cells first
  for (let row = 0; row < rows; row++) {
    for (let col = 0; col < cols; col++) {
      const cellX = col * (cellSize + CONSTANTS.CELL_GAP) + labelSize;
      const cellY = row * (cellSize + CONSTANTS.CELL_GAP) + labelSize;

      // Draw border for all cells
      ctx.strokeStyle = CONSTANTS.EMPTY_CELL_BORDER;
      ctx.lineWidth = 1;
      ctx.strokeRect(cellX, cellY, cellSize, cellSize);

      // Draw hover border for empty cells
      if (
        hoveredCell &&
        hoveredCell.row === row &&
        hoveredCell.col === col &&
        !occupiedCells.has(`${col},${row}`)
      ) {
        ctx.strokeStyle = CONSTANTS.HOVER_BORDER_COLOR;
        ctx.lineWidth = CONSTANTS.HOVER_BORDER_WIDTH;
        ctx.strokeRect(cellX, cellY, cellSize, cellSize);
      }
    }
  }

  // Draw coordinate labels
  ctx.fillStyle = CONSTANTS.LABEL_COLOR;
  ctx.font = CONSTANTS.LABEL_FONT.replace(
    "{size}",
    (labelSize * 0.5).toString()
  );

  // Draw coordinates
  for (let col = 0; col < cols; col++) {
    const x = col * (cellSize + CONSTANTS.CELL_GAP) + cellSize / 2 + labelSize;
    const y = labelSize / 2;
    ctx.fillText(col.toString(), x - 4, y + 4);
  }

  for (let row = 0; row < rows; row++) {
    const x = labelSize / 2;
    const y = row * (cellSize + CONSTANTS.CELL_GAP) + cellSize / 2 + labelSize;
    ctx.fillText(row.toString(), x - 4, y + 4);
  }

  // Helper function for drawing links
  const drawLink = (
    from: { x: number; y: number },
    to: { x: number; y: number }
  ) => {
    if (to.x === null || to.y === null) return;

    const fromCenterX =
      from.x * (cellSize + CONSTANTS.CELL_GAP) + cellSize / 2 + labelSize;
    const fromCenterY =
      from.y * (cellSize + CONSTANTS.CELL_GAP) + cellSize / 2 + labelSize;
    const toCenterX =
      to.x * (cellSize + CONSTANTS.CELL_GAP) + cellSize / 2 + labelSize;
    const toCenterY =
      to.y * (cellSize + CONSTANTS.CELL_GAP) + cellSize / 2 + labelSize;

    // Draw the white outline first
    ctx!.beginPath();
    ctx!.strokeStyle = CONSTANTS.LINK_OUTLINE_COLOR;
    ctx!.lineWidth = CONSTANTS.LINK_OUTLINE_WIDTH;
    ctx!.moveTo(fromCenterX, fromCenterY);
    ctx!.lineTo(toCenterX, toCenterY);
    ctx!.stroke();

    // Draw the actual link line
    ctx!.beginPath();
    ctx!.strokeStyle = CONSTANTS.LINK_COLOR;
    ctx!.lineWidth = CONSTANTS.LINK_WIDTH;
    ctx!.moveTo(fromCenterX, fromCenterY);
    ctx!.lineTo(toCenterX, toCenterY);
    ctx!.stroke();
  };

  for (const obj of memory.objects) {
    const cellX = obj.x * (cellSize + CONSTANTS.CELL_GAP) + labelSize;
    const cellY = obj.y * (cellSize + CONSTANTS.CELL_GAP) + labelSize;

    // Fill cell with appropriate color
    ctx.fillStyle = CONSTANTS.EMPTY_COLOR;
    if (obj.type === "nod") {
      ctx.fillStyle = obj.isCurrent
        ? CONSTANTS.CURRENT_NODE_COLOR
        : CONSTANTS.NODE_COLOR;
    } else if (obj.type === "hed") {
      ctx.fillStyle = CONSTANTS.HEAD_COLOR;
    }

    // Draw cell background
    ctx.fillRect(cellX, cellY, cellSize, cellSize);

    // Draw hover border if this is the hovered cell
    if (hoveredCell && hoveredCell.row === obj.y && hoveredCell.col === obj.x) {
      ctx.strokeStyle = CONSTANTS.HOVER_BORDER_COLOR;
      ctx.lineWidth = CONSTANTS.HOVER_BORDER_WIDTH;
      ctx.strokeRect(cellX, cellY, cellSize, cellSize);
    }
  }

  // Draw all links after cells are drawn
  for (const obj of memory.objects) {
    // Check if this object has a connection
    if (obj.connectsToX != undefined && obj.connectsToY != undefined) {
      drawLink(
        { x: obj.x, y: obj.y },
        { x: obj.connectsToX, y: obj.connectsToY }
      );
    }
  }
}
