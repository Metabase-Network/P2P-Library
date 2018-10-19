// Copyright 2014 The Metabase Authors
// This file is part of vasuki p2p library.
//
// The vasuki library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The vasuki library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package dht

import (
	"container/list"
	"sort"
	"sync"

	"github.com/Metabase-Network/vasuki/node"
)

// BucketSize defines the NodeID, Key, and routing table data structures.
const BucketSize = 16

// RoutingTable contains one bucket list for lookups.
type RoutingTable struct {
	// Current node's ID.
	self node.Def

	buckets []*Bucket
}

// Bucket holds a list of contacts of this node.
type Bucket struct {
	*list.List
	mutex *sync.RWMutex
}

// NewBucket is a Factory method of Bucket, contains an empty list.
func NewBucket() *Bucket {
	return &Bucket{
		List:  list.New(),
		mutex: &sync.RWMutex{},
	}
}

// CreateRoutingTable is a Factory method of RoutingTable containing empty buckets.
func CreateRoutingTable(id node.Def) *RoutingTable {
	table := &RoutingTable{
		self:    id,
		buckets: make([]*Bucket, len(id.NodeID)*8),
	}
	for i := 0; i < len(id.NodeID)*8; i++ {
		table.buckets[i] = NewBucket()
	}

	table.Update(id)

	return table
}

// Self returns the ID of the node hosting the current routing table instance.
func (t *RoutingTable) Self() node.Def {
	return t.self
}

// Update moves a peer to the front of a bucket in the routing table.
func (t *RoutingTable) Update(target node.Def) {
	if len(t.self.NodeID) != len(target.NodeID) {
		return
	}
	pat := target.XorID(t.self.NodeID)
	bucketID := target.XorID(t.self.NodeID).PrefixLen()
	bucket := t.Bucket(bucketID)

	var element *list.Element

	// Find current node in bucket.
	bucket.mutex.Lock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(node.Def).Equals(target) {
			element = e
			break
		}
	}

	if element == nil {
		// Populate bucket if its not full.
		if bucket.Len() <= BucketSize {
			bucket.PushFront(target)
		}
	} else {
		bucket.MoveToFront(element)
	}

	bucket.mutex.Unlock()
}

// GetPeers returns a randomly-ordered, unique list of all peers within the routing network (excluding itself).
func (t *RoutingTable) GetPeers() (peers []node.Def) {
	visited := make(map[string]struct{})
	visited[t.self.PublicKeyHex()] = struct{}{}

	for _, bucket := range t.buckets {
		bucket.mutex.RLock()

		for e := bucket.Front(); e != nil; e = e.Next() {
			id := e.Value.(node.Def)
			if _, seen := visited[id.PublicKeyHex()]; !seen {
				peers = append(peers, id)
				visited[id.PublicKeyHex()] = struct{}{}
			}
		}

		bucket.mutex.RUnlock()
	}

	return
}

// GetPeerAddresses returns a unique list of all peer addresses within the routing network.
func (t *RoutingTable) GetPeerAddresses() (peers []string) {
	visited := make(map[string]struct{})
	visited[t.self.PublicKeyHex()] = struct{}{}

	for _, bucket := range t.buckets {
		bucket.mutex.RLock()

		for e := bucket.Front(); e != nil; e = e.Next() {
			id := e.Value.(node.Def)
			if _, seen := visited[id.PublicKeyHex()]; !seen {
				peers = append(peers, id.Address)
				visited[id.PublicKeyHex()] = struct{}{}
			}
		}

		bucket.mutex.RUnlock()
	}

	return
}

// RemovePeer removes a peer from the routing table with O(bucket_size) time complexity.
func (t *RoutingTable) RemovePeer(target node.Def) bool {
	bucketID := target.XorID(t.self).PrefixLen()
	bucket := t.Bucket(bucketID)

	bucket.mutex.Lock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(node.Def).Equals(target) {
			bucket.Remove(e)

			bucket.mutex.Unlock()
			return true
		}
	}

	bucket.mutex.Unlock()

	return false
}

// PeerExists checks if a peer exists in the routing table with O(bucket_size) time complexity.
func (t *RoutingTable) PeerExists(target node.Def) bool {
	bucketID := target.XorID(t.self).PrefixLen()
	bucket := t.Bucket(bucketID)

	bucket.mutex.Lock()

	defer bucket.mutex.Unlock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(node.Def).Equals(target) {
			return true
		}
	}

	return false
}

// FindClosestPeers returns a list of k(count) peers with smallest XorID distance.
func (t *RoutingTable) FindClosestPeers(target node.Def, count int) (peers []node.Def) {
	if len(t.self.NodeID) != len(target.NodeID) {
		return []node.Def{}
	}
	bucketID := target.XorID(t.self.NodeID).PrefixLen()
	bucket := t.Bucket(bucketID)

	bucket.mutex.RLock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		peers = append(peers, e.Value.(node.Def))
	}

	bucket.mutex.RUnlock()

	for i := 1; len(peers) < count && (bucketID-i >= 0 || bucketID+i < len(t.self.Id)*8); i++ {
		if bucketID-i >= 0 {
			other := t.Bucket(bucketID - i)
			other.mutex.RLock()
			for e := other.Front(); e != nil; e = e.Next() {
				peers = append(peers, e.Value.(node.Def))
			}
			other.mutex.RUnlock()
		}

		if bucketID+i < len(t.self.NodeID)*8 {
			other := t.Bucket(bucketID + i)
			other.mutex.RLock()
			for e := other.Front(); e != nil; e = e.Next() {
				peers = append(peers, e.Value.(node.Def))
			}
			other.mutex.RUnlock()
		}
	}

	// Sort peers by XorID distance.
	sort.Slice(peers, func(i, j int) bool {
		left := peers[i].XorID(target)
		right := peers[j].XorID(target)
		return left.Less(right)
	})

	if len(peers) > count {
		peers = peers[:count]
	}

	return peers
}

// Bucket returns a specific Bucket by ID.
func (t *RoutingTable) Bucket(id int) *Bucket {
	if id >= 0 && id < len(t.buckets) {
		return t.buckets[id]
	}
	return nil
}
