package rs

import (
	"OSS/apiServer/objectstream"
	"fmt"
	"io"
	"log"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")

	}
	log.Println("NewRSGetStream1111111")
	log.Println("locateInfo22222222", locateInfo)//map[0:127.0.0.1:54324 1:127.0.0.1:54325 2:127.0.0.1:54321 3:127.0.0.1:54322 4:127.0.0.1:54323 5:127.0.0.1:54326]
	log.Println("dataServers333333", dataServers)//为空
	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		log.Println("server11", server)
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue

		}
		// TODO 导致reader为空
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		log.Println("reader222", reader)
		if e == nil {
			readers[i] = reader

		}
		log.Println("reader3333", reader)
	}
	log.Println("readers111", readers)
	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		if readers[i] == nil {
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e

			}

		}

	}
	log.Println("writers2222", writers)
	log.Println("writers2222", writers[0])
	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil

}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)

		}

	}

}

func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")

	}
	if offset < 0 {
		panic("only support forward seek")

	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset

		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length

	}
	return offset, nil

}
