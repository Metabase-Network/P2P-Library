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

package hash

import (
	"github.com/Metabase-Network/vasuki/crypto"

	"golang.org/x/crypto/blake2b"
)

var (
	_ crypto.cryptoInterface = (*hash)(nil)
)

func New() *hash {
	return &hash{}
}

// HashBytes hashes the given bytes using the blakle2b hash algorithm.

func Hash(bytes []byte) []byte {
	src := blake2b.Sum256(bytes)
	return src
}

/*
--Encode in String--
func Hash(bytes []byte) []byte {
	src := blake2b.Sum256(bytes)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst[:], src[:])
	return dst
}

*/
