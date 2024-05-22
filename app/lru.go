package app

import (
	"errors"
	"fmt"
	"github.com/umeshkaushik-21122000/golang-lru-cache/types"
	"sync"
	"time"
)

type LRUCacheInterface interface {
	DeleteExpired()
	Remove(key int) error
	Put(key int, value int, expiry int)
	Get(key int) int
	GetAll() []types.KeyValue
}

type doubleLinkedNode struct {
	key        int
	value      int
	expiryTime time.Time
	pre        *doubleLinkedNode
	post       *doubleLinkedNode
}

var (
	head doubleLinkedNode
	tail doubleLinkedNode
)

func addNode(node *doubleLinkedNode) {
	node.pre = &head
	node.post = head.post

	head.post.pre = node
	head.post = node
}

func removeNode(node *doubleLinkedNode) {
	pre := node.pre
	post := node.post

	pre.post = post
	post.pre = pre
}

func moveToHead(node *doubleLinkedNode) {
	removeNode(node)
	addNode(node)
}

func popTail() *doubleLinkedNode {
	res := tail.pre
	removeNode(res)
	return res
}

type LRUCache struct {
	cache    map[int]*doubleLinkedNode
	count    int
	capacity int
	mutex    sync.RWMutex
}

func Constructor(capacity int) LRUCacheInterface {
	head.post = &tail
	tail.pre = &head
	cache := make(map[int]*doubleLinkedNode)
	lru := &LRUCache{
		cache:    cache,
		count:    0,
		capacity: capacity,
		mutex:    sync.RWMutex{},
	}
	return lru
}

func (this *LRUCache) Get(key int) int {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if node, ok := this.cache[key]; !ok {
		return -1
	} else {
		moveToHead(node)
		return node.value
	}
}

func (this *LRUCache) Put(key int, value int, expiry int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if node, ok := this.cache[key]; !ok {
		newNode := doubleLinkedNode{
			key:        key,
			value:      value,
			expiryTime: time.UnixMilli(0),
		}
		if expiry != 0 {
			newNode.expiryTime = time.Now().Add(time.Duration(expiry) * time.Second)
			fmt.Println("----------------------->>", newNode.expiryTime)
		}
		this.cache[key] = &newNode
		addNode(&newNode)
		this.count++
		if this.count > this.capacity {
			tail = *popTail()
			delete(this.cache, tail.key)
			this.count--
		}

	} else {
		node.value = value
		moveToHead(node)
	}
}

func (this *LRUCache) Remove(key int) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	node, ok := this.cache[key]
	if !ok {
		return errors.New("key not found")
	}

	prev := &doubleLinkedNode{}
	post := &doubleLinkedNode{}

	if node.pre != nil {
		prev = node.pre
	}
	if node.post != nil {
		post = node.post
	}
	post.pre = prev
	prev.post = post
	delete(this.cache, node.key)
	this.count--
	return nil
}

func (this *LRUCache) DeleteExpired() {
	now := time.Now().UnixMilli()
	for key, node := range this.cache {
		if node.expiryTime.UnixMilli() == 0 {
			continue
		}
		if node.expiryTime.UnixMilli() < now {
			fmt.Println("delete expired key:", key)
			this.Remove(key)
		}
	}
}

func (this *LRUCache) GetAll() []types.KeyValue {
	temp := &head
	res := []types.KeyValue{}
	for i := 0; i <= this.count; i++ {
		if i == 0 {
			temp = temp.post
			continue
		}
		tempp := types.KeyValue{
			Key:   temp.key,
			Value: temp.value,
		}
		res = append(res, tempp)
		temp = temp.post
	}
	return res
}
