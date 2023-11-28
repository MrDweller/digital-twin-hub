package manufacturer

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	httpserver "github.com/MrDweller/digital-twin-hub/http-server"
	"github.com/gin-gonic/gin"
)

func AdminAuthorization(c *gin.Context) {
	serverMode := os.Getenv("SERVER_MODE")

	switch httpserver.ServerMode(serverMode) {
	case httpserver.UNSECURE_SERVER_MODE:
		break
	case httpserver.SECURE_SERVER_MODE:
		secureModeAdminAuthorization(c)
		break
	default:
		errorString := fmt.Sprintf("the server mode %s, has no implementation", serverMode)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errorString,
		})
		c.Abort()
		break
	}

}

func secureModeAdminAuthorization(c *gin.Context) {
	if c.Request.TLS == nil || len(c.Request.TLS.VerifiedChains) <= 0 || len(c.Request.TLS.VerifiedChains[0]) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{})
		c.Abort()
		return
	}
	var commonName = c.Request.TLS.VerifiedChains[0][0].Subject.CommonName

	if strings.Split(commonName, ".")[0] != "sysop" {
		c.JSON(http.StatusUnauthorized, gin.H{})
		c.Abort()
		return
	}
}
