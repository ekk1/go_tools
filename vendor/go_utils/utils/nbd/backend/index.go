package backend

import (
	"errors"
	"go_utils/utils/myhttp"
	"net/url"
	"strconv"
	"sync"
)

type IndexBackend struct {
	name      string
	size      int64
	c         *myhttp.HTTPClient
	blockSize int64
	lock      sync.Mutex
	url       string
}

func NewIndexBackend(name string, size int64) *IndexBackend {
	c := myhttp.NewHTTPClient()
	c.SetSendForm(true)
	return &IndexBackend{
		name,
		size,
		c,
		1 * 1024 * 1024,
		sync.Mutex{},
		"http://127.0.0.1:9999",
	}
}

func (b *IndexBackend) ReadAt(p []byte, off int64) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	endPosition := int64(len(p)) + off

	if endPosition >= b.size {
		return 0, errors.New("Size overflow")
	}
	firstBlock := off / b.blockSize
	endBlock := endPosition / b.blockSize

	var dataBuffer []byte = []byte{}

	for i := firstBlock; i <= endBlock; i++ {
		formBody := url.Values{}
		formBody.Add("name", b.name)
		formBody.Add("index", strconv.FormatInt(i, 10))
		// ret, err := b.c.SendPost(b.url+"/"+"get", formBody)
		// if err != nil {
		// 	return 0, errors.New("Failed connect backend")
		// }
		// // dataRead := ret.Data()
		// dataBuffer = bytes.Join([][]byte{dataBuffer, dataRead}, []byte{})
	}

	bufferStartIndex := off % b.blockSize

	n := copy(p, dataBuffer[bufferStartIndex:bufferStartIndex+int64(len(p))])
	return n, nil
}

func (b *IndexBackend) WriteAt(p []byte, off int64) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	endPosition := int64(len(p)) + off
	if endPosition >= b.size {
		return 0, errors.New("Size overflow")
	}

	// var dataBuffer []byte = []byte{}

	firstBlock := off / b.blockSize
	firstBlockPrepend := off % b.blockSize

	// endBlock := endPosition / b.blockSize
	// endBlockRemain := (b.blockSize - (endPosition % b.blockSize))

	if firstBlockPrepend != 0 {
		formBody := url.Values{}
		formBody.Add("name", b.name)
		formBody.Add("index", strconv.FormatInt(firstBlock, 10))
		// ret, err := b.c.SendPost(b.url+"/"+"get", formBody)
		// if err != nil {
		// 	return 0, errors.New("Failed connect backend")
		// }
		// dataRead := ret.Data()
		// dataBuffer = bytes.Join([][]byte{dataBuffer, dataRead}, []byte{})
	}

	return 0, nil
}

func (b *IndexBackend) Size() (int64, error) {
	return b.size, nil
}

func (b *IndexBackend) Sync() error {
	return nil
}
