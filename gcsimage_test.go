package gcsimage

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

var bucket Bucket

func TestInitBucket(t *testing.T) {
	//arrange

	//act
	bucket, err := InitBucket(os.Getenv("IMAGES_STORAGE_BUCKET"))

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

	//act

	//assert
}

func TestAdd(t *testing.T) {
	//arrange
	bucket, err := InitBucket(os.Getenv("IMAGES_STORAGE_BUCKET"))
	if err != nil {
		return
	}

	cat := dataFromUrl("https://placekitten.com/500/500")
	empty := make([]byte, 0)

	//act
	good, _ := bucket.Add(cat)
	bad, _ := bucket.Add(empty)

	//assert
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
