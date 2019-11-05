package rest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/pubgo/g/db/mysql"
	"github.com/pubgo/g/db/nosql/redis_g"
	"github.com/pubgo/g/ginx/controller/base/alipay"
	"github.com/pubgo/g/ginx/controller/base/hotnews"
	"github.com/pubgo/g/ginx/controller/dtcp/garnet_services"
	"github.com/pubgo/g/ginx/controller/dtcp/metadata_service"
	"github.com/pubgo/g/ginx/controller/dtcp/signature/signature_service"
	"github.com/pubgo/g/ginx/controller/miks/miks_lib/groupslib"
	"github.com/pubgo/g/ginx/controller/miks/miks_lib/searchlib"
	"github.com/pubgo/g/ginx/controller/miks/miks_mod/miks_cache"
	"github.com/pubgo/g/ginx/middleware/sessions"
	"github.com/pubgo/g/pkg/iosapns"
	"github.com/pubgo/g/pkg/storage"
	"github.com/pubgo/g/pkg/storage/elk"
	"github.com/pubgo/g/pm/cmd/miks_foresee/config"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/tylerb/graceful.v1"
)

var (
	server       *http.Server
	listener     net.Listener
	gracefulFlag = flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
)

func RunApi(cfgFile string) error {
	/**
	初始化配置
	*/
	err := config.InitConfig(cfgFile)
	if err != nil {
		return err
	}
	cfg := config.GetInstance()
	if cfg == nil {
		return errors.New("the global config has not initialized yet")
	}

	if cfg.SvcDomain == "" {
		fmt.Println("the global config svcdomain is empty")
		return errors.New("the global config svcdomain is empty")
	}

	cfg.AppDebug = true
	gin.SetMode(gin.DebugMode)
	fmt.Println("DebuMode:", strings.ToLower(cfg.GINMode))
	if strings.EqualFold(strings.ToLower(cfg.GINMode), "release") {
		cfg.AppDebug = false
		gin.SetMode(gin.ReleaseMode)
	}

	/**
	初始化日志

	/**
	初始化 mysql
	*/
	err = mysql.InitConnDB(cfg)
	if err != nil {
		fmt.Println("InitConnDB error:", err.Error())
		return err
	}

	err = elk.ElasticInit()
	if err != nil {
		fmt.Println("ElasticInit error:", err.Error())
		//return err
	}

	// start article score list
	searchlib.InitScoreList()

	// init cache
	miks_cache.InitCache()

	//miks_lib.InitCount()
	// start storage local or s3 or oss
	err = storage.StorageInit()
	if err != nil {
		fmt.Println("StorageInit start error:", err.Error())
		return err
	}

	hotnews.InitHotNewsQueue()
	// start redis service
	err = redis_g.InitRedis()
	if err != nil {
		fmt.Println("StartRedisService start error:", err.Error())
		return err
	}

	// init ios apns message queue
	iosapns.ApnsSendQueueInit()

	// session manager
	err = sessions.InitSessionManage(cfg.RedisHost, cfg.RedisPwd, cfg.RedisDb, cfg.RedisMaxidle, cfg.RedisMaxactive)
	if err != nil {
		fmt.Println("InitSessionManage start error:", err.Error())
		return err
	}

	//signature_service.UpdateSignatureOldData()

	signature_service.UpdateDnahashFromOldData()

	// signature service
	err = signature_service.BeginSignNoticeService()
	if err != nil {
		fmt.Println("BeginSignNoticeService error:", err.Error())
		return err
	}

	// send metadata
	err = metadata_service.SendMetadataServiceInit()
	if err != nil {
		fmt.Println("SendMetadataServiceInit error:", err.Error())
		return err
	}

	garnet_services.InitServerDiscover()

	//bind middleware and router
	r := gin.New()
	// allow origin

	// router

	//Init Node Group
	if err := groupslib.InitNodeDefaultAdmin(); err != nil {
		log.Println("InitNodeDefaultAdmin:", err)
	}

	if _, err := groupslib.InitNodeDefaultGroup(); err != nil {
		log.Println("InitNodeDefaultGroup:", err)
	}

	// start serving
	server := &graceful.Server{
		Timeout: 300 * time.Second,
		Server: &http.Server{
			Addr:         cfg.SvcHost,
			Handler:      r,
			ReadTimeout:  1900 * time.Second,
			WriteTimeout: 1900 * time.Second,
			IdleTimeout:  1900 * time.Second,
			//ErrorLog:     log.New(logger.LogFile, "\r\n[http]", log.Ldate|log.Ltime),
		},
		//Logger: log.New(logger.LogFile, "\r\n[graceful]", log.Ldate|log.Ltime),
	}
	if *gracefulFlag {
		log.Print("main: Listening to existing file descriptor 3.")
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
	} else {
		log.Print("main: Listening on a new file descriptor.")
		listener, err = net.Listen("tcp", server.Addr)
	}

	if err != nil {
		log.Fatalf("listener error: %v", err)
	}

	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Fatalf("listener error: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//http.HandleFunc("/v1/handle/level", logTool.LogLevel.ServeHTTP)

	http.HandleFunc("/v1/payment/notify", func(rep http.ResponseWriter, req *http.Request) {
		var noti, _ = alipay.GetTradeNotification(req)
		if noti != nil {
			alipay.TradeFinish(noti.OutTradeNo, noti)
			fmt.Println("支付成功")
		} else {
			fmt.Println("支付失败")
		}
		alipay.AckNotification(rep) // 确认收到通知消息
	})

	go func() {
		if err := http.ListenAndServe(":9090", nil); err != nil {
			panic(err)
		}
	}()

	signalHandler(ctx)
	log.Printf("signal end")

	return err
}

func reload() error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}

	f, err := tl.File()
	if err != nil {
		return err
	}

	args := []string{"-graceful"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{f}
	return cmd.Start()
}

func signalHandler(ctx context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-ch
		log.Printf("signal: %v", sig)

		// timeout context for shutdown

		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// stop
			log.Printf("stop")
			signal.Stop(ch)
			if server != nil {
				err := server.Shutdown(ctx)
				if err != nil {
					log.Printf("signalHandler SIGINT Shutdown:%v", err.Error())
					return
				}
			}
			//locTaskSendQueue := tasks.GetSendQueueObj()
			//locTaskSendQueue.Stop(ctx)
			log.Printf("graceful shutdown")
			return
		case syscall.SIGUSR2:
			// reload
			log.Printf("reload")
			err := reload()
			if err != nil {
				log.Fatalf("graceful restart error: %v", err)
			}
			if server != nil {
				err := server.Shutdown(ctx)
				if err != nil {
					log.Printf("signalHandler SIGUSR2 Shutdown:%v", err.Error())
					return
				}
			}
			//locTaskSendQueue := tasks.GetSendQueueObj()
			//locTaskSendQueue.Stop(ctx)
			log.Printf("graceful reload")
			return
		}
	}
}
