package models

import (
	"sync"
	"time"
)

type Node struct {
	Key          int       `json:"key"`
	Value        int       `json:"value"`
	PreviousNode *Node     `json:"previousNode,omitempty"`
	NextNode     *Node     `json:"nextNode,omitempty"`
	Expiration   int       `json:"expiration,omitempty"`
	IssuedAt     time.Time `json:"issuedAt,omitempty"`
}
type Cache struct {
	Cache map[int]Node
	Lock  sync.Mutex
}
