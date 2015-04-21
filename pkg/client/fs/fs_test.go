/*
 * Mini Copy, (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this fs except in compliance with the License.
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

package fs

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/minio-io/check"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestList(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "fs-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root)

	objectPath := filepath.Join(root, "object1")
	fsc := New(objectPath)

	data := "hello"
	binarySum := md5.Sum([]byte(data))
	etag := base64.StdEncoding.EncodeToString(binarySum[:])
	dataLen := int64(len(data))

	writer, err := fsc.Put(etag, dataLen)
	c.Assert(err, IsNil)

	size, err := io.CopyN(writer, bytes.NewBufferString(data), dataLen)
	c.Assert(err, IsNil)
	c.Assert(size, Equals, dataLen)

	objectPath = filepath.Join(root, "object2")
	fsc = New(objectPath)

	writer, err = fsc.Put(etag, dataLen)
	c.Assert(err, IsNil)

	size, err = io.CopyN(writer, bytes.NewBufferString(data), dataLen)
	c.Assert(err, IsNil)
	c.Assert(size, Equals, dataLen)

	fsc = New(root)
	items, err := fsc.List()
	c.Assert(err, IsNil)
	c.Assert(len(items), Equals, 2)
}

func (s *MySuite) TestPutBucket(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "fs-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root)

	bucketPath := filepath.Join(root, "bucket")
	fsc := New(bucketPath)
	err = fsc.PutBucket()
	c.Assert(err, IsNil)
}

func (s *MySuite) TestStatBucket(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "fs-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root)

	bucketPath := filepath.Join(root, "bucket")
	fsc := New(bucketPath)
	err = fsc.PutBucket()
	c.Assert(err, IsNil)
	err = fsc.Stat()
	c.Assert(err, IsNil)
}

func (s *MySuite) TestPutObject(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "fs-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root)

	objectPath := filepath.Join(root, "object")
	fsc := New(objectPath)

	data := "hello"
	binarySum := md5.Sum([]byte(data))
	etag := base64.StdEncoding.EncodeToString(binarySum[:])
	dataLen := int64(len(data))

	writer, err := fsc.Put(etag, dataLen)
	c.Assert(err, IsNil)

	size, err := io.CopyN(writer, bytes.NewBufferString(data), dataLen)
	c.Assert(err, IsNil)
	c.Assert(size, Equals, dataLen)
}

func (s *MySuite) TestGetObject(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "fs-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root)

	objectPath := filepath.Join(root, "object")
	fsc := New(objectPath)

	data := "hello"
	binarySum := md5.Sum([]byte(data))
	etag := base64.StdEncoding.EncodeToString(binarySum[:])
	dataLen := int64(len(data))

	writer, err := fsc.Put(etag, dataLen)
	c.Assert(err, IsNil)

	_, err = io.CopyN(writer, bytes.NewBufferString(data), dataLen)
	c.Assert(err, IsNil)

	reader, size, _, err := fsc.Get()
	c.Assert(err, IsNil)
	var results bytes.Buffer

	_, err = io.CopyN(&results, reader, size)
	c.Assert(err, IsNil)
	c.Assert([]byte(data), DeepEquals, results.Bytes())

}

func (s *MySuite) TestGetObjectMetadata(c *C) {
	root, err := ioutil.TempDir(os.TempDir(), "fs-")
	c.Assert(err, IsNil)
	defer os.RemoveAll(root)

	objectPath := filepath.Join(root, "object")
	fsc := New(objectPath)

	data := "hello"
	binarySum := md5.Sum([]byte(data))
	etag := base64.StdEncoding.EncodeToString(binarySum[:])
	dataLen := int64(len(data))

	writer, err := fsc.Put(etag, dataLen)
	c.Assert(err, IsNil)

	_, err = io.CopyN(writer, bytes.NewBufferString(data), dataLen)
	c.Assert(err, IsNil)

	item, err := fsc.GetObjectMetadata()
	c.Assert(err, IsNil)
	c.Assert(item.Name, Equals, "object")
	c.Assert(item.Size, Equals, dataLen)
}