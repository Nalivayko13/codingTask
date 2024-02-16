package response

import "github.com/Nalivayko13/codingTask/gateway/model"

type VersionResp struct {
	ID            int         `json:"id"`
	StoreID       int         `json:"store_id"`
	VersionNumber int         `json:"version_number" binding:"required"`
	Creator       string      `json:"creator" binding:"required"`
	Owner         string      `json:"owner"  binding:"required"`
	OpenAt        string      `json:"open_at" binding:"required"`
	CloseAt       string      `json:"close_at"  binding:"required"`
	CreatedAt     string      `json:"created_at" binding:"required"`
	Store         model.Store `json:"store"`
}

func CreateVersionResp(version *model.Version, store *model.Store) VersionResp {
	return VersionResp{
		ID:            version.ID,
		StoreID:       version.StoreID,
		VersionNumber: version.VersionNumber,
		Creator:       version.Creator,
		Owner:         version.Owner,
		OpenAt:        version.OpenAt,
		CloseAt:       version.CloseAt,
		CreatedAt:     version.CreatedAt,
		Store:         *store,
	}
}
