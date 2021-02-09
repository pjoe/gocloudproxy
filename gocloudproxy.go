package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
	"gocloud.dev/gcerrors"
)

var storageURL string

func main() {
	storageURL = os.Getenv("STORAGE_URL")
	if len(os.Args) > 1 {
		storageURL = os.Args[1]
	}
	fmt.Printf("STORAGE_URL: %v\n", storageURL)

	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello gocloudproxy")
	})

	http.HandleFunc("/", serveBlobs)

	appPort := os.Getenv("PORT")
	if len(appPort) == 0 {
		appPort = "8080"
	}
	fmt.Printf("listening on port %s\n", appPort)
	err := http.ListenAndServe(":"+appPort, nil)
	if err != nil {
		panic(err)
	}
}

func serveBlobs(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)

	path := strings.TrimPrefix(r.URL.Path, "/")
	if len(path) == 0 {
		path = "index.html"
	}
	status := http.StatusOK
	ctx := r.Context()

	bucket, err := blob.OpenBucket(ctx, storageURL)
	if err != nil {
		status = http.StatusInternalServerError
		http.Error(w, err.Error(), status)
		return
	}
	defer bucket.Close()

	attribs, err := bucket.Attributes(ctx, path)
	if err != nil {
		status = http.StatusInternalServerError
		errCode := gcerrors.Code(err)
		if errCode == gcerrors.NotFound || len(path) == 0 {
			status = http.StatusNotFound
		}
	} else {
		ifNoneMatch := r.Header.Get("If-None-Match")
		ifModifiedSince, err := time.Parse(time.RFC1123, r.Header.Get("If-Modified-Since"))
		if len(ifNoneMatch) > 0 && attribs.ETag == ifNoneMatch {
			status = http.StatusNotModified
		} else if err == nil && !attribs.ModTime.After(ifModifiedSince) {
			status = http.StatusNotModified
		} else {
			setStrHeader(w, "Cache-Control", &attribs.CacheControl)
			setStrHeader(w, "Content-Type", &attribs.ContentType)
			setStrHeader(w, "Content-Disposition", &attribs.ContentDisposition)
			setStrHeader(w, "Content-Encoding", &attribs.ContentEncoding)
			setStrHeader(w, "Content-Language", &attribs.ContentLanguage)
			setIntHeader(w, "Content-Length", &attribs.Size)
			setTimeHeader(w, "Last-Modified", &attribs.ModTime)
			setStrHeader(w, "ETag", &attribs.ETag)
		}
	}

	if status == http.StatusOK {
		w.WriteHeader(status)
		reader, err := bucket.NewReader(ctx, path, nil)
		if err != nil {
			status = http.StatusInternalServerError
			fmt.Printf("Error for '%s': %d, %s\n", path, status, err.Error())
		}
		defer reader.Close()

		io.Copy(w, reader)
	} else if status == http.StatusNotModified {
		w.WriteHeader(status)
	} else {
		noCache := "no-cache"
		setStrHeader(w, "Cache-Control", &noCache)
		w.WriteHeader(status)
		w.Write([]byte(http.StatusText(status)))
	}
}

func setStrHeader(w http.ResponseWriter, key string, value *string) {
	if value != nil && len(*value) > 0 {
		w.Header().Add(key, *value)
	}
}

func setIntHeader(w http.ResponseWriter, key string, value *int64) {
	if value != nil && *value > 0 {
		w.Header().Add(key, strconv.FormatInt(*value, 10))
	}
}

func setTimeHeader(w http.ResponseWriter, key string, value *time.Time) {
	if value != nil && !reflect.DeepEqual(*value, time.Time{}) {
		w.Header().Add(key, value.UTC().Format(http.TimeFormat))
	}
}
