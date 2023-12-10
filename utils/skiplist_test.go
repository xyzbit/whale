// Copyright 2021 hardcore-os Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TestSkipListBasicCRUD(t *testing.T) {
	list := NewSkipList()
	key, val := "", ""
	maxTime := 10000
	for i := 0; i < maxTime; i++ {
		//number := rand.Intn(10000)
		key, val = fmt.Sprintf("Key%d", i), fmt.Sprintf("Val%d", i)
		_ = list.Set(key, val)
		geVal, ok := list.Get(key)
		assert.True(t, ok)
		assert.Equal(t, geVal.V, val)
	}
}

func TestSkipListBasicCRUD2(t *testing.T) {
	const n = 100
	l := NewSkipList()
	for i := 0; i < n; i++ {
		a := strconv.Itoa(i)
		_ = l.Set(a, a)
	}

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		a := strconv.Itoa(i)
		v, ok := l.Get(a)
		assert.True(t, ok)
		assert.Equal(t, a, v.V)
	}
}

// Benchmark_SkipListBasicCRUD-4   	   72715	    146770 ns/op	     142 B/op	       6 allocs/op
// Benchmark_SkipListBasicCRUD-4   	   72780	    157370 ns/op	     142 B/op	       6 allocs/op
func Benchmark_SkipListBasicCRUD(b *testing.B) {
	list := NewSkipList()
	key, val := "", ""
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key, val = fmt.Sprintf("Key%d", i), fmt.Sprintf("Val%d", i)
		_ = list.Set(key, val)
		_, _ = list.Get(key)
	}
}

func TestConcurrentBasic(t *testing.T) {
	const n = 1000
	l := NewSkipList()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a := strconv.Itoa(i)
			l.Set(a, a)
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a := strconv.Itoa(i)
			v, ok := l.Get(a)
			assert.True(t, ok)
			assert.Equal(t, a, v.V)
		}(i)
	}
	wg.Wait()
}

func Benchmark_ConcurrentBasic(b *testing.B) {
	const n = 1000
	l := NewSkipList()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a := strconv.Itoa(i)
			l.Set(a, a)
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a := strconv.Itoa(i)
			v, ok := l.Get(a)
			assert.True(b, ok)
			assert.Equal(b, a, v.V)
		}(i)
	}
	wg.Wait()
}
