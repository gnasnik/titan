package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/linguohua/titan/node"
	"github.com/linguohua/titan/node/modules/dtypes"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/linguohua/titan/api"
	"github.com/linguohua/titan/build"
	lcli "github.com/linguohua/titan/cli"
	cliutil "github.com/linguohua/titan/cli/util"
	"github.com/linguohua/titan/lib/titanlog"
	"github.com/linguohua/titan/lib/ulimit"
	"github.com/linguohua/titan/metrics"
	"github.com/linguohua/titan/node/config"
	"github.com/linguohua/titan/node/repo"
	"github.com/quic-go/quic-go/http3"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/tag"
	"golang.org/x/xerrors"
)

var log = logging.Logger("main")

const FlagLocatorRepo = "locator-repo"

// TODO remove after deprecation period
const FlagLocatorRepoDeprecation = "locatorrepo"

func main() {
	api.RunningNodeType = api.NodeLocator
	titanlog.SetupLogLevels()

	local := []*cli.Command{
		runCmd,
	}

	local = append(local, lcli.CommonCommands...)

	app := &cli.App{
		Name:                 "titan-locator",
		Usage:                "Titan locator node",
		Version:              build.UserVersion(),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagLocatorRepo,
				Aliases: []string{FlagLocatorRepoDeprecation},
				EnvVars: []string{"TITAN_LOCATION_PATH", "LOCATION_PATH"},
				Value:   "~/.titanlocator", // TODO: Consider XDG_DATA_HOME
				Usage:   fmt.Sprintf("Specify locator repo path. flag %s and env TITAN_EDGE_PATH are DEPRECATION, will REMOVE SOON", FlagLocatorRepoDeprecation),
			},
			&cli.StringFlag{
				Name:    "panic-reports",
				EnvVars: []string{"TITAN_PANIC_REPORT_PATH"},
				Hidden:  true,
				Value:   "~/.titanlocator", // should follow --repo default
			},
		},

		After: func(c *cli.Context) error {
			if r := recover(); r != nil {
				// Generate report in TITAN_LOCATOR_PATH and re-raise panic
				build.GeneratePanicReport(c.String("panic-reports"), c.String(FlagLocatorRepo), c.App.Name)
				panic(r)
			}
			return nil
		},
		Commands: append(local, lcli.LocationCmds...),
	}
	app.Setup()
	app.Metadata["repoType"] = repo.Locator

	if err := app.Run(os.Args); err != nil {
		log.Errorf("%+v", err)
		return
	}
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "Start titan locator node",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "geodb-path",
			Usage: "geodb path, example: --geodb-path=../../geoip/geolite2_city/city.mmdb",
			Value: "../../geoip/geolite2_city/city.mmdb",
		},
		&cli.StringFlag{
			Name:  "accesspoint-db",
			Usage: "mysql db, example: --accesspoint-db=user01:sql001@tcp(127.0.0.1:3306)/test",
			Value: "user01:sql001@tcp(127.0.0.1:3306)/test",
		},
	},

	Before: func(cctx *cli.Context) error {
		return nil
	},
	Action: func(cctx *cli.Context) error {
		log.Info("Starting titan locator node")

		limit, _, err := ulimit.GetLimit()
		switch {
		case err == ulimit.ErrUnsupported:
			log.Errorw("checking file descriptor limit failed", "error", err)
		case err != nil:
			return xerrors.Errorf("checking fd limit: %w", err)
		default:
			if limit < build.DefaultFDLimit {
				return xerrors.Errorf("soft file descriptor limit (ulimit -n) too low, want %d, current %d", build.DefaultFDLimit, limit)
			}
		}

		repoPath := cctx.String(FlagLocatorRepo)
		r, err := repo.NewFS(repoPath)
		if err != nil {
			return err
		}

		ok, err := r.Exists()
		if err != nil {
			return err
		}
		if !ok {
			if err := r.Init(repo.Locator); err != nil {
				return err
			}
		}

		lr, err := r.Lock(repo.Locator)
		if err != nil {
			return err
		}

		cfg, err := lr.Config()
		if err != nil {
			return err
		}

		locatorCfg := cfg.(*config.LocatorCfg)

		err = lr.Close()
		if err != nil {
			return err
		}

		var locatorAPI api.Locator
		stop, err := node.New(cctx.Context,
			node.Locator(&locatorAPI),
			node.Base(),
			node.Repo(r),
			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("accesspoint-db") },
				node.Override(new(dtypes.DatabaseAddress), func() dtypes.DatabaseAddress {
					return dtypes.DatabaseAddress(cctx.String("accesspoint-db"))
				})),
			node.ApplyIf(func(s *node.Settings) bool { return cctx.IsSet("geodb-path") },
				node.Override(new(dtypes.GeoDBPath), func() dtypes.GeoDBPath {
					return dtypes.GeoDBPath(cctx.String("geodb-path"))
				})),
		)
		if err != nil {
			return xerrors.Errorf("creating node: %w", err)
		}

		address := locatorCfg.ListenAddress
		addrSplit := strings.Split(address, ":")
		if len(addrSplit) < 2 {
			return fmt.Errorf("listen address %s is error", address)
		}

		ctx := lcli.ReqContext(cctx)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		handler := LocatorHandler(locatorAPI, true)
		srv := &http.Server{
			Handler: handler,
			BaseContext: func(listener net.Listener) context.Context {
				ctx, _ := tag.New(context.Background(), tag.Upsert(metrics.APIInterface, "titan-edge"))
				return ctx
			},
		}

		udpPacketConn, err := net.ListenPacket("udp", address)
		if err != nil {
			return err
		}
		defer udpPacketConn.Close()

		httpClient := cliutil.NewHttp3Client(udpPacketConn, locatorCfg.InsecureSkipVerify, locatorCfg.CaCertificatePath)
		jsonrpc.SetHttp3Client(httpClient)

		go startUDPServer(udpPacketConn, handler, locatorCfg)

		go func() {
			<-ctx.Done()
			log.Warn("Shutting down...")
			if err := srv.Shutdown(context.TODO()); err != nil {
				log.Errorf("shutting down RPC server failed: %s", err)
			}
			stop(ctx)
			// udpPacketConn.Close()
			log.Warn("Graceful shutdown successful")
		}()

		nl, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}

		log.Infof("Titan locator server listen on %s", address)

		return srv.Serve(nl)
	},
}

func startUDPServer(conn net.PacketConn, handler http.Handler, locatorCfg *config.LocatorCfg) error {
	var tlsConfig *tls.Config
	if locatorCfg.InsecureSkipVerify {
		config, err := generateTLSConfig()
		if err != nil {
			log.Errorf("startUDPServer, generateTLSConfig error:%s", err.Error())
			return err
		}
		tlsConfig = config
	} else {
		cert, err := tls.LoadX509KeyPair(locatorCfg.CaCertificatePath, locatorCfg.PrivateKeyPath)
		if err != nil {
			log.Errorf("startUDPServer, LoadX509KeyPair error:%s", err.Error())
			return err
		}

		tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: false,
		}
	}

	srv := http3.Server{
		TLSConfig: tlsConfig,
		Handler:   handler,
	}

	return srv.Serve(conn)
}

func generateTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		InsecureSkipVerify: true,
	}, nil
}
