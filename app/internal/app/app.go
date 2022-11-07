package app

import (
	"context"
	"fmt"
	"net"
	postgressql "regulations_writable_service/internal/adapters/db/postgresql"
	"regulations_writable_service/internal/config"
	"regulations_writable_service/internal/pb"
	"regulations_writable_service/internal/service"
	"regulations_writable_service/pkg/client/postgresql"
	"time"

	"github.com/i-b8o/logging"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logger := logging.GetLogger(config.Logger.LogLevel)

	logger.Print("Postgres initializing")
	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.PostgreUsername, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}
	chapterAdapter := postgressql.NewChapterStorage(pgClient)
	paragraphAdapter := postgressql.NewParagraphStorage(pgClient)
	regAdapter := postgressql.NewRegulationStorage(pgClient)

	regulationGrpcService := service.NewWritableRegulationGRPCService(regAdapter, chapterAdapter, paragraphAdapter, logger)

	// read ca's cert, verify to client's certificate
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// caPem, err := ioutil.ReadFile(homeDir + "/certs/ca-cert.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create cert pool and append ca's cert
	// certPool := x509.NewCertPool()
	// if !certPool.AppendCertsFromPEM(caPem) {
	// 	log.Fatal(err)
	// }

	// // read server cert & key
	// serverCert, err := tls.LoadX509KeyPair(homeDir+"/certs/server-cert.pem", homeDir+"/certs/server-key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // configuration of the certificate what we want to
	// conf := &tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// 	ClientAuth:   tls.RequireAndVerifyClientCert,
	// 	ClientCAs:    certPool,
	// }

	// //create tls certificate
	// tlsCredentials := credentials.NewTLS(conf)

	// grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))
	grpcServer := grpc.NewServer()
	pb.RegisterWritableRegulationGRPCServer(grpcServer, regulationGrpcService)

	return App{cfg: config, grpcServer: grpcServer}, nil
}

func (a *App) Run(ctx context.Context) error {
	address := fmt.Sprintf("%s:%s", a.cfg.GRPC.IP, a.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return a.grpcServer.Serve(listener)

}
