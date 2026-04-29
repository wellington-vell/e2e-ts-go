package auth

import (
	"time"

	"server/internal/lib"

	"github.com/Authula/authula"

	authulaconfig "github.com/Authula/authula/config"
	authulamodels "github.com/Authula/authula/models"
	accesscontrolplugin "github.com/Authula/authula/plugins/access-control"
	accesscontrolplugintypes "github.com/Authula/authula/plugins/access-control/types"
	adminplugin "github.com/Authula/authula/plugins/admin"
	adminplugintypes "github.com/Authula/authula/plugins/admin/types"
	csrfplugin "github.com/Authula/authula/plugins/csrf"
	emailplugin "github.com/Authula/authula/plugins/email"
	emailpasswordplugin "github.com/Authula/authula/plugins/email-password"
	emailpasswordtypes "github.com/Authula/authula/plugins/email-password/types"
	emailplugintypes "github.com/Authula/authula/plugins/email/types"
	sessionplugin "github.com/Authula/authula/plugins/session"
	// ratelimitplugin "github.com/Authula/authula/plugins/rate-limit"
	// organizationsplugin "github.com/authula/authula/plugins/organizations"
	// organizationsplugintypes "github.com/authula/authula/plugins/organizations/types"
)

func NewAuthula() (*authula.Auth, error) {
	serverUrl := lib.Env.ServerURL
	secret := lib.Env.AuthSecret
	dbURL := lib.Env.DatabaseURL
	cors := lib.Env.CorsOrigin

	config := authulaconfig.NewConfig(
		authulaconfig.WithAppName("e2e-ts-go"),
		authulaconfig.WithBaseURL(serverUrl),
		authulaconfig.WithBasePath("/auth"),
		authulaconfig.WithSecret(secret),
		authulaconfig.WithDatabase(authulamodels.DatabaseConfig{
			Provider: "postgres",
			URL:      dbURL,
		}),
		authulaconfig.WithSession(authulamodels.SessionConfig{
			CookieName:         "session_token",
			ExpiresIn:          24 * time.Hour,
			UpdateAge:          5 * time.Minute,
			CookieMaxAge:       24 * time.Hour,
			Secure:             false,
			HttpOnly:           true,
			SameSite:           "lax",
			AutoCleanup:        true,
			CleanupInterval:    1 * time.Minute,
			MaxSessionsPerUser: 5,
		}),
		authulaconfig.WithSecurity(authulamodels.SecurityConfig{
			TrustedOrigins: []string{cors},
			CORS: authulamodels.CORSConfig{
				AllowCredentials: true,
				AllowedOrigins:   []string{cors},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Content-Type", "Set-Cookie", "Cookie"},
				MaxAge:           24 * time.Hour,
			},
		}),
		// authulaconfig.WithLogger(authulamodels.LoggerConfig{
		// 	Level: "debug",
		// }),
		authulaconfig.WithRouteMappings(
			[]authulamodels.RouteMapping{
				// Core routes
				{
					Paths:   []string{"GET:/me", "POST:/sign-out"},
					Plugins: []string{"session.auth"},
				},
				// // Email/Password Plugin
				{
					Paths: []string{
						"POST:/email-password/sign-up",
						"POST:/email-password/sign-in",
						"POST:/email-password/change-password",
					},
				},
				{
					Paths: []string{
						"GET:/email-password/verify-email",
						"POST:/email-password/send-email-verification",
						"POST:/email-password/request-password-reset",
						"POST:/email-password/request-email-change",
					},
					Plugins: []string{"session.auth"},
				},

				// // ADMIN Plugin
				// User Management
				{
					Paths: []string{
						"POST:/admin/users",
						"GET:/admin/users",
						"GET:/admin/users/{user_id}",
						"PATCH:/admin/users/{user_id}",
						"DELETE:/admin/users/{user_id}",
					},
					Plugins: []string{"session.auth"},
				},
				// Account Management
				{
					Paths: []string{
						"POST:/admin/users/{user_id}/accounts",
						"GET:/admin/users/{user_id}/accounts",
						"GET:/admin/accounts/{id}",
						"PATCH:/admin/accounts/{id}",
						"DELETE:/admin/accounts/{id}",
					},
					Plugins: []string{"session.auth"},
				},
				// User State Management
				{
					Paths: []string{
						"GET:/admin/users/{user_id}/state",
						"POST:/admin/users/{user_id}/state",
						"PATCH:/admin/users/{user_id}/state",
						"DELETE:/admin/users/{user_id}/state",
						"GET:/admin/users/states/banned",
						"POST:/admin/users/{user_id}/ban",
						"POST:/admin/users/{user_id}/unban",
					},
					Plugins: []string{"session.auth"},
				},
				// Session State Management
				{
					Paths: []string{
						"GET:/admin/sessions/{session_id}/state",
						"POST:/admin/sessions/{session_id}/state",
						"PATCH:/admin/sessions/{session_id}/state",
						"DELETE:/admin/sessions/{session_id}/state",
						"POST:/admin/sessions/{session_id}/revoke",
						"GET:/admin/sessions/states/revoked",
						"GET:/admin/users/{user_id}/sessions",
					},
					Plugins: []string{"session.auth"},
				},
				// Impersonation Management
				{
					Paths: []string{
						"GET:/admin/impersonations",
						"GET:/admin/impersonations/{impersonation_id}",
						"POST:/admin/impersonations",
						"POST:/admin/impersonations/{impersonation_id}/stop",
					},
					Plugins: []string{"session.auth"},
				},

				// // ACCESS CONTROL Plugin
				// Roles Management
				{
					Paths: []string{
						"POST:/access-control/roles",
						"GET:/access-control/roles",
						"GET:/access-control/roles/by-name/{role_name}",
						"GET:/access-control/roles/{role_id}",
						"PATCH:/access-control/roles/{role_id}",
						"DELETE:/access-control/roles/{role_id}",
					},
					Plugins: []string{"session.auth"},
				},
				// Permissions Management
				{
					Paths: []string{
						"POST:/access-control/permissions",
						"GET:/access-control/permissions",
						"GET:/access-control/permissions/{permission_id}",
						"PATCH:/access-control/permissions/{permission_id}",
						"DELETE:/access-control/permissions/{permission_id}",
					},
					Plugins: []string{"session.auth"},
				},
				// Role-Permission Mapping
				{
					Paths: []string{
						"POST:/access-control/roles/{role_id}/permissions",
						"GET:/access-control/roles/{role_id}/permissions",
						"PUT:/access-control/roles/{role_id}/permissions",
						"DELETE:/access-control/roles/{role_id}/permissions/{permission_id}",
					},
					Plugins: []string{"session.auth"},
				},
				// User Access & Permissions
				{
					Paths: []string{
						"GET:/access-control/users/{user_id}/roles",
						"POST:/access-control/users/{user_id}/roles",
						"PUT:/access-control/users/{user_id}/roles",
						"DELETE:/access-control/users/{user_id}/roles/{role_id}",
						"GET:/access-control/users/{user_id}/permissions",
						"POST:/access-control/users/{user_id}/permissions/check",
					},
					Plugins: []string{"session.auth"},
				},

				// // Organizations Plugin
				// // Organizations
				// {
				// 	Paths: []string{
				// 		"POST:/organizations",
				// 		"GET:/organizations",
				// 		"GET:/organizations/{organization_id}",
				// 		"PATCH:/organizations/{organization_id}",
				// 		"DELETE:/organizations/{organization_id}",
				// 	},
				// 	Plugins: []string{"session.auth"},
				// },
				// // Invitations
				// {
				// 	Paths: []string{
				// 		"POST:/organizations/{organization_id}/invitations",
				// 		"GET:/organizations/{organization_id}/invitations",
				// 		"GET:/organizations/{organization_id}/invitations/{invitation_id}",
				// 		"PATCH:/organizations/{organization_id}/invitations/{invitation_id}",
				// 		"POST:/organizations/{organization_id}/invitations/{invitation_id}/accept",
				// 		"POST:/organizations/{organization_id}/invitations/{invitation_id}/reject",
				// 	},
				// 	Plugins: []string{"session.auth"},
				// },
				// // Members
				// {
				// 	Paths: []string{
				// 		"POST:/organizations/{organization_id}/members",
				// 		"GET:/organizations/{organization_id}/members",
				// 		"GET:/organizations/{organization_id}/members/{member_id}",
				// 		"PATCH:/organizations/{organization_id}/members/{member_id}",
				// 		"DELETE:/organizations/{organization_id}/members/{member_id}",
				// 	},
				// 	Plugins: []string{"session.auth"},
				// },
				// // Teams
				// {
				// 	Paths: []string{
				// 		"POST:/organizations/{organization_id}/teams",
				// 		"GET:/organizations/{organization_id}/teams",
				// 		"PATCH:/organizations/{organization_id}/teams/{team_id}",
				// 		"DELETE:/organizations/{organization_id}/teams/{team_id}",
				// 	},
				// 	Plugins: []string{"session.auth"},
				// },
				// // Team Members
				// {
				// 	Paths: []string{
				// 		"POST:/organizations/{organization_id}/teams/{team_id}/members",
				// 		"GET:/organizations/{organization_id}/teams/{team_id}/members",
				// 		"GET:/organizations/{organization_id}/teams/{team_id}/members/{member_id}",
				// 		"DELETE:/organizations/{organization_id}/teams/{team_id}/members/{member_id}",
				// 	},
				// 	Plugins: []string{"session.auth"},
				// },
			},
		),
	)

	authInstance := authula.New(&authula.AuthConfig{
		Config: config,
		Plugins: []authulamodels.Plugin{
			// This cannot be disabled
			emailplugin.New(emailplugintypes.EmailPluginConfig{
				Enabled:     true,
				Provider:    emailplugintypes.ProviderSMTP,
				FromAddress: "email@domain.com",
			}),
			emailpasswordplugin.New(emailpasswordtypes.EmailPasswordPluginConfig{
				Enabled:                     true,
				MinPasswordLength:           8,
				MaxPasswordLength:           128,
				DisableSignUp:               false,
				RequireEmailVerification:    false,
				AutoSignIn:                  true,
				SendEmailOnSignUp:           false,
				SendEmailOnSignIn:           false,
				EmailVerificationExpiresIn:  24 * time.Hour,
				PasswordResetExpiresIn:      time.Hour,
				RequestEmailChangeExpiresIn: time.Hour,
			}),
			sessionplugin.New(sessionplugin.SessionPluginConfig{
				Enabled: true,
			}),
			accesscontrolplugin.New(accesscontrolplugintypes.AccessControlPluginConfig{
				Enabled: true,
			}),
			adminplugin.New(adminplugintypes.AdminPluginConfig{
				Enabled:                   true,
				ImpersonationMaxExpiresIn: 15 * time.Minute,
			}),
			csrfplugin.New(csrfplugin.CSRFPluginConfig{
				Enabled:                false,
				CookieName:             "csrf_token",
				HeaderName:             "X-CSRF-TOKEN",
				Secure:                 lib.Env.NodeEnv == lib.NodeEnvProduction,
				SameSite:               "strict",
				EnableHeaderProtection: true,
			}),
			// ratelimitplugin.New(ratelimitplugin.RateLimitPluginConfig{
			// 	Enabled:  false,
			// 	Provider: ratelimitplugin.RateLimitProviderRedis,
			// }),
			// organizationsplugin.New(&organizationsplugintypes.OrganizationsPluginConfig{
			// 	Enabled:                          true,
			// 	OrganizationsLimit:               10,
			// 	MembersLimit:                     100,
			// 	InvitationsLimit:                 100,
			// 	InvitationExpiresIn:              24 * time.Hour,
			// 	RequireEmailVerifiedOnInvitation: true,
			// }),
		},
	})

	return authInstance, nil
}
