// Types for memory objects and visualization configuration

/**
 * Represents a single object in memory with its position and properties
 */
export interface MemoryObject {
  // Position in the grid
  x: number;
  y: number;

  // Type of memory object
  type: "nod" | "hed"; // Node or Head

  // Whether this node is currently active
  isCurrent?: boolean;

  // Connection coordinates (for linked structures)
  connectsToX?: number | null;
  connectsToY?: number | null;
}

/**
 * Represents the complete state of memory
 */
export interface MemoryState {
  // Array of all objects currently in memory
  objects: MemoryObject[];
}

/**
 * Configuration options for the memory visualization
 */
export interface VisualizationConfig {
  // Width of the visualization area in pixels
  width: number;

  // Height of the visualization area in pixels
  height: number;

  // Size of the coordinate labels (optional)
  labelSize?: number;
}
