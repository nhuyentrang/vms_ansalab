package models

type VideoDownLoad struct {
	StorageBucket string `json:"storageBucket,omitempty"`
	Video         string `json:"video,omitempty"`
}

type ListVideoDownLoad struct {
	VideoDownLoad []VideoDownLoad `json:"videoDownLoad,omitempty"`
}
