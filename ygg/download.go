package ygg

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/juliensalinas/torrengo/core"
)

func parseDescPage(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", fmt.Errorf("could not load html response into GoQuery: %v", err)
	}

	// file url is located in the 1st <a> of the 2nd <td> of the second <tbody> of class infos-torrents
	fileURL, ok := doc.Find(".infos-torrent tbody").Eq(1).Find("tr td").Eq(1).Find("a").Eq(0).First().Attr("href")
	if !ok {
		return "", fmt.Errorf("could not find a torrent file on the description page")
	}

	return fileURL, nil
}

// FindAndDlFile authenticates user, opens the torrent description page,
// and downloads the torrent file.
// Returns the local path of downloaded torrent file.
// A user timeout is set.
func FindAndDlFile(descURL string, userID string, userPass string, timeout time.Duration, client *http.Client) (string, error) {
	// Set timeout
	client.Timeout = timeout

	// Authenticate user and create http client that handles cookie and timeout
	client, err := authUser(userID, userPass, client)
	if err != nil {
		return "", fmt.Errorf("error while authenticating: %v", err)
	}

	// Fetch url
	resp, err := core.Fetch(descURL, client)
	if err != nil {
		return "", fmt.Errorf("error while fetching url: %v", err)
	}
	defer resp.Body.Close()

	// Parse html response
	filePath, err := parseDescPage(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while parsing torrent description page: %v", err)
	}

	fileURL := "https://" + baseURL + filePath
	filePathOnDisk, err := core.DlFile(fileURL, client)
	if err != nil {
		return "", fmt.Errorf("error while downloading torrent file: %v", err)
	}

	return filePathOnDisk, nil
}
