package hikivision

import (
	"encoding/xml"
)

// Define the XML structure as Go structs
type CMSearchDescription struct {
	XMLName              xml.Name             `xml:"CMSearchDescription,omitempty"`
	XMLVersion           string               `xml:"version,attr"`
	XMLNamespace         string               `xml:"xmlns,attr"`
	SearchID             string               `xml:"searchID"`
	TrackList            []TrackNVR           `xml:"trackList"`
	TimeSpanList         []TimeSpan           `xml:"timeSpanList>timeSpan"`
	MaxResults           int                  `xml:"maxResults"`
	SearchResultPosition int                  `xml:"searchResultPostion"`
	MetadataList         []MetadataDescriptor `xml:"metadataList"`
}

type TrackNVR struct {
	TrackID int `xml:"trackID"`
}
type MetadataDescriptor struct {
	MetadataDescriptor string `xml:"metadataDescriptor"`
}

// Define the XML structure as Go structs
type CMSearchResult struct {
	XMLName      xml.Name          `xml:"CMSearchResult,omitempty"`
	XMLVersion   string            `xml:"version,attr"`
	XMLNamespace string            `xml:"xmlns,attr"`
	SearchID     string            `xml:"searchID"`
	Response     bool              `xml:"responseStatus"`
	StatusStrg   string            `xml:"responseStatusStrg"`
	NumOfMatches int               `xml:"numOfMatches"`
	MatchList    []SearchMatchItem `xml:"matchList>searchMatchItem"`
}

type SearchMatchItem struct {
	SourceID               string                 `xml:"sourceID"`
	TrackID                string                 `xml:"trackID"`
	TimeSpan               TimeSpan               `xml:"timeSpan"`
	MediaSegmentDescriptor MediaSegmentDescriptor `xml:"mediaSegmentDescriptor"`
	MetadataMatches        []MetadataMatch        `xml:"metadataMatches>metadataDescriptor"`
}

type TimeSpan struct {
	StartTime string `xml:"startTime"`
	EndTime   string `xml:"endTime"`
}

type MediaSegmentDescriptor struct {
	ContentType string `xml:"contentType"`
	CodecType   string `xml:"codecType"`
	PlaybackURI string `xml:"playbackURI"`
}

type MetadataMatch struct {
	Descriptor string `xml:",chardata"`
}
