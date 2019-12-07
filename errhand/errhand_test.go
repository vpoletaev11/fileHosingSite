package errhand

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrhandWithUsername(t *testing.T) {
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
	InternalError("errhand", "testErrhandWithUsername", "username", fmt.Errorf("testing error"), w)
	writer.Close()

	assert.Equal(t, time.Now().Format("2006/01/02 15:04:05")+" errhand.testErrhandWithUsername() [user: username] INTERNAL ERROR: testing error\n", <-out)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestErrhandWithoutUsername(t *testing.T) {
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
	InternalError("errhand", "testErrhandWithUsername", "", fmt.Errorf("testing error"), w)
	writer.Close()

	assert.Equal(t, time.Now().Format("2006/01/02 15:04:05")+" errhand.testErrhandWithUsername() INTERNAL ERROR: testing error\n", <-out)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Equal(t, "INTERNAL ERROR. Please try later\n", bodyString)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
