package meta

import (
	"github.com/noah-blockchain/coinExplorer-tools/models"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
)

type Resource struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IconUrl     *string `json:"icon_url"`
	SiteUrl     *string `json:"site_url"`
}

func (r Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)

	return Resource{
		Name:        validator.Name,
		Description: validator.Description,
		IconUrl:     validator.IconUrl,
		SiteUrl:     validator.SiteUrl,
	}
}
