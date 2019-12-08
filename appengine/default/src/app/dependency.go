package app

import (
	"firebase.google.com/go/messaging"

	"github.com/rabee-inc/push/appengine/default/src/config"
	"github.com/rabee-inc/push/appengine/default/src/handler/api"
	"github.com/rabee-inc/push/appengine/default/src/handler/worker"
	"github.com/rabee-inc/push/appengine/default/src/lib/cloudfirestore"
	"github.com/rabee-inc/push/appengine/default/src/lib/cloudtasks"
	"github.com/rabee-inc/push/appengine/default/src/lib/deploy"
	"github.com/rabee-inc/push/appengine/default/src/lib/internalauth"
	"github.com/rabee-inc/push/appengine/default/src/lib/jsonrpc2"
	"github.com/rabee-inc/push/appengine/default/src/lib/log"
	"github.com/rabee-inc/push/appengine/default/src/repository"
	"github.com/rabee-inc/push/appengine/default/src/service"
)

// Dependency ... 依存性
type Dependency struct {
	Log                  *log.Middleware
	InternalAuth         *internalauth.Middleware
	JSONRPC2Handler      *jsonrpc2.Handler
	EntryAction          *api.EntryAction
	LeaveAction          *api.LeaveAction
	SendByUsersAction    *api.SendByUsersAction
	SendByAllUsersAction *api.SendByAllUsersAction
	GetReserve           *api.GetReserveAction
	ListReserve          *api.ListReserveAction
	CreateReserve        *api.CreateReserveAction
	UpdateReserve        *api.UpdateReserveAction
	SendHandler          *worker.SendHandler
}

// Inject ... 依存性を注入する
func (d *Dependency) Inject(e *Environment) {
	// FCM
	var env string
	appIDs := config.GetFCMAppIDs()
	fcmClients := map[string]*messaging.Client{}
	fcmServerKeys := map[string]string{}
	if deploy.IsProduction() {
		env = "production"
	} else {
		env = "staging"
	}
	for _, appID := range appIDs {
		fcmEnv := config.GetFCMEnv(env, appID)
		fcmClients[appID] = config.GetClient(env, appID)
		fcmServerKeys[appID] = fcmEnv.ServerKey
	}

	// Client
	tCli := cloudtasks.NewClient(e.Port, e.Deploy, e.ProjectID, "default", e.LocationID, e.InternalAuthToken)
	var lCli log.Writer
	if deploy.IsLocal() {
		lCli = log.NewWriterStdout()
	} else {
		lCli = log.NewWriterStackdriver(e.ProjectID)
	}
	fCli := cloudfirestore.NewClient(e.ProjectID)

	// Repository
	tRepo := repository.NewToken(fCli)
	fRepo := repository.NewFcm(fcmClients, fcmServerKeys)
	rRepo := repository.NewReserve(fCli)

	// Service
	rgSvc := service.NewRegister(tRepo, fRepo)
	sSvc := service.NewSender(tRepo, fRepo, rRepo, tCli, fCli)
	rSvc := service.NewReserve(rRepo, fCli)

	// Middleware
	d.Log = log.NewMiddleware(lCli, e.MinLogSeverity)
	d.InternalAuth = internalauth.NewMiddleware(e.InternalAuthToken)

	// Handler
	d.JSONRPC2Handler = jsonrpc2.NewHandler()
	d.SendHandler = worker.NewSendHandler(sSvc, rSvc)

	// Action
	d.EntryAction = api.NewEntryAction(rgSvc)
	d.LeaveAction = api.NewLeaveAction(rgSvc)
	d.SendByUsersAction = api.NewSendByUsersAction(tCli)
	d.SendByAllUsersAction = api.NewSendByAllUsersAction(sSvc)
	d.GetReserve = api.NewGetReserveAction(rSvc)
	d.ListReserve = api.NewListReserveAction(rSvc)
	d.CreateReserve = api.NewCreateReserveAction(rSvc)
	d.UpdateReserve = api.NewUpdateReserveAction(rSvc)
}
