package chart

import (
	"time"

	"github.com/noah-blockchain/noah-explorer-api/internal/helpers"
	"github.com/noah-blockchain/noah-explorer-api/internal/resource"
	"github.com/noah-blockchain/noah-explorer-api/internal/reward"
)

type RewardResource struct {
	Time   string `json:"time"`
	Amount string `json:"amount"`
}

func (RewardResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := model.(reward.ChartData)

	return RewardResource{
		Time:   data.Time.Format(time.RFC3339),
		Amount: helpers.QNoahStr2Noah(data.Amount),
	}
}
