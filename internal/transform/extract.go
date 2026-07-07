package transform

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sxijyoti/whiskey/internal/parser"
)

var defuddleClient = &http.Client{}

func FetchMarkdown(
	target string,
) ([]byte, error) {

	endpoint :=
		"https://defuddle.md/" +
		url.PathEscape(target)
	
	fmt.Println("[defuddle] fetching:", endpoint)
	
	req, err := http.NewRequest(
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"Whiskey/1.0",
	)

	resp, err := defuddleClient.Do(
		req,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf(
			"defuddle: %s",
			resp.Status,
		)
	}

	body, err := io.ReadAll(
		resp.Body,
	)
	if err != nil {
		return nil, err
	}

	doc, err := parser.ParseFrontmatter(
		body,
	)

	if err != nil {

		// Not a Whiskey-style frontmatter.
		// Just return the body unchanged.
		return body, nil
	}

	return []byte(
		doc.Body,
	), nil
}