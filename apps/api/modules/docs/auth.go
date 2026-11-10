package docs

import (
	"net/http"

	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
)

type meResponse struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type sessionUser struct {
	Email     string `json:"email"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type sessionResponse struct {
	User *sessionUser `json:"user"`
}

func authModule() Module {
	const session = "session cookie or bearer token required"
	return Module{
		Name:        "auth",
		Description: "Sign-in, served by the shared porte kit. Pastes never require it; the CLI holds a session through it.",
		Routes: []Route{
			{
				Method:       "GET",
				Path:         porte.RouteConfig,
				Summary:      "Describe the sign-in modes",
				Description:  "Returns {sso_only, oidc_enabled} so a client knows which sign-in paths exist.",
				Auth:         "public",
				ResponseBody: porte.ConfigResponse{},
			},
			{
				Method:      "GET",
				Path:        porte.RouteLogin,
				Summary:     "Start the single sign-on flow",
				Description: "Redirects to the identity provider. A CLI adds ?flow=cli&port=N to be redirected back to a loopback port with a one-time code.",
				Auth:        "public",
				Status:      http.StatusFound,
			},
			{
				Method:      "GET",
				Path:        porte.RouteCallback,
				Summary:     "Complete the single sign-on flow",
				Description: "Validates the provider response, opens a session and redirects to the app, or to the CLI loopback port.",
				Auth:        "public",
				Status:      http.StatusFound,
			},
			{
				Method:       "POST",
				Path:         porte.RouteExchange,
				Summary:      "Trade a CLI login code for a token",
				Description:  "Consumes the single-use code from the CLI flow and returns a bearer token.",
				Auth:         "public",
				RequestBody:  porte.ExchangeRequest{},
				ResponseBody: porte.ExchangeResponse{},
				Errors: []Error{
					{Status: 400, Code: "invalid_argument", Description: "The body carried no code."},
					{Status: 401, Code: "unauthenticated", Description: "The code is unknown, used or expired."},
				},
			},
			{
				Method:       "POST",
				Path:         porte.RouteLogout,
				Summary:      "Sign out",
				Auth:         session,
				ResponseBody: porte.LogoutResponse{},
			},
			{
				Method:       "GET",
				Path:         "/auth/me",
				Summary:      "Get the signed-in user",
				Auth:         session,
				ResponseBody: meResponse{},
				Errors:       []Error{{Status: 401, Code: "unauthenticated", Description: "No valid session."}},
			},
			{
				Method:       "GET",
				Path:         "/auth/session",
				Summary:      "Get the signed-in user, or null",
				Description:  "Answers 200 with {\"user\": null} when there is no session, so a browser can tell whether to show the user menu without a failed request.",
				Auth:         "public",
				ResponseBody: sessionResponse{},
			},
			{
				Method:      "POST",
				Path:        porte.RouteSyncProfile,
				Summary:     "Refresh the profile from the identity provider",
				Auth:        session,
				Description: "Mounted only when single sign-on is configured.",
			},
			{
				Method:      "POST",
				Path:        porte.RouteBackchannelLogout,
				Summary:     "Accept a back-channel logout",
				Description: "Called by the identity provider to end the sessions of a user signed out elsewhere.",
				Auth:        "signed logout token from the identity provider",
			},
			invitationsRoute(),
		},
	}
}

func invitationsRoute() Route {
	route := kit.InvitationsRoute()
	route.Description += " Mounted only when single sign-on is configured."
	route.ResponseBody = kit.InvitationsResponse{}
	return route
}
