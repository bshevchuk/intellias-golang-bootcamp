package parser

import "encoding/xml"

func ParseRss(data []byte) (Rss, error) {
	var rss Rss
	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}
