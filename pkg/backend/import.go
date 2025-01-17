package backend

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fupas/commons/pkg/util"
	ds "github.com/fupas/platform/pkg/platform"
	"github.com/labstack/echo/v4"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/platform"
	"google.golang.org/appengine"
)

const (
	// ImportTask route to ImportTaskEndpoint
	ImportTask = "/import"

	// full canonical route
	importTaskWithPrefix = "/_t/import"
)

type (
	// ContentMetadata keeps basic data on resource
	ContentMetadata struct {
		Size        int64
		Duration    int64
		ContentType string
		Etag        string
		Timestamp   int64
	}
)

// ImportTaskEndpoint implements async file import
func ImportTaskEndpoint(c echo.Context) error {
	var req *a.Import = new(a.Import)

	err := c.Bind(req)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		return c.NoContent(http.StatusOK)
	}

	// FIXME does it make sense to retry? If not, send StatusOK
	status := importResource(appengine.NewContext(c.Request()), req.Source, req.Dest)
	return c.NoContent(status)
}

func importResource(ctx context.Context, src, dest string) int {
	resp, err := http.Get(src)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		platform.ReportError(fmt.Errorf("can not retrieve '%s': %s", src, resp.Status))
		return http.StatusBadRequest
	}

	meta := extractMetadataFromResponse(resp)
	obj := ds.Storage().Bucket(a.BucketCDN).Object(dest)
	writer := obj.NewWriter(ctx)
	writer.ContentType = meta.ContentType
	defer writer.Close()

	// transfer using a buffer
	buffer := make([]byte, 65536)
	l, err := io.CopyBuffer(writer, resp.Body, buffer)

	// error handling & verification
	if err != nil {
		platform.ReportError(fmt.Errorf("can not transfer '%s': %v", dest, err))
		return http.StatusBadRequest
	}
	if l != meta.Size {
		platform.ReportError(fmt.Errorf("error transfering '%s': expected %d, reveived %d", src, meta.Size, l))
		return http.StatusBadRequest
	}

	// update the inventory
	parent := strings.Split(dest, "/")[0]

	temp := a.Asset{
		URI: src,
		Rel: a.ResourceTypeImport,
	}
	name := strings.Split(temp.FingerprintURI(parent), "/")[1]

	duration := int64(0) // FIXME implement it

	if err := UpdateAssetResource(ctx, name, util.Checksum(src), a.ResourceAsset, parent, dest, meta.ContentType, meta.Size, duration); err != nil {
		platform.ReportError(fmt.Errorf("error updating inventory: %v", err))
		return http.StatusBadRequest
	}

	return http.StatusOK
}

// extractMetadataFromResponse extracts the metadata from http.Response
func extractMetadataFromResponse(resp *http.Response) *ContentMetadata {
	meta := ContentMetadata{
		ContentType: resp.Header.Get("content-type"),
		Etag:        resp.Header.Get("etag"),
	}
	l, err := strconv.ParseInt(resp.Header.Get("content-length"), 10, 64)
	if err == nil {
		meta.Size = l
	}
	// expects 'Wed, 30 Dec 2020 14:14:26 GM'
	t, err := time.Parse(time.RFC1123, resp.Header.Get("date"))
	if err == nil {
		meta.Timestamp = t.Unix()
	}
	return &meta
}
