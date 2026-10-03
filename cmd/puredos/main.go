package main

import (
 "context"
 "log"
 "os"
 "os/signal"
 "syscall"
 "time"
 "github.com/louisoff84/PureDDosSecurity/internal/agent"
)
func main() {
 cfg:=agent.LoadConfig()
 logger:=log.New(os.Stdout,"puredos ",log.LstdFlags|log.Lmicroseconds)
 a,err:=agent.New(cfg,logger); if err!=nil { logger.Fatal(err) }
 ctx,cancel:=signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM); defer cancel()
 if err:=a.Run(ctx); err!=nil { logger.Fatal(err) }
 time.Sleep(100*time.Millisecond)
}
