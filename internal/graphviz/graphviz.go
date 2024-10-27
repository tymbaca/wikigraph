package graphviz

import (
	"math"
	"math/rand/v2"

	"github.com/tymbaca/wikigraph/internal/model"
)

type Vector2 struct{ X, Y float32 }

// Subtract two vectors
func (v1 Vector2) Subtract(v2 Vector2) Vector2 {
	return Vector2{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

// Add two vectors
func (v1 Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

// Scale a vector by a scalar
func (v1 Vector2) Scale(scalar float32) Vector2 {
	return Vector2{X: v1.X * scalar, Y: v1.Y * scalar}
}

// Normalize the vector to unit length
func (v1 Vector2) Normalize() Vector2 {
	magnitude := float32(math.Sqrt(float64(v1.X*v1.X + v1.Y*v1.Y)))
	if magnitude == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{X: v1.X / magnitude, Y: v1.Y / magnitude}
}

// Distance between two vectors
func (v1 Vector2) Distance(v2 Vector2) float32 {
	return float32(math.Sqrt(float64((v2.X-v1.X)*(v2.X-v1.X) + (v2.Y-v1.Y)*(v2.Y-v1.Y))))
}

func ForceDirLayout(graph model.Graph) map[int]Vector2 {
	const (
		k              = 0.01   // Attractive force constant
		repulsionConst = 0.1    // Lower repulsion constant to reduce NaNs
		minDistance    = 0.01   // Small minimum distance to avoid division by zero
		gridSize       = 1000.0 // Size of the simulation area
		numIterations  = 10     // Increased number of iterations
		maxForce       = 10.0   // Maximum force cap to prevent NaNs
	)

	coolingFactor := float32(0.95) // Cooling factor to reduce forces over time

	// Initialize positions randomly
	positions := make(map[int]Vector2)
	for id := range graph {
		positions[id] = Vector2{X: rand.Float32() * gridSize, Y: rand.Float32() * gridSize}
	}

	// Main iteration loop
	for iter := 0; iter < numIterations; iter++ {
		// Initialize forces map
		forces := make(map[int]Vector2)
		for id := range graph {
			forces[id] = Vector2{X: 0, Y: 0}
		}

		// Apply repulsive forces between all pairs
		for id1, pos1 := range positions {
			for id2, pos2 := range positions {
				if id1 != id2 {
					displacement := pos1.Subtract(pos2)
					distance := displacement.Distance(Vector2{}) + minDistance // Avoid zero distance
					repulsiveForce := repulsionConst / (distance * distance)
					force := displacement.Normalize().Scale(float32(repulsiveForce))

					// Cap the force to avoid overflow
					if force.X > maxForce {
						force.X = maxForce
					}
					if force.Y > maxForce {
						force.Y = maxForce
					}

					// Accumulate repulsive force for each node
					forces[id1] = forces[id1].Add(force)
					forces[id2] = forces[id2].Add(force.Scale(-1))
				}
			}
		}

		// Apply attractive forces between connected nodes
		for id, article := range graph {
			for _, childID := range article.Childs {
				if childPos, exists := positions[childID]; exists {
					displacement := childPos.Subtract(positions[id])
					distance := displacement.Distance(Vector2{}) + minDistance // Avoid zero distance
					attractiveForce := k * distance
					force := displacement.Normalize().Scale(float32(attractiveForce))

					// Cap the force to avoid overflow
					if force.X > maxForce {
						force.X = maxForce
					}
					if force.Y > maxForce {
						force.Y = maxForce
					}

					// Accumulate attractive force for each node
					forces[id] = forces[id].Add(force)
					forces[childID] = forces[childID].Add(force.Scale(-1))
				}
			}
		}

		// Update positions with capped force and apply cooling
		for id, pos := range positions {
			force := forces[id].Scale(coolingFactor)
			positions[id] = pos.Add(force)
		}

		// Reduce the cooling factor gradually
		coolingFactor *= coolingFactor
	}

	return positions
}
