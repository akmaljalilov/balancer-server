package consistenHashing

import (
	"fmt"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"strconv"
	"velox/server/src/service"
	"velox/server/src/service/replic"
)

/*
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
*/
//
var RingDB *consistent.Consistent

type myMember string

func (m myMember) String() string {
	return string(m)
}

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}
func Start() {
	cfg := consistent.Config{
		PartitionCount:    4,
		ReplicationFactor: 2,
		Load:              1.25,
		Hasher:            hasher{},
	}
	members := []consistent.Member{myMember("localhost:3307")}

	for i := 8; i < 11; i++ {
		d := strconv.Itoa(i)
		if i < 10 {
			d = "0" + d
		}

		member := myMember(fmt.Sprintf("localhost:33%s", d))
		members = append(members, member)
	}
	RingDB = consistent.New(members, cfg)
}

func GetServerByKey(key string) *replic.Server {
	server := RingDB.LocateKey([]byte(key))
	return service.MapServer[server.String()]
}
