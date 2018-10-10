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

package stor

import (
	"testing"
)

var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

func TestPath(t *testing.T) {
	path := "./testing"
	n, err := Start(path)
	if err != nil {
		t.Errorf("error while initialising DB")
	} else {
		t.Log("Init Test completed successfully")

		putErr := n.Put([]byte("testKey"), []byte("TestValue"))
		if putErr != nil {
			t.Errorf("Error while storing key")
		} else {
			t.Log("Put Operation successfull")
		}

		val, getErr := n.Get([]byte("testKey"))
		if getErr != nil {
			t.Errorf("Error while Fetching key")
		} else {
			t.Logf("Get Operation successfull %d", val)
		}

		has, hasErr := n.Has([]byte("testKey"))
		if hasErr != nil {
			t.Errorf("Error while checking key %t", has)
		} else {
			t.Logf("Get Operation successfull %t", has)
		}

	}
}
