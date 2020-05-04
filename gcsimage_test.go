package gcsimage

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAdd(t *testing.T) {
	//arrange
	image := DefaultGCSImage("p", "b")
	cat := bytes.NewReader(readerFromUrl("https://placekitten.com/500/500"))
	empty := bytes.NewReader(make([]byte, 0))

	//act
	good := image.Add(cat)
	bad := image.Add(empty)

	//assert
	if good != nil {
		t.Errorf("Can't load image")
	}
	if bad == nil {
		t.Errorf("Should not load nil")
	}
}

func TestGet(t *testing.T) {

}

func readerFromUrl(url string) []byte {
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
