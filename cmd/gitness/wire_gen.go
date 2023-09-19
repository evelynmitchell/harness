// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/harness/gitness/cli/server"
	"github.com/harness/gitness/encrypt"
	"github.com/harness/gitness/events"
	"github.com/harness/gitness/gitrpc"
	server3 "github.com/harness/gitness/gitrpc/server"
	"github.com/harness/gitness/gitrpc/server/cron"
	check2 "github.com/harness/gitness/internal/api/controller/check"
	"github.com/harness/gitness/internal/api/controller/connector"
	"github.com/harness/gitness/internal/api/controller/execution"
	"github.com/harness/gitness/internal/api/controller/githook"
	logs2 "github.com/harness/gitness/internal/api/controller/logs"
	"github.com/harness/gitness/internal/api/controller/pipeline"
	"github.com/harness/gitness/internal/api/controller/plugin"
	"github.com/harness/gitness/internal/api/controller/principal"
	pullreq2 "github.com/harness/gitness/internal/api/controller/pullreq"
	"github.com/harness/gitness/internal/api/controller/repo"
	"github.com/harness/gitness/internal/api/controller/secret"
	"github.com/harness/gitness/internal/api/controller/service"
	"github.com/harness/gitness/internal/api/controller/serviceaccount"
	"github.com/harness/gitness/internal/api/controller/space"
	"github.com/harness/gitness/internal/api/controller/system"
	"github.com/harness/gitness/internal/api/controller/template"
	"github.com/harness/gitness/internal/api/controller/trigger"
	"github.com/harness/gitness/internal/api/controller/user"
	webhook2 "github.com/harness/gitness/internal/api/controller/webhook"
	"github.com/harness/gitness/internal/auth/authn"
	"github.com/harness/gitness/internal/auth/authz"
	"github.com/harness/gitness/internal/bootstrap"
	events3 "github.com/harness/gitness/internal/events/git"
	events2 "github.com/harness/gitness/internal/events/pullreq"
	"github.com/harness/gitness/internal/pipeline/canceler"
	"github.com/harness/gitness/internal/pipeline/commit"
	"github.com/harness/gitness/internal/pipeline/file"
	"github.com/harness/gitness/internal/pipeline/manager"
	plugin2 "github.com/harness/gitness/internal/pipeline/plugin"
	"github.com/harness/gitness/internal/pipeline/runner"
	"github.com/harness/gitness/internal/pipeline/scheduler"
	"github.com/harness/gitness/internal/pipeline/triggerer"
	"github.com/harness/gitness/internal/router"
	server2 "github.com/harness/gitness/internal/server"
	"github.com/harness/gitness/internal/services"
	"github.com/harness/gitness/internal/services/codecomments"
	"github.com/harness/gitness/internal/services/exporter"
	"github.com/harness/gitness/internal/services/importer"
	"github.com/harness/gitness/internal/services/job"
	"github.com/harness/gitness/internal/services/metric"
	"github.com/harness/gitness/internal/services/pullreq"
	trigger2 "github.com/harness/gitness/internal/services/trigger"
	"github.com/harness/gitness/internal/services/webhook"
	"github.com/harness/gitness/internal/sse"
	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/internal/store/cache"
	"github.com/harness/gitness/internal/store/database"
	"github.com/harness/gitness/internal/store/logs"
	"github.com/harness/gitness/internal/url"
	"github.com/harness/gitness/livelog"
	"github.com/harness/gitness/lock"
	"github.com/harness/gitness/pubsub"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/check"
)

// Injectors from wire.go:

func initSystem(ctx context.Context, config *types.Config) (*server.System, error) {
	databaseConfig := server.ProvideDatabaseConfig(config)
	db, err := database.ProvideDatabase(ctx, databaseConfig)
	if err != nil {
		return nil, err
	}
	principalUID := check.ProvidePrincipalUIDCheck()
	spacePathTransformation := store.ProvidePathTransformation()
	spacePathStore := database.ProvideSpacePathStore(db, spacePathTransformation)
	spacePathCache := cache.ProvidePathCache(spacePathStore, spacePathTransformation)
	spaceStore := database.ProvideSpaceStore(db, spacePathCache, spacePathStore)
	principalInfoView := database.ProvidePrincipalInfoView(db)
	principalInfoCache := cache.ProvidePrincipalInfoCache(principalInfoView)
	membershipStore := database.ProvideMembershipStore(db, principalInfoCache, spacePathStore)
	permissionCache := authz.ProvidePermissionCache(spaceStore, membershipStore)
	authorizer := authz.ProvideAuthorizer(permissionCache, spaceStore)
	principalUIDTransformation := store.ProvidePrincipalUIDTransformation()
	principalStore := database.ProvidePrincipalStore(db, principalUIDTransformation)
	tokenStore := database.ProvideTokenStore(db)
	controller := user.ProvideController(db, principalUID, authorizer, principalStore, tokenStore, membershipStore)
	serviceController := service.NewController(principalUID, authorizer, principalStore)
	bootstrapBootstrap := bootstrap.ProvideBootstrap(config, controller, serviceController)
	authenticator := authn.ProvideAuthenticator(principalStore, tokenStore)
	provider, err := url.ProvideURLProvider(config)
	if err != nil {
		return nil, err
	}
	pathUID := check.ProvidePathUIDCheck()
	repoStore := database.ProvideRepoStore(db, spacePathCache, spacePathStore)
	pipelineStore := database.ProvidePipelineStore(db)
	gitrpcConfig, err := server.ProvideGitRPCClientConfig()
	if err != nil {
		return nil, err
	}
	gitrpcInterface, err := gitrpc.ProvideClient(gitrpcConfig)
	if err != nil {
		return nil, err
	}
	triggerStore := database.ProvideTriggerStore(db)
	encrypter, err := encrypt.ProvideEncrypter(config)
	if err != nil {
		return nil, err
	}
	jobStore := database.ProvideJobStore(db)
	pubsubConfig := pubsub.ProvideConfig(config)
	universalClient, err := server.ProvideRedis(config)
	if err != nil {
		return nil, err
	}
	pubSub := pubsub.ProvidePubSub(pubsubConfig, universalClient)
	executor := job.ProvideExecutor(jobStore, pubSub)
	lockConfig := server.ProvideLockConfig(config)
	mutexManager := lock.ProvideMutexManager(lockConfig, universalClient)
	jobScheduler, err := job.ProvideScheduler(jobStore, executor, mutexManager, pubSub, config)
	if err != nil {
		return nil, err
	}
	streamer := sse.ProvideEventsStreaming(pubSub)
	repository, err := importer.ProvideRepoImporter(config, provider, gitrpcInterface, db, repoStore, pipelineStore, triggerStore, encrypter, jobScheduler, executor, streamer)
	if err != nil {
		return nil, err
	}
	repoController := repo.ProvideController(config, db, provider, pathUID, authorizer, repoStore, spaceStore, pipelineStore, principalStore, gitrpcInterface, repository)
	executionStore := database.ProvideExecutionStore(db)
	stageStore := database.ProvideStageStore(db)
	schedulerScheduler, err := scheduler.ProvideScheduler(stageStore, mutexManager)
	if err != nil {
		return nil, err
	}
	stepStore := database.ProvideStepStore(db)
	cancelerCanceler := canceler.ProvideCanceler(executionStore, streamer, repoStore, schedulerScheduler, stageStore, stepStore)
	commitService := commit.ProvideCommitService(gitrpcInterface)
	checkStore := database.ProvideCheckStore(db, principalInfoCache)
	fileService := file.ProvideFileService(gitrpcInterface)
	triggererTriggerer := triggerer.ProvideTriggerer(executionStore, checkStore, stageStore, db, pipelineStore, fileService, schedulerScheduler, repoStore)
	executionController := execution.ProvideController(db, authorizer, executionStore, cancelerCanceler, commitService, triggererTriggerer, repoStore, stageStore, pipelineStore)
	logStore := logs.ProvideLogStore(db, config)
	logStream := livelog.ProvideLogStream(config)
	logsController := logs2.ProvideController(db, authorizer, executionStore, repoStore, pipelineStore, stageStore, stepStore, logStore, logStream)
	secretStore := database.ProvideSecretStore(db)
	connectorStore := database.ProvideConnectorStore(db)
	templateStore := database.ProvideTemplateStore(db)
	exporterRepository, err := exporter.ProvideSpaceExporter(provider, gitrpcInterface, repoStore, jobScheduler, executor, encrypter, streamer)
	if err != nil {
		return nil, err
	}
	spaceController := space.ProvideController(config, db, provider, streamer, pathUID, authorizer, spacePathStore, pipelineStore, secretStore, connectorStore, templateStore, spaceStore, repoStore, principalStore, repoController, membershipStore, repository, exporterRepository)
	pipelineController := pipeline.ProvideController(db, pathUID, repoStore, triggerStore, authorizer, pipelineStore)
	secretController := secret.ProvideController(db, pathUID, encrypter, secretStore, authorizer, spaceStore)
	triggerController := trigger.ProvideController(db, authorizer, triggerStore, pathUID, pipelineStore, repoStore)
	connectorController := connector.ProvideController(db, pathUID, connectorStore, authorizer, spaceStore)
	templateController := template.ProvideController(db, pathUID, templateStore, authorizer, spaceStore)
	pluginStore := database.ProvidePluginStore(db)
	pluginController := plugin.ProvideController(db, pluginStore)
	pullReqStore := database.ProvidePullReqStore(db, principalInfoCache)
	pullReqActivityStore := database.ProvidePullReqActivityStore(db, principalInfoCache)
	codeCommentView := database.ProvideCodeCommentView(db)
	pullReqReviewStore := database.ProvidePullReqReviewStore(db)
	pullReqReviewerStore := database.ProvidePullReqReviewerStore(db, principalInfoCache)
	pullReqFileViewStore := database.ProvidePullReqFileViewStore(db)
	eventsConfig, err := server.ProvideEventsConfig()
	if err != nil {
		return nil, err
	}
	eventsSystem, err := events.ProvideSystem(eventsConfig, universalClient)
	if err != nil {
		return nil, err
	}
	reporter, err := events2.ProvideReporter(eventsSystem)
	if err != nil {
		return nil, err
	}
	migrator := codecomments.ProvideMigrator(gitrpcInterface)
	readerFactory, err := events3.ProvideReaderFactory(eventsSystem)
	if err != nil {
		return nil, err
	}
	eventsReaderFactory, err := events2.ProvideReaderFactory(eventsSystem)
	if err != nil {
		return nil, err
	}
	repoGitInfoView := database.ProvideRepoGitInfoView(db)
	repoGitInfoCache := cache.ProvideRepoGitInfoCache(repoGitInfoView)
	pullreqService, err := pullreq.ProvideService(ctx, config, readerFactory, eventsReaderFactory, reporter, gitrpcInterface, db, repoGitInfoCache, repoStore, pullReqStore, pullReqActivityStore, codeCommentView, migrator, pullReqFileViewStore, pubSub, provider)
	if err != nil {
		return nil, err
	}
	pullreqController := pullreq2.ProvideController(db, provider, authorizer, pullReqStore, pullReqActivityStore, codeCommentView, pullReqReviewStore, pullReqReviewerStore, repoStore, principalStore, pullReqFileViewStore, gitrpcInterface, reporter, mutexManager, migrator, pullreqService)
	webhookConfig := server.ProvideWebhookConfig(config)
	webhookStore := database.ProvideWebhookStore(db)
	webhookExecutionStore := database.ProvideWebhookExecutionStore(db)
	webhookService, err := webhook.ProvideService(ctx, webhookConfig, readerFactory, eventsReaderFactory, webhookStore, webhookExecutionStore, repoStore, pullReqStore, provider, principalStore, gitrpcInterface, encrypter)
	if err != nil {
		return nil, err
	}
	webhookController := webhook2.ProvideController(webhookConfig, db, authorizer, webhookStore, webhookExecutionStore, repoStore, webhookService, encrypter)
	eventsReporter, err := events3.ProvideReporter(eventsSystem)
	if err != nil {
		return nil, err
	}
	githookController := githook.ProvideController(db, authorizer, principalStore, repoStore, eventsReporter)
	serviceaccountController := serviceaccount.NewController(principalUID, authorizer, principalStore, spaceStore, repoStore, tokenStore)
	principalController := principal.ProvideController(principalStore)
	checkController := check2.ProvideController(db, authorizer, repoStore, checkStore, gitrpcInterface)
	systemController := system.NewController(principalStore, config)
	apiHandler := router.ProvideAPIHandler(config, authenticator, repoController, executionController, logsController, spaceController, pipelineController, secretController, triggerController, connectorController, templateController, pluginController, pullreqController, webhookController, githookController, serviceaccountController, controller, principalController, checkController, systemController)
	gitHandler := router.ProvideGitHandler(config, provider, repoStore, authenticator, authorizer, gitrpcInterface)
	webHandler := router.ProvideWebHandler(config)
	routerRouter := router.ProvideRouter(config, apiHandler, gitHandler, webHandler)
	serverServer := server2.ProvideServer(config, routerRouter)
	executionManager := manager.ProvideExecutionManager(config, executionStore, pipelineStore, provider, streamer, fileService, logStore, logStream, checkStore, repoStore, schedulerScheduler, secretStore, stageStore, stepStore, principalStore)
	client := manager.ProvideExecutionClient(executionManager, config)
	pluginManager := plugin2.ProvidePluginManager(config, pluginStore)
	runtimeRunner, err := runner.ProvideExecutionRunner(config, client, pluginManager, executionManager)
	if err != nil {
		return nil, err
	}
	poller := runner.ProvideExecutionPoller(runtimeRunner, config, client)
	serverConfig, err := server.ProvideGitRPCServerConfig()
	if err != nil {
		return nil, err
	}
	goGitRepoProvider := server3.ProvideGoGitRepoProvider()
	cacheCache := server3.ProvideLastCommitCache(serverConfig, universalClient, goGitRepoProvider)
	gitAdapter, err := server3.ProvideGITAdapter(goGitRepoProvider, cacheCache)
	if err != nil {
		return nil, err
	}
	grpcServer, err := server3.ProvideServer(serverConfig, gitAdapter)
	if err != nil {
		return nil, err
	}
	cronManager := cron.ProvideManager(serverConfig)
	triggerConfig := server.ProvideTriggerConfig(config)
	triggerService, err := trigger2.ProvideService(ctx, triggerConfig, triggerStore, commitService, pullReqStore, repoStore, pipelineStore, triggererTriggerer, readerFactory, eventsReaderFactory)
	if err != nil {
		return nil, err
	}
	collector, err := metric.ProvideCollector(config, principalStore, repoStore, pipelineStore, executionStore, jobScheduler, executor)
	if err != nil {
		return nil, err
	}
	servicesServices := services.ProvideServices(webhookService, pullreqService, triggerService, jobScheduler, collector)
	serverSystem := server.NewSystem(bootstrapBootstrap, serverServer, poller, grpcServer, pluginManager, cronManager, servicesServices)
	return serverSystem, nil
}
