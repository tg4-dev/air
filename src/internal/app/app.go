package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/tg4-dev/air/src/internal/discovery"
	"github.com/tg4-dev/air/src/internal/engines"
	"github.com/tg4-dev/air/src/internal/nodeinfo"
	"github.com/tg4-dev/air/src/internal/utils"
)

const shutdownTimeout = 5 * time.Second
const browseTimeout = 10 * time.Second

type App struct {
	logger     *slog.Logger
	node       *discovery.Node
	infoServer *nodeinfo.Server
	engine     engines.Engine
	browser    *discovery.Browser
	advertiser *discovery.Advertiser
}

func New() (*App, error) {
	logger := utils.NewLogger()

	node, err := discovery.NewNode()
	if err != nil {
		return nil, fmt.Errorf("create node: %w", err)
	}

	statusProvider := nodeinfo.NewStatusProvider(node.ID(), node.Hostname(), time.Now())

	// TODO should be filled with config json
	workingEngineName := "ollama"

	return &App{
		logger:     logger,
		node:       node,
		infoServer: nodeinfo.NewServer(statusProvider),
		engine:     engines.WorkingEngine(workingEngineName),
		browser:    &discovery.Browser{},
	}, nil
}

// Run starts all components and blocks until ctx is cancelled, then shuts them down.
func (a *App) Run(ctx context.Context) error {
	addr := ":" + strconv.Itoa(a.node.Port())
	if err := a.infoServer.Start(addr); err != nil {
		return fmt.Errorf("start nodeinfo HTTP server on %s: %w", addr, err)
	}
	a.logger.Info("nodeinfo HTTP server listening", "addr", addr)
	defer a.shutdown()

	a.engine.TestRequest()

	advertiser, err := discovery.NewAdvertiser(*a.node)
	if err != nil {
		return fmt.Errorf("start mDNS advertiser: %w", err)
	}
	a.advertiser = advertiser

	browseCtx, cancel := context.WithTimeout(ctx, browseTimeout)
	defer cancel()
	a.browser.Update(browseCtx)
	for _, peer := range a.browser.GetPeers() {
		a.logger.Info("peer found", "peer", fmt.Sprintf("%+v", peer))
	}

	<-ctx.Done()
	a.logger.Info("shutting down")
	return nil
}

func (a *App) shutdown() {
	if a.advertiser != nil {
		a.advertiser.Shutdown()
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := a.infoServer.Shutdown(ctx); err != nil {
		a.logger.Error("nodeinfo HTTP shutdown", "err", err)
	}
}
