package handlers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ramnath.1998/apica-app/models"
)

var (
	Cache      map[int]models.Node
	Count      int
	Capacity   int
	Head, Tail models.Node
	Lock       sync.Mutex
)

func RemoveExpired() {
	count := 0
	for {
		Lock.Lock()
		time.Sleep(1 * time.Second)
		count++
		log.Println("Count", count+1)
		for key, value := range Cache {
			if time.Since(value.IssuedAt) >= time.Duration(value.Expiration)*time.Second {
				delete(Cache, key)
				fmt.Println("Removed expired key:", key)
			}
		}
		Lock.Unlock()
	}
}

func GetTheListFromHeadToTail() []models.Node {
	Lock.Lock()
	defer Lock.Unlock()

	node := Head
	resultList := []models.Node{}

	for node.NextNode != nil {
		node, exist := Cache[node.NextNode.Key]
		if exist {
			resultList = append(resultList, models.Node{Key: node.Key, Value: node.Value, Expiration: node.Expiration})

		}
	}

	return resultList
}

func NodeUpdateCache(node *models.Node) {

	if node.NextNode == nil || node.PreviousNode == nil {
		return
	} else {
		Cache[node.Key] = *node
	}

}

func AddNode(node models.Node) models.Node {

	nextOfHead := Head.NextNode
	nextOfHead.PreviousNode = &node
	node.NextNode = Head.NextNode
	Head.NextNode = &node
	node.PreviousNode = &Head
	NodeUpdateCache(nextOfHead)
	NodeUpdateCache(&node)

	return node
}

func RemoveNode(node models.Node) {

	previousNode := node.PreviousNode
	nextNode := node.NextNode

	previousNode.NextNode = nextNode
	nextNode.PreviousNode = previousNode
	NodeUpdateCache(previousNode)
	NodeUpdateCache(nextNode)
	delete(Cache, node.Key)
}

func MoveToHead(node models.Node) models.Node {

	RemoveNode(node)
	node = AddNode(node)

	return node
}

func PopTail() {

	nodeToBePopped := Tail.PreviousNode
	RemoveNode(*nodeToBePopped)

}

func InitializeCache(capacity int) {

	Cache = make(map[int]models.Node)
	Count = 0
	Capacity = capacity
	Head = models.Node{}
	Head.PreviousNode = nil

	Tail = models.Node{}
	Tail.NextNode = nil

	Head.NextNode = &Tail
	Tail.PreviousNode = &Head
}

func UpdateCache(key int, value int, expiration int) models.Node {
	Lock.Lock()
	defer Lock.Unlock()

	node, exists := Cache[key]

	if exists {
		node.Value = value
		node.IssuedAt = time.Now().Local().UTC()
		node = MoveToHead(node)
	} else {
		Count = Count + 1
		if Count > Capacity {
			PopTail()
		}
		newNode := models.Node{
			Key:        key,
			Value:      value,
			IssuedAt:   time.Now().UTC(),
			Expiration: expiration,
		}
		node = AddNode(newNode)

	}

	return models.Node{
		Key:        node.Key,
		Value:      node.Value,
		Expiration: node.Expiration,
		IssuedAt:   node.IssuedAt,
	}
}
