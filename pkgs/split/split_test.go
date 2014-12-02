/*
 * Mini Object Storage, (C) 2014 Minio, Inc.
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

package split

import (
	"bufio"
	"bytes"
	"log"
	"strconv"
	"testing"

	. "gopkg.in/check.v1"
)

type MySuite struct{}

var _ = Suite(&MySuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *MySuite) TestFileSplit(c *C) {
	b := Split{}
	err := b.GenChunks("TESTFILE", "1KB")
	c.Assert(err, IsNil)
}

func (s *MySuite) TestSplitStream(c *C) {
	var bytesBuffer bytes.Buffer
	bytesWriter := bufio.NewWriter(&bytesBuffer)
	for i := 0; i < 100; i++ {
		bytesWriter.Write([]byte(strconv.Itoa(i)))
	}
	bytesWriter.Flush()
	log.Println(strconv.Itoa(bytesBuffer.Len()))
	ch := make(chan ByteMessage)
	reader := bytes.NewReader(bytesBuffer.Bytes())
	go SplitStream(reader, 25, ch)
	var resultsBuffer bytes.Buffer
	resultsWriter := bufio.NewWriter(&resultsBuffer)
	for chunk := range ch {
		resultsWriter.Write(chunk.Data)
	}
	resultsWriter.Flush()
	c.Assert(bytes.Compare(bytesBuffer.Bytes(), resultsBuffer.Bytes()), Equals, 0)
}

func (s *MySuite) TestFileSplit2(c *C) {
	err := SplitFilesWithPrefix("TESTFILE", 1024, "TESTPREFIX")
	c.Assert(err, IsNil)
}
