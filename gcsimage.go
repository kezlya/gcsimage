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

type Anchor int

const (
	Center Anchor = iota
	TopLeft
	Top
	TopRight
	Left
	Right
	BottomLeft
	Bottom
	BottomRight
)

type Bucket struct {
	handle  *storage.BucketHandle
	quality int
}

func InitBucket(bucket string, quality int) (*Bucket, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &Bucket{
		handle:  client.Bucket(bucket),
		quality: quality,
	}, nil
}

func (b *Bucket) Get(id string, anchor Anchor, width, height int) ([]byte, error) {
	ctx := context.Background()
	if width <= 0 && height <= 0 {
		return b.getOriginal(id)
	}

	key := fmt.Sprintf("%s-%d-%d", id, width, height)
	reader, err := b.handle.Object(key).NewReader(ctx)
	if err == nil {
		return ioutil.ReadAll(reader)
	}

	reader, err = b.handle.Object(id).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	original, errImg := imaging.Decode(reader, imaging.AutoOrientation(true))
	if errImg != nil {
		return nil, errImg
	}

	modified := imaging.Fill(original, width, height, imaging.Anchor(anchor), imaging.Lanczos)
	buf := new(bytes.Buffer)
	errEnc := jpeg.Encode(buf, modified, &jpeg.Options{Quality: b.quality})
	if errEnc != nil {
		return nil, errEnc
	}

	data := buf.Bytes()
	errSave := b.Save(key, data)
	if errSave != nil {
		return nil, errSave
	}

	return data, nil
}

func (b *Bucket) getOriginal(id string) ([]byte, error) {
	reader, err := b.handle.Object(id).NewReader(context.Background())
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, errBytes := buf.ReadFrom(reader)
	if errBytes != nil {
		return nil, errBytes
	}

	return buf.Bytes(), nil
}

func (b *Bucket) Add(data []byte) (string, error) {
	id := uuid.New().String()
	err := b.Save(id, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *Bucket) Save(key string, data []byte) error {
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	ctx := context.Background()
	writer := b.handle.Object(key).NewWriter(ctx)
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
