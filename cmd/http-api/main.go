package main

import (
	"context"
	"log"
	"os"

	"simple-app/app/api/http"
	"simple-app/config"
	"simple-app/internal/pkg/response"
	"simple-app/internal/pkg/sqldb"
	lnRepo "simple-app/internal/repository/loan"
	lnuc "simple-app/internal/usecase/loan"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Get()
	ctx := context.Background()

	/* initialize resource like db, elastic, redis, httpclient etc here */

	/* initialize services */

	// initialize db
	db, err := sqldb.Connect(ctx, sqldb.DBConfig{
		Driver:             "postgres",
		MasterDSN:          cfg.Databases.Postgres.Master,
		FollowerDSN:        cfg.Databases.Postgres.Slave,
		MaxOpenConnections: cfg.Databases.Postgres.MaxCon,
		Retry:              cfg.Databases.Postgres.Retry,
	})
	if err != nil {
		log.Fatal("Could not get Database connection :" + err.Error())
		return
	}

	if os.Getenv("APP_ENV") == "development" {
		response.EnableStackTrace(true)
	}

	/* initialize repo */
	loanRepo := lnRepo.New(lnRepo.Param{
		DB: db,
	})

	/* initialize usecase */
	loanUc := lnuc.New(&loanRepo)

	/* initialize http handler */

	http.Init(http.Dependencies{
		LoanUC: loanUc,
	})

	// run server
	http.Run(cfg.App.Port)
}
