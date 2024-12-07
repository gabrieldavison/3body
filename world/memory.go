// world/memory.go
package world

import (
	"fmt"
	"sync"
)

// Memory2D represents a 2D grid of nodes and heads
type Memory2D struct {
	mem  [][]*Nod
	heds []*Hed
	mu   sync.RWMutex // Protects concurrent access
}

// NewMemory2D creates a new 2D memory grid
func NewMemory2D(rows, cols int) *Memory2D {
	mem := make([][]*Nod, rows)
	for i := range mem {
		mem[i] = make([]*Nod, cols)
	}

	return &Memory2D{
		mem:  mem,
		heds: make([]*Hed, 0),
	}
}

// AddNod adds a node to the memory grid
func (m *Memory2D) AddNod(x, y int, nod *Nod) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.checkBounds(x, y); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	m.mem[y][x] = nod
	return nil
}

// AddHed adds a head to memory
func (m *Memory2D) AddHed(x, y int, hed *Hed) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.checkBounds(x, y); err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	// Remove any existing head with the same ID
	for i := len(m.heds) - 1; i >= 0; i-- {
		if m.heds[i].ID() == hed.ID() {
			m.heds = append(m.heds[:i], m.heds[i+1:]...)
		}
	}

	m.heds = append(m.heds, hed)
	return nil
}

// GetNod retrieves a node from the grid
func (m *Memory2D) GetNod(x, y int) (*Nod, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := m.checkBounds(x, y); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	nod := m.mem[y][x]
	if nod == nil {
		return nil, fmt.Errorf("no node at coordinates (%d,%d)", x, y)
	}

	return nod, nil
}

// GetHed retrieves a head by coordinates
func (m *Memory2D) GetHed(x, y int) (*Hed, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if err := m.checkBounds(x, y); err != nil {
		return nil, fmt.Errorf("invalid coordinates: %w", err)
	}

	coordStr := fmt.Sprintf("%d,%d", x, y)
	for _, hed := range m.heds {
		if hed.ID() == coordStr {
			return hed, nil
		}
	}

	return nil, fmt.Errorf("no head at coordinates (%d,%d)", x, y)
}

// Bang triggers all heads
func (m *Memory2D) Bang() []error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errors []error
	for _, hed := range m.heds {
		if err := hed.Bang(); err != nil {
			errors = append(errors, fmt.Errorf("head %s error: %w", hed.ID(), err))
		}
	}

	return errors
}

// Dimensions returns the size of the grid
func (m *Memory2D) Dimensions() (rows, cols int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.mem), len(m.mem[0])
}

// checkBounds verifies coordinates are within the grid
func (m *Memory2D) checkBounds(x, y int) error {
	if y < 0 || y >= len(m.mem) || x < 0 || x >= len(m.mem[0]) {
		return fmt.Errorf("coordinates (%d,%d) out of bounds", x, y)
	}
	return nil
}

// Helper functions to create IDs
func NodID(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func HedID(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func (m *Memory2D) GetHeads() []*Hed {
	// No need for locking since we're just reading a snapshot
	// and don't need perfect consistency for visualization
	heads := make([]*Hed, len(m.heds))
	copy(heads, m.heds)
	return heads
}

func (m *Memory2D) GetGrid() [][]*Nod {
	grid := make([][]*Nod, len(m.mem))
	for i := range m.mem {
		grid[i] = make([]*Nod, len(m.mem[i]))
		copy(grid[i], m.mem[i])
	}
	return grid
}

func (m *Memory2D) ClearMemory() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get current dimensions
	rows, cols := len(m.mem), len(m.mem[0])

	// Recreate empty grid with same dimensions
	m.mem = make([][]*Nod, rows)
	for i := range m.mem {
		m.mem[i] = make([]*Nod, cols)
	}

	// Clear heads slice
	m.heds = make([]*Hed, 0)
}
