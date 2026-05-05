// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package routers

import (
	"context"
	"net/http"
	"reflect"
	"runtime"

	"code.aegit.io/aegit/models"
	authmodel "code.aegit.io/aegit/models/auth"
	"code.aegit.io/aegit/modules/cache"
	"code.aegit.io/aegit/modules/eventsource"
	"code.aegit.io/aegit/modules/git"
	"code.aegit.io/aegit/modules/git/gitcmd"
	"code.aegit.io/aegit/modules/log"
	"code.aegit.io/aegit/modules/markup"
	"code.aegit.io/aegit/modules/markup/external"
	"code.aegit.io/aegit/modules/setting"
	"code.aegit.io/aegit/modules/ssh"
	"code.aegit.io/aegit/modules/storage"
	"code.aegit.io/aegit/modules/svg"
	"code.aegit.io/aegit/modules/system"
	"code.aegit.io/aegit/modules/translation"
	"code.aegit.io/aegit/modules/util"
	"code.aegit.io/aegit/modules/web"
	"code.aegit.io/aegit/modules/web/routing"
	actions_router "code.aegit.io/aegit/routers/api/actions"
	packages_router "code.aegit.io/aegit/routers/api/packages"
	apiv1 "code.aegit.io/aegit/routers/api/v1"
	"code.aegit.io/aegit/routers/common"
	"code.aegit.io/aegit/routers/private"
	web_routers "code.aegit.io/aegit/routers/web"
	actions_service "code.aegit.io/aegit/services/actions"
	asymkey_service "code.aegit.io/aegit/services/asymkey"
	"code.aegit.io/aegit/services/auth"
	"code.aegit.io/aegit/services/auth/source/oauth2"
	"code.aegit.io/aegit/services/automerge"
	"code.aegit.io/aegit/services/cron"
	feed_service "code.aegit.io/aegit/services/feed"
	indexer_service "code.aegit.io/aegit/services/indexer"
	"code.aegit.io/aegit/services/mailer"
	mailer_incoming "code.aegit.io/aegit/services/mailer/incoming"
	markup_service "code.aegit.io/aegit/services/markup"
	repo_migrations "code.aegit.io/aegit/services/migrations"
	mirror_service "code.aegit.io/aegit/services/mirror"
	"code.aegit.io/aegit/services/oauth2_provider"
	packages_spec "code.aegit.io/aegit/services/packages/pkgspec"
	pull_service "code.aegit.io/aegit/services/pull"
	release_service "code.aegit.io/aegit/services/release"
	repo_service "code.aegit.io/aegit/services/repository"
	"code.aegit.io/aegit/services/repository/archiver"
	"code.aegit.io/aegit/services/task"
	"code.aegit.io/aegit/services/uinotification"
	"code.aegit.io/aegit/services/webhook"
)

func mustInit(fn func() error) {
	err := fn()
	if err != nil {
		ptr := reflect.ValueOf(fn).Pointer()
		fi := runtime.FuncForPC(ptr)
		log.Fatal("%s failed: %v", fi.Name(), err)
	}
}

func mustInitCtx(ctx context.Context, fn func(ctx context.Context) error) {
	err := fn(ctx)
	if err != nil {
		ptr := reflect.ValueOf(fn).Pointer()
		fi := runtime.FuncForPC(ptr)
		log.Fatal("%s(ctx) failed: %v", fi.Name(), err)
	}
}

func syncAppConfForGit(ctx context.Context) error {
	runtimeState := new(system.RuntimeState)
	if err := system.AppState.Get(ctx, runtimeState); err != nil {
		return err
	}

	updated := false
	if runtimeState.LastAppPath != setting.AppPath {
		log.Info("AppPath changed from '%s' to '%s'", runtimeState.LastAppPath, setting.AppPath)
		runtimeState.LastAppPath = setting.AppPath
		updated = true
	}
	if runtimeState.LastCustomConf != setting.CustomConf {
		log.Info("CustomConf changed from '%s' to '%s'", runtimeState.LastCustomConf, setting.CustomConf)
		runtimeState.LastCustomConf = setting.CustomConf
		updated = true
	}

	if updated {
		log.Info("re-sync repository hooks ...")
		mustInitCtx(ctx, repo_service.SyncRepositoryHooks)

		log.Info("re-write ssh public keys ...")
		mustInitCtx(ctx, asymkey_service.RewriteAllPublicKeys)

		return system.AppState.Set(ctx, runtimeState)
	}
	return nil
}

func InitWebInstallPage(ctx context.Context) {
	translation.InitLocales(ctx)
	setting.LoadSettingsForInstall()
	mustInit(svg.Init)
}

// InitWebInstalled is for the global configuration of an installed instance
func InitWebInstalled(ctx context.Context) {
	mustInit(git.InitFull)
	log.Info("Git version: %s (home: %s)", git.DefaultFeatures().VersionInfo(), gitcmd.HomeDir())
	if !git.DefaultFeatures().SupportHashSha256 {
		log.Warn("sha256 hash support is disabled - requires Git >= 2.42." + util.Iif(git.DefaultFeatures().UsingGogit, " Gogit is currently unsupported.", ""))
	}

	// Setup i18n
	translation.InitLocales(ctx)

	setting.LoadSettings()
	mustInit(storage.Init)

	mailer.NewContext(ctx)
	mustInit(cache.Init)
	mustInit(feed_service.Init)
	mustInit(uinotification.Init)
	mustInitCtx(ctx, archiver.Init)

	external.RegisterRenderers()
	markup.Init(markup_service.FormalRenderHelperFuncs())

	mustInitCtx(ctx, common.InitDBEngine)
	log.Info("ORM engine initialization successful!")
	mustInit(system.Init)
	mustInitCtx(ctx, oauth2.Init)
	mustInitCtx(ctx, oauth2_provider.Init)
	mustInit(release_service.Init)

	mustInitCtx(ctx, models.Init)
	mustInitCtx(ctx, authmodel.Init)
	mustInitCtx(ctx, repo_service.Init)
	mustInit(packages_spec.InitManager)

	// Booting long running goroutines.
	mustInit(indexer_service.Init)

	mirror_service.InitSyncMirrors()
	mustInit(webhook.Init)
	mustInit(pull_service.Init)
	mustInit(automerge.Init)
	mustInit(task.Init)
	mustInit(repo_migrations.Init)
	eventsource.GetManager().Init()
	mustInitCtx(ctx, mailer_incoming.Init)

	mustInitCtx(ctx, syncAppConfForGit)

	mustInit(ssh.Init)

	auth.Init()
	mustInit(svg.Init)

	mustInitCtx(ctx, actions_service.Init)

	mustInit(repo_service.InitLicenseClassifier)

	// Finally start up the cron
	cron.Init(ctx)
}

// NormalRoutes represents non install routes
func NormalRoutes() *web.Router {
	r := web.NewRouter()
	r.BeforeRouting(common.ProtocolMiddlewares()...)

	r.AfterRouting(common.MaintenanceModeHandler())

	r.Mount("/", web_routers.Routes())
	r.Mount("/api/v1", apiv1.Routes())
	r.Mount("/api/internal", private.Routes())

	r.Post("/-/fetch-redirect", common.FetchRedirectDelegate)

	if setting.Packages.Enabled {
		// This implements package support for most package managers
		r.Mount("/api/packages", packages_router.CommonRoutes())
		// This implements the OCI API, this container registry "/v2" endpoint must be in the root of the site.
		// If site admin deploys Gitea in a sub-path, they must configure their reverse proxy to map the "https://host/v2" endpoint to Gitea.
		r.Mount("/v2", packages_router.ContainerRoutes())
	}

	if setting.Actions.Enabled {
		prefix := "/api/actions"
		r.Mount(prefix, actions_router.Routes(prefix))

		// TODO: Pipeline api used for runner internal communication with gitea server. but only artifact is used for now.
		// In Github, it uses ACTIONS_RUNTIME_URL=https://pipelines.actions.githubusercontent.com/fLgcSHkPGySXeIFrg8W8OBSfeg3b5Fls1A1CwX566g8PayEGlg/
		// TODO: this prefix should be generated with a token string with runner ?
		prefix = "/api/actions_pipeline"
		r.Mount(prefix, actions_router.ArtifactsRoutes(prefix))
		prefix = actions_router.ArtifactV4RouteBase
		r.Mount(prefix, actions_router.ArtifactsV4Routes(prefix))
	}

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		defer routing.RecordFuncInfo(req.Context(), routing.GetFuncInfo(http.NotFound, "GlobalNotFound"))()
		http.NotFound(w, req)
	})
	return r
}
