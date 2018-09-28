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

package node

import (
	"crypto/ecdsa"
	"errors"

	"github.com/Metabase-Network/vasuki/common"
	"github.com/Metabase-Network/vasuki/crypto"
)

var errInvalidPubkey = errors.New("Error Generating Key for Node ")
var node Node

type Node struct {
	NodeId     []byte
	NodeAddr   []byte
	NodePubKey ecdsa.PublicKey
	NodePvkKey *ecdsa.PrivateKey
}

func InitNode(hex string) (Node, error) {
	res, err := setPrivateKey(hex)
	if err != nil {
		return Node{}, err
	}
	node = Node{NodePvkKey: res, NodePubKey: res.PublicKey, NodeAddr: CalcNodeAddr(res.PublicKey).Bytes(), NodeId: CalcNodeID(CalcNodeAddr(res.PublicKey))}
	return node, nil
}

func setPrivateKey(hexkey string) (*ecdsa.PrivateKey, error) {
	res, err := crypto.HexToECDSA(hexkey)
	return res, err
}

func CalcNodeAddr(puk ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(puk)
}

func CalcNodeID(addr common.Address) []byte {
	return crypto.Keccak256(addr.Bytes())
}
