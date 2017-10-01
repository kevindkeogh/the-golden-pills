package main

type Image struct {
	URL    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
	Frames string `json:"frames"`
}

type Giphy struct {
	Type               string `json:"type"`
	Id                 string `json:"id"`
	URL                string `json:"url"`
	Tags               string `json:"tags"`
	BitlyGifUrl        string `json:"bitly_gif_url"`
	BitlyFullscreenUrl string `json:"bitly_fullscreen_url"`
	BitlyTiledUrl      string `json:"bitly_tiled_url"`
	Images             struct {
		Original               Image `json:"original"`
		FixedHeight            Image `json:"fixed_height"`
		FixedHeightStill       Image `json:"fixed_height_still"`
		FixedHeightDownsampled Image `json:fixed_height_downsampled"`
		FixedWidth             Image `json:"fixed_width"`
		FixedWidthStill        Image `json:"fixed_width_still"`
		FixedWithDownsampled   Image `json:"fixed_width_downsampled"`
	} `json:"images"`
}
