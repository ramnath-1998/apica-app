package handlers

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ramnath.1998/apica-app/models"
)

type CacheStruct struct {
	Cache map[int]models.Node
	Lock  sync.Mutex
}

var (
	AppCache   = &CacheStruct{}
	Count      int
	Capacity   int
	Head, Tail models.Node
	Lock       sync.Mutex
)

func (cache *CacheStruct) SetCache() {

	AppCache = cache

}

func (cache *CacheStruct) RemoveExpired() {

	count := 0
	for {
		cache.Lock.Lock()

		time.Sleep(1 * time.Second)
		count++
		log.Println("Count", count+1)
		for key, value := range cache.Cache {
			if time.Since(value.IssuedAt) >= time.Duration(value.Expiration)*time.Second {
				delete(cache.Cache, key)
				fmt.Println("Removed expired key:", key)
			}
		}
		cache.Lock.Unlock()
	}
}

func (cache *CacheStruct) GetTheListFromHeadToTail() []models.Node {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	node := Head
	resultList := []models.Node{}

	for range cache.Cache {
		nextNode, exist := cache.Cache[node.NextNode.Key]
		if exist {
			resultList = append(resultList, models.Node{Key: nextNode.Key, Value: nextNode.Value, IssuedAt: nextNode.IssuedAt, Expiration: nextNode.Expiration})
		}
		node = nextNode
	}

	return resultList
}

func (cache *CacheStruct) NodeUpdateCache(node *models.Node) {

	if node.NextNode == nil || node.PreviousNode == nil {
		return
	} else {
		cache.Cache[node.Key] = *node
	}

}

func (cache *CacheStruct) AddNode(node models.Node) models.Node {

	nextOfHead := Head.NextNode
	nextOfHead.PreviousNode = &node
	node.NextNode = Head.NextNode
	Head.NextNode = &node
	node.PreviousNode = &Head
	cache.NodeUpdateCache(nextOfHead)
	cache.NodeUpdateCache(&node)

	return node
}

func (cache *CacheStruct) RemoveNode(node models.Node) {

	previousNode := node.PreviousNode
	nextNode := node.NextNode

	previousNode.NextNode = nextNode
	nextNode.PreviousNode = previousNode
	cache.NodeUpdateCache(previousNode)
	cache.NodeUpdateCache(nextNode)
	delete(cache.Cache, node.Key)
}

func (cache *CacheStruct) MoveToHead(node models.Node) models.Node {

	cache.RemoveNode(node)
	node = cache.AddNode(node)

	return node
}

func (cache *CacheStruct) PopTail() {

	nodeToBePopped := Tail.PreviousNode
	cache.RemoveNode(*nodeToBePopped)

}

func (cache *CacheStruct) InitializeCache(capacity int) {

	cache.Cache = make(map[int]models.Node)
	Count = 0
	Capacity = capacity
	Head = models.Node{}
	Head.PreviousNode = nil

	Tail = models.Node{}
	Tail.NextNode = nil

	Head.NextNode = &Tail
	Tail.PreviousNode = &Head
	cache.SetCache()

	fmt.Println(cache)
	fmt.Println(AppCache.Cache)

}

func (cache *CacheStruct) UpdateCache(key int, value int, expiration int) models.Node {
	Lock.Lock()
	defer Lock.Unlock()

	node, exists := cache.Cache[key]

	if exists {
		node.Value = value
		node.IssuedAt = time.Now().Local().UTC()
		node = cache.MoveToHead(node)
	} else {
		Count = Count + 1
		if Count > Capacity {
			cache.PopTail()
		}
		newNode := models.Node{
			Key:        key,
			Value:      value,
			IssuedAt:   time.Now().UTC(),
			Expiration: expiration,
		}
		node = cache.AddNode(newNode)

	}

	return models.Node{
		Key:        node.Key,
		Value:      node.Value,
		Expiration: node.Expiration,
		IssuedAt:   node.IssuedAt,
	}
}
