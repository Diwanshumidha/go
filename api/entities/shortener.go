package entities

type ShortenerParams struct {
	UID int `uri:"uid" binding:"required"`
}

type ShortenerPost struct {
	Url string `json:"url" binding:"required"`
}
