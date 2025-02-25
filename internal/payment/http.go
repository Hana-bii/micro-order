package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}

}

func (p *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", p.handleWebhook)
}

func (p *PaymentHandler) handleWebhook(c *gin.Context) {
	logrus.Info("Got webhook from stripe")
}
