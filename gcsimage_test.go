package gcsimage

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestInitBucket(t *testing.T) {
	//arrange

	//act

	//assert
}

func TestGet(t *testing.T) {
	//arrange

	//act

	//assert
}

func TestAdd(t *testing.T) {
	//arrange
	bucket, _ := InitBucket("test")
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
		if err != nil {
			return bodyBytes
		}
	}

	return nil
}
