package errhand

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vpoletaev11/fileHostingSite/test"
)

func TestErrhand(t *testing.T) {
	w := httptest.NewRecorder()

	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	InternalError(fmt.Errorf("testing error"), w)
	writer.Close()

	assert.Equal(t, time.Now().Format("2006/01/02 15:04:05")+" github.com/vpoletaev11/fileHostingSite/errhand.TestErrhand():38 INTERNAL ERROR: testing error\n", <-out)

	test.AssertBodyEqual(t, "INTERNAL ERROR. Please try later\n", w.Body)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
