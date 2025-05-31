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
)



func main() {
      //load config 
      cfg := config.MustLoad()

      //database setup 
   
      // setup router 
     router := http.NewServeMux()
 

     router.HandleFunc("POST /api/students",student.New())
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
   

   err:= server.Shutdown(ctx)

   if err != nil {
      slog.Error("Failed to shutdown server gracefully", slog.String("error", err.Error()))
   } else {
      slog.Info("Server shutdown gracefully")
   }

 }