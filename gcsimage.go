package gcsimage

import (
	"bytes"
	"cloud.google.com/go/storage"
	c "context"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"image/jpeg"
	"image/png"
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

type Ext string

const (
	JPG Ext = "jpg"
	PNG     = "png"
)

type Bucket struct {
	handle  *storage.BucketHandle
	quality int
}

func InitBucket(ctx c.Context, bucket string, jpegQuality int) (*Bucket, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		handle:  client.Bucket(bucket),
		quality: jpegQuality,
	}, nil
}

func (b *Bucket) Get(ctx c.Context, id string, ext Ext, anchor Anchor, width, height int) ([]byte, error) {
	if width <= 0 && height <= 0 {
		return b.getOriginal(ctx, id)
	}

	key := fmt.Sprintf("%s-%d-%d-%s", id, width, height, ext)
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

	switch ext {
	case PNG:
		err = png.Encode(buf, modified)
	case JPG:
		err = jpeg.Encode(buf, modified, &jpeg.Options{Quality: b.quality})
	}
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()
	errSave := b.Save(ctx, key, data)
	if errSave != nil {
		return nil, errSave
	}

	return data, nil
}

func (b *Bucket) getOriginal(ctx c.Context, id string) ([]byte, error) {
	reader, err := b.handle.Object(id).NewReader(ctx)
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

func (b *Bucket) Add(ctx c.Context, data []byte) (string, error) {
	id := uuid.New().String()
	err := b.Save(ctx, id, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *Bucket) Save(ctx c.Context, key string, data []byte) error {
	if len(data) == 0 {
		return errors.New("data is empty")
	}

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
