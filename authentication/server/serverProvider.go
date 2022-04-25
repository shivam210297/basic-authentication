package server

import (
	"Assignment/providers"
	"Assignment/providers/cacheprovider"
	"Assignment/providers/configprovider"
	"Assignment/providers/dbhelpprovider"
	"Assignment/providers/dbprovider"
	"Assignment/providers/middlewareprovider"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	defaultServerRequestTimeoutMinutes      = 2
	defaultServerReadHeaderTimeoutSeconds   = 30
	defaultServerRequestWriteTimeoutMinutes = 30
)

type Server struct {
	DBHelper      providers.DBHelpProvider
	PSQL          providers.PSQLProvider
	Config        providers.ConfigProvider
	httpServer    *http.Server
	cache         map[string]interface{}
	CacheProvider providers.CacheProvider
	Middlewares   providers.MiddlewareProvider
}

func SrvInit() *Server {
	// initialising ConfigProvider
	secretConfig := configprovider.NewConfigProvider()

	// reading configs
	secretConfig.Read()

	// initialising PSQLProvider
	db := dbprovider.NewPSQLProvider(os.Getenv("CONNECTION_STRING"), secretConfig.GetInt("MAX_CONNECTION"), secretConfig.GetInt("MAX_IDLE_CONNECTION"))

	// dbHelpProvider contains all db related helper functions
	dbHelper := dbhelpprovider.NewDBHelper(db.DB())

	// initialising in memory cache
	cache := make(map[string]interface{})

	// initialising in cache provider that contains the caching method
	cacheProvider := cacheprovider.NewCacheProvider(cache)

	// newMiddlewareProvider provides middlewares for the routes
	newMiddlewareProvider := middlewareprovider.NewMiddleware(db.DB())

	return &Server{
		PSQL:          db,
		DBHelper:      dbHelper,
		cache:         cache,
		CacheProvider: cacheProvider,
		Middlewares:   newMiddlewareProvider,
	}
}

// Start - starting the server and injecting routes
func (srv *Server) Start() {

	addr := "localhost:8082"
	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           srv.InjectRoutes(),
		ReadTimeout:       defaultServerRequestTimeoutMinutes * time.Minute,
		ReadHeaderTimeout: defaultServerReadHeaderTimeoutSeconds * time.Second,
		WriteTimeout:      defaultServerRequestWriteTimeoutMinutes * time.Minute,
	}
	srv.httpServer = httpSrv

	tokenDetail, err := srv.DBHelper.PopulateCache()
	if err != nil {
		logrus.Errorf("%v", err)
		return
	}

	for i := range tokenDetail {
		err = srv.CacheProvider.Set(tokenDetail[i].Token, tokenDetail[i])
		if err != nil {
			logrus.Errorf("Start: not able to set cache: %v", err)
		}
	}

	logrus.Info("Server running at PORT ", addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
		return
	}
}

// Stop - closing the server
func (srv *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Info("closing server...")
	_ = srv.httpServer.Shutdown(ctx)
	logrus.Info("Done")
}
