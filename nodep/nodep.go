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

package nodep

import (
	"encoding/hex"

	"github.com/Metabase-Network/vasuki/crypto"
	"github.com/syndtr/goleveldb/leveldb"
)

// nodePks  struct for nodepks
type nodePks struct {
	StorPath string
	Stordb   *leveldb.DB
}

// start function to initilize nodep
func start(path string) nodePks {
	var ret nodePks
	db, err := leveldb.OpenFile(path, nil)
	defer db.Close()
	if err == nil {
		ret = nodePks{StorPath: path, Stordb: db}
	} else {
		ret = nodePks{}
	}
	return ret
}

// FetchKey function to fetch keys from local storage
func (n nodePks) fetchKey() string {
	var ret string
	if n.Stordb != nil {
		res, err := n.Stordb.Get([]byte("vasukiNodeKey"), nil)
		if err != nil {
			ret = hex.EncodeToString(res)
		} else {
			ret = ""
		}
	}
	return ret
}

// isKey function to check weather valid key is applied or not
func (n nodePks) iskey() bool {
	var ret1 bool
	if n.Stordb != nil {
		res, err := n.Stordb.Get([]byte("vasukiNodeKey"), nil)
		if err != nil {
			ret := hex.EncodeToString(res)
			_, err1 := crypto.HexToECDSA(ret)
			if err1 != nil {
				ret1 = true
			} else {
				ret1 = false
			}
		} else {
			ret1 = false
		}
	}
	return ret1
}

// setKey function to apply a key if key is not present
func (n nodePks) setkey() bool {
	var ret bool
	if n.Stordb != nil {
		key, _ := crypto.GenerateKey()
		err := n.Stordb.Put([]byte("vasukiNodeKey"), crypto.FromECDSA(key), nil)
		if err != nil {
			ret = true
		} else {
			ret = false
		}
	} else {
		ret = false
	}
	return ret
}
