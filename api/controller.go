package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ibraheemacara/tezos-delegation-service/db"
	"github.com/ibraheemacara/tezos-delegation-service/types"
)

type Controller struct {
	db db.DBInterface
}

func NewController(db db.DBInterface) *Controller {
	return &Controller{db: db}
}

func (ctr *Controller) GetDelegations(ctx *gin.Context) {
	year := ctx.Param("year")
	if year == "" {
		delegations, err := ctr.db.GetDelegations()
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		var data []types.Delegation
		for _, delegation := range delegations {
			data = append(data, types.Delegation{
				Delegator: delegation.Delegator,
				Timestamp: delegation.Timestamp,
				Block:     delegation.Block,
				Amount:    delegation.Amount,
			})
		}
		ctx.JSON(200, types.DelegationsResponse{Delegations: data})
	} else {
		delegations, err := ctr.db.GetDelegationsByYear(year)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		var data []types.Delegation
		for _, delegation := range delegations {
			data = append(data, types.Delegation{
				Delegator: delegation.Delegator,
				Timestamp: delegation.Timestamp,
				Block:     delegation.Block,
				Amount:    delegation.Amount,
			})
		}
		ctx.JSON(200, types.DelegationsResponse{Delegations: data})
	}
}

// func (ctr *Controller) GetDelegationsByYear(ctx *gin.Context) {
// 	year := ctx.Param("year")
// 	delegations, err := ctr.db.GetDelegationsByYear(year)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{"error": "Internal server error"})
// 		return
// 	}

// 	ctx.JSON(200, delegations)
// }
