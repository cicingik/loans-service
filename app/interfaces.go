package app

import "github.com/cicingik/loans-service/repository/postgre"

type (
	// IWebApplication ...
	IWebApplication interface {
		// ManageDb
		ManageHTTPServer
		Start() error
		Stop() error
	}

	// ManageDb ...
	ManageDb interface {
		SetDb(*postgre.DbEngine)
		GetDb() *postgre.DbEngine
	}

	// DeliveryHTTPEngine ...
	DeliveryHTTPEngine interface {
		Serve() error
	}

	// ManageHTTPServer ...
	ManageHTTPServer interface {
		SetHttpServer(*DeliveryHTTPEngine)
		GetHttpServer() *DeliveryHTTPEngine
	}
)
