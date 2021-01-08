package gcsimage

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

var background = context.Background()

func TestInitBucket(t *testing.T) {
	//arrange

	//act
	bucket, err := InitBucket(background, os.Getenv("IMAGES_STORAGE_BUCKET"), 85)

	//assert
	if err != nil {
		log.Fatalln(err)
	}
	if bucket == nil {
		log.Fatalln("fail connect to gcs bucket")
	}
}

func TestGet(t *testing.T) {
	//arrange
	bucket, _ := InitBucket(background, os.Getenv("IMAGES_STORAGE_BUCKET"), 85)

	//act
	goodJPG, ok := bucket.Get(background, "cat", JPG, TopRight, 10, 10)
	goodPNG, ok := bucket.Get(background, "cat", PNG, TopRight, 10, 10)
	bad, notOk := bucket.Get(background, "", JPG, TopRight, 10, 10)

	//assert
	if goodJPG == nil && ok != nil {
		t.Errorf("fail to get jpg image")
	}
	if goodPNG == nil && ok != nil {
		t.Errorf("fail to get png image")
	}

	if bad != nil && notOk == nil {
		t.Errorf("Should error on bad id")
	}
}

func TestAdd(t *testing.T) {
	//arrange
	bucket, _ := InitBucket(background, os.Getenv("IMAGES_STORAGE_BUCKET"), 85)

	cat := dataFromUrl("https://placekitten.com/500/500")
	empty := make([]byte, 0)

	//act
	err := bucket.Save(background, "cat", cat)
	good, _ := bucket.Add(background, cat)
	bad, _ := bucket.Add(background, empty)

	//assert
	if err != nil {
		t.Errorf("fail to save image")
	}

	if good == "" {
		t.Errorf("fail to add image")
	}

	if bad != "" {
		t.Errorf("Should not add empty image")
	}
}

func dataFromUrl(url string) []byte {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return bodyBytes
		}
	}

	return nil
}
