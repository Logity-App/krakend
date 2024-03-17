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
	router.GET("/", homePage)
	router.GET("/sso/verify-new-phone-number/:phone", recipesHandler.VerifyNewPhoneNumber)
	router.GET("/sso/send-sms-code", recipesHandler.SendSmsCode)
	router.POST("/sso/sign-up-by-phone", recipesHandler.SignUpByPhone)
	router.POST("/sso/verify-phone-number", recipesHandler.VerifyPhoneNumber)
	router.POST("/sso/sign-in-by-phone", recipesHandler.SignInByPhone)

	// Start the server
	router.Run()
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is my home page")
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

func (h *SsoHandler) SendSmsCode(c *gin.Context) {

	res, err := h.client.SendSmsCode(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

func (h *SsoHandler) SignUpByPhone(c *gin.Context) {
	phone := c.Param("phone")
	birthdayDate := c.Param("birthdayDate")
	defaultTag := c.Param("defaultTag")

	res, err := h.client.SignUpByPhone(context.Background(), phone, birthdayDate, defaultTag)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

func (h *SsoHandler) VerifyPhoneNumber(c *gin.Context) {
	phone := c.Param("phone")

	res, err := h.client.VerifyPhoneNumber(context.Background(), phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}

func (h *SsoHandler) SignInByPhone(c *gin.Context) {
	phone := c.Param("phone")
	smsCode := c.Param("smsCode")

	res, err := h.client.SignInByPhone(context.Background(), phone, smsCode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, res)
}
