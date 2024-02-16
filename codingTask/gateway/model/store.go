package model

type Store struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name" binding:"required"`
	Address   string `json:"address" db:"address" binding:"required"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type Version struct {
	ID            int    `json:"id" db:"id"`
	StoreID       int    `json:"store_id" db:"store_id" binding:"required"`
	VersionNumber int    `json:"version_number" db:"version_number" binding:"required"`
	Creator       string `json:"creator" db:"creator" binding:"required"`
	Owner         string `json:"owner" db:"owner" binding:"required"`
	OpenAt        string `json:"open_at" db:"open_at" binding:"required"`
	CloseAt       string `json:"close_at" db:"close_at" binding:"required"`
	CreatedAt     string `json:"created_at" db:"created_at"`
}
