package api

import (
	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
)

type ticketService interface {
	CreateTicket(ticket *model.Ticket) error
}

type TicketAPI struct {
	service ticketService
}

func NewTicketAPI(service ticketService) *TicketAPI {
	return &TicketAPI{
		service: service,
	}
}

func (a *TicketAPI) CreateTicket(ctx *gin.Context) {
	ticket := &model.Ticket{}
	if err := ctx.ShouldBindJSON(ticket); err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	if err := a.service.CreateTicket(ticket); err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(201, ticket)
}
