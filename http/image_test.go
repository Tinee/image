package http

import (
	"encoding/json"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Tinee/prog-image"
	"github.com/Tinee/prog-image/mock"
)

func TestImageHandler_handlePost(t *testing.T) {
	fakeID := "FakeSuccessID"
	svcMock := getMockStorage()
	png := getMockImage("github.png")
	h := Handler{}
	h.ImageHandler = &ImageHandler{
		Storage: svcMock,
	}

	svcMock.SaveImageFn = func(r io.Reader) (string, error) {
		return "FakeSuccessID", nil
	}

	req, err := http.NewRequest(http.MethodPost, "/image", png)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !svcMock.SaveImageInvoked {
		t.Fatal("Expected SaveImageInvoked to have been invoked")
	}
	if rec.Code != http.StatusOK {
		t.Fatal("Expected status code to be OK( 200 )")
	}

	res := make(map[string]string)
	json.NewDecoder(rec.Body).Decode(&res)
	if res["data"] != fakeID {
		t.Fatalf("Expected response DATA to be %v but got %v", fakeID, res["data"])
	}
}

func TestImageHandler_handleGetWithIDAndConvertToJPEG(t *testing.T) {
	empty := strings.NewReader("")
	bs, _ := ioutil.ReadAll(getMockImage("github.png"))
	fakeImage := progimage.Image{
		Body:        bs,
		ID:          "fakeImageID",
		ContentType: "image/png",
	}

	svcMock := getMockStorage()
	h := Handler{}
	h.ImageHandler = &ImageHandler{
		Storage: svcMock,
	}
	svcMock.GetFn = func(id string) (*progimage.Image, error) {
		return &fakeImage, nil
	}

	req, err := http.NewRequest(http.MethodGet, "/image?id="+fakeImage.ID+".jpeg", empty)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected code to be %v but got %v", http.StatusOK, rec.Code)
	}

	if !svcMock.GetInvoked {
		t.Fatal("Expected GetInvoked to have been invoked")
	}
	var response struct {
		Data progimage.Image
	}
	json.NewDecoder(rec.Body).Decode(&response)

	if response.Data.ContentType != "image/jpeg" {
		t.Fatalf("Expected contentType to be image/jpeg but got %v", response.Data.ContentType)
	}
}

func TestImageHandler_handleGetWithIDAndConvertToPNG(t *testing.T) {
	empty := strings.NewReader("")
	bs, _ := ioutil.ReadAll(getMockImage("github.jpg"))
	fakeImage := progimage.Image{
		Body:        bs,
		ID:          "fakeImageID",
		ContentType: "image/jpeg",
	}

	svcMock := getMockStorage()
	h := Handler{}
	h.ImageHandler = &ImageHandler{
		Storage: svcMock,
	}
	svcMock.GetFn = func(id string) (*progimage.Image, error) {
		return &fakeImage, nil
	}

	req, err := http.NewRequest(http.MethodGet, "/image?id="+fakeImage.ID+".png", empty)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected code to be %v but got %v", http.StatusOK, rec.Code)
	}

	if !svcMock.GetInvoked {
		t.Fatal("Expected GetInvoked to have been invoked")
	}
	var response struct {
		Data progimage.Image
	}
	json.NewDecoder(rec.Body).Decode(&response)

	if response.Data.ContentType != "image/png" {
		t.Fatalf("Expected contentType to be image/png but got %v", response.Data.ContentType)
	}
}

func TestImageHandler_notFound(t *testing.T) {
	h := Handler{}
	r := strings.NewReader("")
	req, err := http.NewRequest(http.MethodGet, "/somethingThatDoesNotExist", r)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected to be %v but got %v", http.StatusNotFound, rec.Code)
	}
}

func getMockStorage() *mock.Storage {
	return &mock.Storage{}
}

func getMockImage(name string) *os.File {
	f, _ := os.Open("../testdata/" + name)
	return f
}
