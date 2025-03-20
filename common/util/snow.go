package util

import (
	"github.com/bwmarrin/snowflake"
)

// Snowflake represents a snowflake ID generator
type Snowflake struct {
	node *snowflake.Node
}

// NewSnowflake generates a new Snowflake instance
func NewSnowflake() *Snowflake {
	// Implement the actual snowflake generation logic here
	node, _ := snowflake.NewNode(1)
	return &Snowflake{node: node}
}

// Int64 returns the snowflake ID as an int64
func (s *Snowflake) Int64() int64 {
	return s.node.Generate().Int64()
}

// String returns the snowflake ID as a string
func (s *Snowflake) String() string {
	return s.node.Generate().String()
}
