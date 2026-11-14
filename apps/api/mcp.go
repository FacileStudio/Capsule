package main

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/FacileStudio/Capsule/apps/api/internal/env"
	"github.com/FacileStudio/Capsule/apps/api/modules/auth"
	"github.com/FacileStudio/Capsule/apps/api/modules/mcp"
	"github.com/FacileStudio/porte/mcpauth"
	mcpkit "github.com/FacileStudio/tronc/mcp"
)

// mcpSetup is what the MCP server needs besides the services buildRouter
// builds. A nil setup leaves /mcp unmounted.
type mcpSetup struct {
	config  mcp.Config
	pending mcpkit.PostgresPending
}

// newMCPSetup prepares the MCP server when MCP_RESOURCE is set and SSO is
// on, and nothing otherwise. The verifier runs Registre discovery here, so a
// wrong issuer fails at boot like the OIDC kit does. A person's first MCP
// call creates their account through the OIDC kit, as a first web login
// would.
func newMCPSetup(ctx context.Context, appEnv env.Config, sqlDB *sql.DB, kits *auth.Kits, appLogger *slog.Logger) (*mcpSetup, error) {
	if appEnv.OIDC == nil || appEnv.OIDC.MCPResource == "" {
		return nil, nil
	}
	verifier, err := mcpauth.Verifier(ctx, mcpauth.Config{
		Issuer:      appEnv.OIDC.Issuer,
		Resource:    appEnv.OIDC.MCPResource,
		ScopePrefix: "capsule",
		Identities:  kits.Identities,
		Provision:   true,
		Accounts:    kits.OIDC,
		Logger:      appLogger,
	})
	if err != nil {
		return nil, err
	}
	pending := mcpkit.PostgresPending{DB: sqlDB}
	appLogger.Info("mcp server on", slog.String("resource", appEnv.OIDC.MCPResource))
	return &mcpSetup{config: mcp.Config{
		Resource: appEnv.OIDC.MCPResource,
		Issuer:   appEnv.OIDC.Issuer,
		Version:  version,
		Origins:  appEnv.CORSAllowedOrigins,
		StateKey: mcpkit.StateKey("capsule", appEnv.OIDC.ClientSecret),
		Verifier: verifier,
		Pending:  pending,
		Logger:   appLogger,
	}, pending: pending}, nil
}

// keepPruned prunes resolved and expired pending actions until ctx ends. run
// starts it once startup succeeded, so no prune outlives the database pool.
func (s *mcpSetup) keepPruned(ctx context.Context, logger *slog.Logger) {
	if s == nil {
		return
	}
	go mcpkit.KeepPruned(ctx, s.pending, 0, 0, logger)
}

// service builds the MCP server over deps, or nothing for a nil setup.
func (s *mcpSetup) service(deps mcp.Deps) (*mcp.Service, error) {
	if s == nil {
		return nil, nil
	}
	return mcp.NewService(deps, s.config)
}
