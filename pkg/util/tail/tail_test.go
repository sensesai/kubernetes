/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tail

import (
	"bytes"
	"strings"
	"testing"
)

func TestTail(t *testing.T) {
	line := strings.Repeat("a", blockSize)
	testBytes := []byte(line + "\n" +
		line + "\n" +
		line + "\n" +
		line + "\n" +
		line[blockSize/2:]) // incomplete line

	for c, test := range []struct {
		n     int64
		start int64
	}{
		{n: -1, start: 0},
		{n: 0, start: int64(len(line)+1) * 4},
		{n: 1, start: int64(len(line)+1) * 3},
		{n: 9999, start: 0},
	} {
		t.Logf("TestCase #%d: %+v", c, test)
		r := bytes.NewReader(testBytes)
		s, err := FindTailLineStartIndexV2(r, test.n)
		if err != nil {
			t.Error(err)
		}
		if s != test.start {
			t.Errorf("%d != %d", s, test.start)
		}
	}
}

func BenchmarkFindTailLineStartIndex(b *testing.B) {
	lineCount := 200
	lineSize := 100 // assume average log line is 100 bytes
	line := strings.Repeat("a", lineSize) + "\n"
	testBytes := make([]byte, 0, lineCount*lineSize)
	for i := 0; i < lineCount; i++ {
		testBytes = append(testBytes, line...)
	}

	b.Log("Benchmark last 10 lines")
	r := bytes.NewReader(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := FindTailLineStartIndex(r, 10)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}

func BenchmarkFindTailLineStartIndexV2(b *testing.B) {
	lineCount := 200
	lineSize := 100 // assume average log line is 100 bytes
	line := strings.Repeat("a", lineSize) + "\n"
	testBytes := make([]byte, 0, lineCount*lineSize)
	for i := 0; i < lineCount; i++ {
		testBytes = append(testBytes, line...)
	}

	b.Log("Benchmark last 10 lines")
	r := bytes.NewReader(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := FindTailLineStartIndexV2(r, 10)
		if err != nil {
			b.Error(err)
		}
	}
	b.StopTimer()
}
