package main

import
(
	"commentsService/cmd/crud/app"
	"commentsService/pkg/crud/services/comments"
	"context"
	"flag"
	"github.com/AbduvokhidovRustamzhon/mux2/pkg/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	hostF   = flag.String("host", "", "Server host")
	portF   = flag.String("port", "", "Server port")
	dsnF    = flag.String("dsn", "", "Postgres DSN")
)
var (
	EHOST   = "HOST"
	EPORT   = "PORT"
	EDSN    = "DATABASE_URL"
)

func main() {
	flag.Parse()
	host, ok := FlagOrEnv(*hostF, EHOST)
	if !ok {
		log.Panic("can't get host")
	}
	port, ok := FlagOrEnv(*portF, EPORT)
	if !ok {
		log.Panic("can't get port")
	}
	dsn, ok := FlagOrEnv(*dsnF, EDSN)
	if !ok {
		log.Panic("can't get dsn")
	}
	addr := net.JoinHostPort(host, port)
	start(addr, dsn)
}

func start(addr string, dsn string) {
	router := mux.NewExactMux()
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	templatesPath := filepath.Join("web", "templates")
	assetsPath := filepath.Join("web", "assets")
	mediaPath := filepath.Join("web", "media")

	commentsSvc := comments.NewCommentsSvc(pool)
	server := app.NewServer(
		router,
		pool,
		commentsSvc,
		templatesPath,
		assetsPath,
		mediaPath,
	)
	server.InitRoutes()

	panic(http.ListenAndServe(addr, server))
}
func FlagOrEnv(flag string, envKey string) (string, bool) {
	if flag != "" {
		return flag, true
	}
	return os.LookupEnv(envKey)
}