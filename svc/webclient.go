package svc

import (
	"fmt"

	bus "github.com/autonomouskoi/core-tinygo"
)

// WebclientStaticDownload retrieves a file from the specified URL using a GET
// request, caching it. If the URL is already cached, no download is performed.
// The return value is the path relative to the root of the AK web service where
// the file can be downloaded from.
func WebclientStaticDownload(url string, timeoutMS uint64) (string, error) {
	req := WebclientStaticDownloadRequest{URL: url}
	b, err := req.MarshalVT()
	if err != nil {
		return "", fmt.Errorf("marshalling request: %w", err)
	}
	reply, err := bus.WaitForReply(&bus.BusMessage{
		Type:    int32(MessageType_WEBCLIENT_STATIC_DOWNLOAD_REQ),
		Message: b,
	}, timeoutMS)
	if err != nil {
		return "", fmt.Errorf("waiting for reply: %w", err)
	}
	if reply.Error != nil {
		return "", reply.Error
	}
	resp := WebclientStaticDownloadResponse{}
	if err := resp.UnmarshalVT(reply.GetMessage()); err != nil {
		return "", fmt.Errorf("unmarshalling response: %w", err)
	}
	return resp.Path, nil
}
