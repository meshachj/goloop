/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package db

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testDatabase_GetSetDelete(t *testing.T, backend BackendType) {
	dir, err := ioutil.TempDir("", string(backend))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	testDB, _ := openDatabase(backend, "test", dir)
	defer testDB.Close()

	key := []byte("hello")
	key2 := []byte("hell")
	value := []byte("world")

	bucket, _ := testDB.GetBucket("hello")

	// check it has nothing before set
	has, err := bucket.Has(key)
	assert.NoError(t, err)
	assert.False(t, has)

	// SET valid value
	err = bucket.Set(key, value)
	assert.NoError(t, err)

	// GET returns same value
	result, _ := bucket.Get(key)
	assert.Equal(t, value, result)

	// HAS returns true
	has, err = bucket.Has(key)
	assert.NoError(t, err)
	assert.True(t, has)

	// DELETE value
	err = bucket.Delete(key)
	assert.NoError(t, err)

	// GET returns nothing
	result, err = bucket.Get(key)
	assert.NoError(t, err)
	assert.Nil(t, result)

	// HAS returns false
	has, err = bucket.Has(key)
	assert.NoError(t, err)
	assert.False(t, has)

	// SET with empty bytes
	err = bucket.Set(key2, []byte{})
	assert.NoError(t, err)

	// HAS returns false
	has, err = bucket.Has(key2)
	assert.NoError(t, err)
	assert.True(t, has)
}

func TestDatabase_GetSetDelete(t *testing.T) {
	for be, _ := range backends {
		t.Run(string(be), func(t *testing.T) {
			testDatabase_GetSetDelete(t, be)
		})
	}
}

func testDatabase_SetReopenGet(t *testing.T, backend BackendType) {
	if backend == MapDBBackend {
		return
	}
	dir, err := ioutil.TempDir("", string(backend))
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	key := []byte("hello")
	key2 := []byte("hell")
	value := []byte("world")

	buckets := []BucketID{"hello", MerkleTrie, BytesByHash}
	testDB, err := openDatabase(backend, "test", dir)
	assert.NoError(t, err)
	defer func() {
		if testDB != nil {
			testDB.Close()
		}
	}()

	for _, id := range buckets {
		bucket, err := testDB.GetBucket(id)
		assert.NoError(t, err)
		err = bucket.Set(key, value)
		assert.NoError(t, err)
	}
	err = testDB.Close()
	testDB = nil
	assert.NoError(t, err)

	testDB, err = openDatabase(backend, "test", dir)

	for _, id := range buckets {
		bucket, err := testDB.GetBucket(id)
		assert.NoError(t, err)
		stored, err := bucket.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, stored)

		stored, err = bucket.Get(key2)
		assert.NoError(t, err)
		assert.Nil(t, stored)

		has, err := bucket.Has(key2)
		assert.NoError(t, err)
		assert.False(t, has)
	}
}

func TestDatabase_SetReopenGet(t *testing.T) {
	for be, _ := range backends {
		t.Run(string(be), func(t *testing.T) {
			testDatabase_SetReopenGet(t, be)
		})
	}
}
