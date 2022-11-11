package rs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func testEncodeDecode(t *testing.T, p []byte) {
	writers := make([]io.Writer, ALL_SHARDS)
	readers := make([]io.Reader, ALL_SHARDS)
	for i := range writers {
		writers[i], _ = os.Create(fmt.Sprintf("/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ut_%d", i))

	}
	enc := NewEncoder(writers)
	length := len(p)
	for count := 0; count != length; {
		n, e := enc.Write(p[count:])
		if e != nil {
			t.Error(e)

		}
		count += n

	}
	enc.Flush()
	for i := range writers {
		writers[i].(*os.File).Close()
		writers[i] = nil
		readers[i], _ = os.Open(fmt.Sprintf("/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ut_%d", i))

	}
	readers[1] = nil
	readers[4] = nil
	writers[1], _ = os.Create("/home/ling/go/src/Orginone/pkg/minio/objectData/temp/repair_1")
	writers[4], _ = os.Create("/home/ling/go/src/Orginone/pkg/minio/objectData/temp/repair_4")
	dec := NewDecoder(readers, writers, int64(length))
	b := make([]byte, length+10)
	count := 0
	for {
		n, e := dec.Read(b[count:])
		count += n
		if e == io.EOF {
			break

		}

	}
	if count != length {
		t.Error(count, length)

	}
	if !reflect.DeepEqual(b[:count], p) {
		t.Error("not match")

	}
	output, e := exec.Command("diff", "/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ut_1", "/home/ling/go/src/Orginone/pkg/minio/objectData/temp/repair_1").Output()
	if len(output) != 0 {
		t.Error(output, e)

	}
	output, e = exec.Command("diff", "/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ut_4", "/home/ling/go/src/Orginone/pkg/minio/objectData/temp/repair_4").Output()
	if len(output) != 0 {
		t.Error(output, e)

	}
	t.Log("dec.cache", dec.cache)
	// dec.enc.Split()

	// Create shards and load the data.
	var err error
	shards := make([][]byte, ALL_SHARDS)
	for i := range shards {
		infn := fmt.Sprintf("%s%d", "/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ut_", i)
		fmt.Println("Opening", infn)
		shards[i], err = ioutil.ReadFile(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
		}
	}
	t.Log("shards", shards)

	// Join the shards and write them
	outfn := "/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ut.txt"

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)

	err = dec.enc.Join(f, shards, len(shards[0])*4)
	if err != nil {
		t.Log("err", err)
	}
	t.Log("b[:count]", string(b[:count]))
	t.Log("p", string(p))
}

func TestEncodeDecode(t *testing.T) {
	p := []byte("123456789nihk你好,.,.")
	testEncodeDecode(t, p)
	// p = []byte("123")
	// testEncodeDecode(t, p)
	// p = []byte("12345")
	// testEncodeDecode(t, p)
	// p = make([]byte, 9999)
	// fillRandom(p)
	// testEncodeDecode(t, p)
	// p = make([]byte, 99999)
	// fillRandom(p)
	// testEncodeDecode(t, p)

}

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8

		}

	}

}

func testEncodeDecode2(t *testing.T, a []byte) {
	var err error
	writers := make([]io.Writer, ALL_SHARDS)
	readers := make([]io.Reader, ALL_SHARDS)
	for i := range writers {
		writers[i], err = os.Create(fmt.Sprintf("/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ling_%d", i))
		if err != nil {
			t.Log("err", err)
		}
	}
	enc := NewEncoder(writers)
	length := len(a)
	t.Log("length111",length)
	for count := 0; count != length; {
		n, err := enc.Write(a[count:])
		if err != nil {
			t.Error(err)
		}
		count += n
	}
	enc.Flush()
	for i := range writers {
		writers[i].(*os.File).Close()
		writers[i] = nil
	}
	// --------------------------------
	for i := range writers {
		readers[i], _ = os.Open(fmt.Sprintf("/home/ling/go/src/Orginone/pkg/minio/objectData/temp/ling_%d", i))

	}
	t.Log("writers111",writers)
	dec := NewDecoder(readers, writers, int64(length))
	aa := make([]byte, length+10)
	count := 0
	for {
		n, err := dec.Read(aa[count:])
		count += n
		if err != nil {
			break
		}
	}
	if count != length {
		t.Error("err count!=length 11111111111")
	}
	t.Log("aa", aa)
	t.Log("string(aa)", string(aa))

}
func TestEncodeDecodeFuncMain(t *testing.T) {
	a:=[]byte("1.接受r.bod  y文件数据")
	testEncodeDecode2(t,a)
}
