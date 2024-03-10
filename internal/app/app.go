package app

import (
	"database/sql"
	"forum/internal/config"
	delivery "forum/internal/delivery/http"
	server "forum/internal/server"
	"forum/internal/service"
	"forum/internal/storage"
	"forum/pkg/storage/sqlite"
	"log"
)

func Run(cfgPath string) {

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal("config load: ", err)
	}

	db, err := sqlite.NewStorage(cfg.Database.PathToDb)
	if err != nil {
		log.Fatal("storage init: ", err)
	}
	defer db.Conn.Close()
	err = db.ApplyMigrations()
	if err != nil {
		log.Fatal("storage migrations: ", err)
	}
	log.Println("storage migrations applied")

	storages := storage.NewStorages(db.Conn)

	services := service.NewServices(storages)

	handler := delivery.NewHandler(services)

	srv := server.New(&cfg, handler.InitRoutes())
	log.Printf("server started at http://%s:%s", cfg.Server.Host, cfg.Server.Port)
	err = srv.Run()
	if err != nil {
		log.Fatal("server run: ", err)
	}
	defer func(Conn *sql.DB) {
		err = Conn.Close()
		if err != nil {

		}
	}(db.Conn)
}
