package gcsimage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"errors"

	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"image/jpeg"
	"io/ioutil"
)

type Bucket struct {
	handle *storage.BucketHandle
}

func InitBucket(bucket string) (*Bucket, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Bucket{handle: client.Bucket(bucket)}, nil
}

func (b *Bucket) Get(id string, width, height int) ([]byte, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s-%d-%d", id, width, height)

	reader, err := b.handle.Object(key).NewReader(ctx)
	if err == nil {
		return ioutil.ReadAll(reader)
	}

	reader, err = b.handle.Object(id).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	original, errImg := jpeg.Decode(reader)
	if errImg != nil {
		return nil, errImg
	}

	modified := imaging.Resize(original, width, height, imaging.Lanczos)
	buf := new(bytes.Buffer)
	errEnc := jpeg.Encode(buf, modified, nil)
	if errEnc != nil {
		return nil, errEnc
	}

	data := buf.Bytes()
	errSave := b.save(key, data)
	if errSave != nil {
		return nil, errSave
	}

	return data, nil
}

func (b *Bucket) Add(data []byte) (string, error) {
	sniffFormat()
	fixOrientation()

	id := uuid.New().String()
	err := b.save(id, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *Bucket) save(key string, data []byte) error {
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	ctx := context.Background()
	writer := b.handle.Object(key).NewWriter(ctx)
	writer.ContentType = "image/jpeg"
	_, errWrite := writer.Write(data)
	if errWrite != nil {
		return errWrite
	}

	errClose := writer.Close()
	if errClose != nil {
		return errClose
	}

	return nil
}

func sniffFormat() {
	//TODO implement
	// https://golang.org/src/image/format.go
}

func fixOrientation() {
	//TODO implement
	//  imaging . readOrientation(r io.Reader) orientation
}
