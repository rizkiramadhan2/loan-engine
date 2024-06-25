package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simple-app/app"
	"simple-app/app/api/http/handler"

	"github.com/gin-gonic/gin"
)

type Server struct {
	handler *handler.Handler
}

type Dependencies struct {
	LoanUC app.LoanUseCase
}

var (
	s Server
)

// Init will initialize this http package
func Init(deps Dependencies) {
	// add more uc here
	h := handler.New(deps.LoanUC)

	s = Server{
		handler: h,
	}
}

// Run run http server
func Run(port string) {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "You know, for checking...")
	})

	loans := r.Group("/loans")
	loans.GET("", s.handler.GetList)
	loans.GET("/:id/detail", s.handler.GetDetail)

	loans.POST("", s.handler.CreateLoan)
	loans.PATCH("/:id/approve", s.handler.ApproveLoan)
	loans.POST("/:id/invest", s.handler.InvestLoan)
	loans.POST("/:id/disburse", s.handler.DisburseLoan)

	/* End of registering router */

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Graceful Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
