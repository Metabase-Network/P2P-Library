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
	"github.com/syndtr/goleveldb/leveldb"
)

// S  struct for S
type S struct {
	StorPath string
	Stordb   *leveldb.DB
}

// Start initilizes S Struct
func Start(path string) (S, error) {
	var ret S
	db, err := leveldb.OpenFile(path, nil)
	ret = S{StorPath: path, Stordb: db}
	return ret, err
}

// Put Stores the keys and their respective values
func (s S) Put(key []byte, value []byte) error {
	return s.Stordb.Put(key, value, nil)
}

// Get Stores the keys and their respective values
func (s S) Get(key []byte) ([]byte, error) {
	ret, err := s.Stordb.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Has Checks if the given key is repesent in db
func (s S) Has(key []byte) (bool, error) {
	return s.Stordb.Has(key, nil)
}

// Delete deletes the key from the database
func (s S) Delete(key []byte) error {
	return s.Stordb.Delete(key, nil)
}
