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

package crypto

import (
	"encoding/hex"
	"errors"
)

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

var (
	PrivateKeySizeErr = errors.New("private key length does not equal expected key length")
)

func (k *KeyPair) Sign(ci cryptoInterface, message []byte) ([]byte, error) {
	if len(k.PrivateKey) != ci.PrivateKeySize() {
		return nil, PrivateKeySizeErr
	}

	message = ci.HashBytes(message)

	signature := ci.Sign(k.PrivateKey, message)
	return signature, nil
}

// PrivateKeyHex returns the hex representation of the private key.
func (k *KeyPair) PrivateKeyHex() string {
	return hex.EncodeToString(k.PrivateKey)
}

// PublicKeyHex returns the hex representation of the public key.
func (k *KeyPair) PublicKeyHex() string {
	return hex.EncodeToString(k.PublicKey)
}

// String returns the private and public key pair.
func (k *KeyPair) String() (string, string) {
	return k.PrivateKeyHex(), k.PublicKeyHex()
}

// FromPrivateKey returns a KeyPair given a signature policy and private key.
func FromPrivateKey(sp cryptoInterface, privateKey string) (*KeyPair, error) {
	rawPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	return fromPrivateKeyBytes(sp, rawPrivateKey)
}

func fromPrivateKeyBytes(ci cryptoInterface, rawPrivateKey []byte) (*KeyPair, error) {
	if len(rawPrivateKey) != ci.PrivateKeySize() {
		return nil, PrivateKeySizeErr
	}

	rawPublicKey, err := ci.PrivateToPublic(rawPrivateKey)
	if err != nil {
		return nil, err
	}

	keyPair := &KeyPair{
		PrivateKey: rawPrivateKey,
		PublicKey:  rawPublicKey,
	}

	return keyPair, nil
}

// Verify returns true if the given signature was generated using the given public key, message, signature policy, and hash policy.
func Verify(ci cryptoInterface, publicKey []byte, message []byte, signature []byte) bool {
	// Public key must be a set size.
	if len(publicKey) != ci.PublicKeySize() {
		return false
	}

	message = ci.HashBytes(message)
	return ci.Verify(publicKey, message, signature)
}
