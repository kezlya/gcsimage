package gcsimage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"io"
	"log"
)

type GCSImage struct {
	project, bucketName string
	bucket              *storage.BucketHandle
	cache               map[string]string
}

func DefaultGCSImage(project, bucket string) *GCSImage {

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	return &GCSImage{
		project:    project,
		bucketName: bucket,
		bucket:     client.Bucket(bucket),
		cache:      make(map[string]string, 0),
	}
}

func (i *GCSImage) Get(id string, width, height int) string {
	key := key(id, width, height)
	url, ok := i.cache[key]
	if ok {
		return url
	}

	if i.existInGCS(key) {
		url = getUrl(i.bucketName, key)
		i.cache[key] = url
		return url
	}

	url, ok = i.cache[id]
	if ok {
		org, err := i.read(url)
		if err != nil {
			modified := imaging.Resize(org, width, height, imaging.Lanczos)
			var w io.Writer
			imaging.Encode(w, modified, imaging.JPEG)
		}

		//Store
		//Add key
		return url
	}

	if i.existInGCS(id) {
		url = getUrl(i.bucketName, id)
		i.cache[key] = url
		//i.resize(width, height)
		//Store
		//Add key

		return url
	}

	return ""
}
func (i *GCSImage) existInGCS(key string) bool {
	_, err := i.bucket.Object(key).Attrs(context.Background())
	return err != nil
}

func (i *GCSImage) Add(file io.Reader) (err error) {

	//TODO: Add sniff format
	//i.org, err = jpeg.Decode(file)
	if err != nil {
		return err
	}
	return nil
}

func (i *GCSImage) read(key string) (image.Image, error) {
	reader, err := i.bucket.Object(key).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(reader)
}

func key(id string, width, height int) string {
	return fmt.Sprintf("%s-%d-%d", id, width, height)
}

func getUrl(bucket, key string) string {
	return fmt.Sprintf("//storage.googleapis.com/%s/%s", bucket, key)
}
