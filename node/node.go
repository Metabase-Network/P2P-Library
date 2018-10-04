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
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"

	"github.com/Metabase-Network/vasuki/common"
	"github.com/Metabase-Network/vasuki/crypto"
)

var errInvalidPubkey = errors.New("Error Generating Key for Node ")

//nodeDef Struct for node
type nodeDef struct {
	nodeID     []byte
	NodeAddr   []byte
	NodePubKey ecdsa.PublicKey
	NodePvkKey *ecdsa.PrivateKey
}

//InitNode Initialise Node
func InitNode(hex string) (nodeDef, error) {
	res, err := setPrivateKey(hex)
	if err != nil {
		return nodeDef{}, err
	}
	return nodeDef{NodePvkKey: res, NodePubKey: res.PublicKey, NodeAddr: CalcNodeAddr(res.PublicKey).Bytes(), nodeID: CalcNodeID(CalcNodeAddr(res.PublicKey))}, nil
}

//setPrivateKey Sets Private Keys for the Node
func setPrivateKey(hexkey string) (*ecdsa.PrivateKey, error) {
	res, err := crypto.HexToECDSA(hexkey)
	return res, err
}

//CalcNodeAddr generates NodeAddress
func CalcNodeAddr(puk ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(puk)
}

//CalcNodeID generates NodeID
func CalcNodeID(addr common.Address) []byte {
	return crypto.Keccak256(addr.Bytes())
}

//Equals Compares 2 Node ID
func (node nodeDef) Equals(other []byte) bool {
	return bytes.Equal(node.nodeID, other)
}

//XorID XOR's ID
func (node nodeDef) XorID(other []byte) []byte {
	result := make([]byte, len(node.nodeID))

	for i := 0; i < len(node.nodeID) && i < len(other); i++ {
		result[i] = node.nodeID[i] ^ other[i]
	}
	return result
}

//AddressHex Converts the address to hex String
func (node nodeDef) AddressHex() string {
	return hex.EncodeToString(node.NodeAddr)
}

//IDHex Converts the Node ID to Hex String
func (node nodeDef) IDHex() string {
	return hex.EncodeToString(node.nodeID)
}

//ExportPvk Exports private keys in hex format
func (node nodeDef) ExportPvk() string {
	return hex.EncodeToString(crypto.FromECDSA(node.NodePvkKey))
}
