package consistenHashing

import "math"

type ConsistentHashing struct {
	Ring       []string
	Nodes      []string
	Weights    map[string]int
	sortedKeys []string
}

func NewConsistentHashing() *ConsistentHashing {
	consh := &ConsistentHashing{}
	//consh.
	return consh
}

func (c ConsistentHashing) generateCircle() {
	// Generate the circle

	totalWeight := 0
	for _, node := range c.Nodes {
		weight := c.Weights[node]
		if weight == 0 {
			weight = 1
		}
		totalWeight += weight
		factor := math.Floor(float64((40 * len(c.Nodes) * weight) / totalWeight))
		println(factor)
	}

}
