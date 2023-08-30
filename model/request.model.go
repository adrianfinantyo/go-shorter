package model

type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ShortURL string `json:"short_url" binding:"required"`
}