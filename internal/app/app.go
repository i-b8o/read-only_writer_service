package app

import (
	"context"
	"fmt"
	"net"
	"read-only_writer_service/internal/config"
	chapter_controller "read-only_writer_service/internal/controller/chapter"
	doc_controller "read-only_writer_service/internal/controller/doc"
	paragraph_controller "read-only_writer_service/internal/controller/paragraph"
	postgressql "read-only_writer_service/internal/data_providers/db/postgresql"
	"read-only_writer_service/internal/domain/service"
	chapter_usecase "read-only_writer_service/internal/domain/usecase/chapter"
	doc_usecase "read-only_writer_service/internal/domain/usecase/doc"
	paragraph_usecase "read-only_writer_service/internal/domain/usecase/paragraph"
	"read-only_writer_service/pkg/client/postgresql"
	"time"

	"github.com/i-b8o/logging"
	pb_writable "github.com/i-b8o/read-only_contracts/pb/writer/v1"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	logger     logging.Logger
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	logger := logging.GetLogger(config.AppConfig.LogLevel)

	logger.Print("Postgres initializing")
	pgConfig := postgresql.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := postgresql.NewClient(context.Background(), 5, time.Second*5, pgConfig)
	if err != nil {
		logger.Fatal(err)
	}
	regAdapter := postgressql.NewDocStorage(pgClient)
	chapterAdapter := postgressql.NewChapterStorage(pgClient)
	paragraphAdapter := postgressql.NewParagraphStorage(pgClient)

	docService := service.NewDocService(regAdapter)
	chapterService := service.NewChapterService(chapterAdapter)
	paragraphService := service.NewParagraphService(paragraphAdapter)

	docUsecase := doc_usecase.NewDocUsecase(docService, chapterService, paragraphService)
	chapterUsecase := chapter_usecase.NewChapterUsecase(chapterService)
	paragraphUsecase := paragraph_usecase.NewParagraphUsecase(paragraphService)

	docGrpcService := doc_controller.NewWriterDocGrpcController(docUsecase, logger)
	chapterGrpcService := chapter_controller.NewWriterChapterGrpcController(chapterUsecase, logger)
	paragraphGrpcService := paragraph_controller.NewWritableDocGRPCService(paragraphUsecase, logger)
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
	pb_writable.RegisterWriterDocGRPCServer(grpcServer, docGrpcService)
	pb_writable.RegisterWriterChapterGRPCServer(grpcServer, chapterGrpcService)
	pb_writable.RegisterWriterParagraphGRPCServer(grpcServer, paragraphGrpcService)

	return App{cfg: config, grpcServer: grpcServer, logger: logger}, nil
}

func (a *App) Run(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", a.cfg.GRPC.IP, a.cfg.GRPC.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	a.logger.Printf("started server on %s", address)
	return a.grpcServer.Serve(listener)

}
