package util

import (
	"fmt"
	"hash/fnv"
)

func HashCode(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func HashCodeAsString(s string) string {
	return fmt.Sprintf("%d", HashCode(s))
}
