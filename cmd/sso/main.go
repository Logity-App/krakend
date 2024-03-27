package main

import (
	"fmt"
	ssogrpc "github.com/Logity-App/sso/internal/clients/sso/grpc"
	config "github.com/Logity-App/sso/internal/pkg/congig"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
)

func main() {
	// TODO : add logger

	// TODO : add validation

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		panic("fail sso")
	}

	ssoClient, err := ssogrpc.New(
		context.Background(),
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		fmt.Println(err.Error())
		panic("fail sso")

	}

	// Create Gin router
	router := gin.Default()

	// TODO : вынести handlers
	recipesHandler := NewRecipesHandler(ssoClient)
	// Register Routes
	router.GET("/sso/verify-new-phone-number/:phone", recipesHandler.VerifyNewPhoneNumber)
	router.POST("/sso/send-sms-code", recipesHandler.SendSmsCode)
	router.POST("/sso/sign-up-by-phone", recipesHandler.SignUpByPhone)
	router.POST("/sso/verify-phone-number", recipesHandler.VerifyPhoneNumber)
	router.POST("/sso/sign-in-by-phone", recipesHandler.SignInByPhone)

	// Start the server
	router.Run()
}

type SsoHandler struct {
	client *ssogrpc.Client
}

func NewRecipesHandler(client *ssogrpc.Client) *SsoHandler {
	return &SsoHandler{
		client: client,
	}
}

func (h *SsoHandler) VerifyNewPhoneNumber(c *gin.Context) {
	phone := c.Param("phone")

	res, err := h.client.VerifyNewPhoneNumber(context.Background(), phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

type SmsRequest struct {
	Phone   string `json:"phone" binding:"required"`
	SmsCode string `json:"smsCode" binding:"required"`
}

func (h *SsoHandler) SendSmsCode(c *gin.Context) {
	var req SmsRequest
	c.BindJSON(&req)
	res, err := h.client.SendSmsCode(context.Background(), req.Phone, req.SmsCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

type SignUpByPhoneRequest struct {
	Phone        string `json:"phone" binding:"required"`
	BirthdayDate string `json:"birthdayDate" binding:"required"`
	DefaultTag   string `json:"defaultTag" binding:"required"`
}

func (h *SsoHandler) SignUpByPhone(c *gin.Context) {
	var req SignUpByPhoneRequest
	c.BindJSON(&req)

	res, err := h.client.SignUpByPhone(context.Background(), req.Phone, req.BirthdayDate, req.DefaultTag)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

type VerifyPhoneNumberRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *SsoHandler) VerifyPhoneNumber(c *gin.Context) {
	var req VerifyPhoneNumberRequest
	c.BindJSON(&req)

	res, err := h.client.VerifyPhoneNumber(context.Background(), req.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

type SignInByPhoneRequest struct {
	Phone   string `json:"phone" binding:"required"`
	SmsCode string `json:"smsCode" binding:"required"`
}

func (h *SsoHandler) SignInByPhone(c *gin.Context) {
	var req SignInByPhoneRequest
	c.BindJSON(&req)

	res, err := h.client.SignInByPhone(context.Background(), req.Phone, req.SmsCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}
