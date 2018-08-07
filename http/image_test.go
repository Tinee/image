package http

import (
	"encoding/json"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

func getMockStorage() *mock.Storage {
	return &mock.Storage{}
}

func getMockImage(name string) *os.File {
	f, _ := os.Open("../testdata/" + name)
	return f
}
