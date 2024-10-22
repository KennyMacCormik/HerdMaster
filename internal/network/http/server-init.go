package http

import (
	"github.com/KennyMacCormik/HerdMaster/internal/network/http/routes"
	"github.com/gin-gonic/gin"
)

func initGin(MaxConn int) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(clientConnLimiter(MaxConn))
	_ = router.SetTrustedProxies(nil)

	routes.DogsHandlers(router)
	routes.OwnersHandlers(router)

	return router
}

// clientConnLimiter limits the number of goroutines actually doing the job.
// Neither gin nor http.Server allows to prevent goroutines from spawning, but we can hold them.
func clientConnLimiter(MaxConn int) func(c *gin.Context) {
	limiter := make(chan struct{}, MaxConn)
	return func(c *gin.Context) {
		for {
			select {
			case limiter <- struct{}{}:
				c.Next()
				<-limiter
				return
			}
		}
	}
}
