package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SOMBIT4/students-api-using-GO/internal/config"
	"github.com/SOMBIT4/students-api-using-GO/internal/http/handlers"
	
	"github.com/SOMBIT4/students-api-using-GO/internal/storage/sqlite"
)



func main() {
      //load config 
      cfg := config.MustLoad()

      //database setup 
   
      storage, err:= sqlite.New(cfg)

      if err != nil {
         log.Fatal(err)
      }

      slog.Info(("Storage initialized successfully"), slog.String("env",cfg.Env),slog.String("version","1.0.0"))
      // setup router 
     router := http.NewServeMux()
 

     router.HandleFunc("POST /api/students",student.New(storage))
     router.HandleFunc("GET /api/students/{id}", student.GetbyId(storage))
     router.HandleFunc("GET /api/students", student.GetList(storage))
     router.HandleFunc("PUT /api/students/{id}", student.UpdateById(storage))
     router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))
     // setup server
   server :=http.Server{
      Addr: cfg.Addr,
      Handler: router,
   }
   slog.Info("server stated",slog.String("address", cfg.Addr))
   fmt.Printf("Server is running on %s",cfg.HTTPServer.Addr)
   

   done := make(chan os.Signal,1 )
   signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) 
   go func ()  {
      err:= server.ListenAndServe()
   if err != nil {
      log.Fatal("Failed to start server: ")
   }
   }()
   
   <-done 
   
   slog.Info("Shutting down server...")
  
  ctx, cancel:= context.WithTimeout(context.Background(), 5 * time.Second)
   defer cancel()
   

   err= server.Shutdown(ctx)

   if err != nil {
      slog.Error("Failed to shutdown server gracefully", slog.String("error", err.Error()))
   } else {
      slog.Info("Server shutdown gracefully")
   }

 }