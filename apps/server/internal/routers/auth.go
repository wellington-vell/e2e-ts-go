package routers

import (
	"net/http"

	actypes "github.com/Authula/authula/plugins/access-control/types"
	admtypes "github.com/Authula/authula/plugins/admin/types"
	eptypes "github.com/Authula/authula/plugins/email-password/types"
	"github.com/go-chi/chi/v5"

	"server/internal/models"
)

var authProxy http.Handler

// SetAuthProxy sets the authula handler that all wrapper routes will forward to.
func SetAuthProxy(handler http.Handler) {
	authProxy = handler
}

// RegisterAuthRoutes registers explicit chi routes for every authula endpoint so
// swaggo can discover them and generate OpenAPI definitions.
func RegisterAuthRoutes(r chi.Router) {
	// Core
	r.Get("/auth/me", AuthGetMe)
	r.Post("/auth/sign-out", AuthSignOut)

	// Email-Password
	r.Post("/auth/email-password/sign-up", AuthEmailPasswordSignUp)
	r.Post("/auth/email-password/sign-in", AuthEmailPasswordSignIn)
	r.Get("/auth/email-password/verify-email", AuthEmailPasswordVerifyEmail)
	r.Post("/auth/email-password/send-email-verification", AuthEmailPasswordSendEmailVerification)
	r.Post("/auth/email-password/request-password-reset", AuthEmailPasswordRequestPasswordReset)
	r.Post("/auth/email-password/change-password", AuthEmailPasswordChangePassword)
	r.Post("/auth/email-password/request-email-change", AuthEmailPasswordRequestEmailChange)

	// Admin - Users
	r.Post("/auth/admin/users", AuthAdminCreateUser)
	r.Get("/auth/admin/users", AuthAdminGetAllUsers)
	r.Get("/auth/admin/users/{user_id}", AuthAdminGetUserByID)
	r.Patch("/auth/admin/users/{user_id}", AuthAdminUpdateUser)
	r.Delete("/auth/admin/users/{user_id}", AuthAdminDeleteUser)

	// Admin - Accounts
	r.Post("/auth/admin/users/{user_id}/accounts", AuthAdminCreateAccount)
	r.Get("/auth/admin/users/{user_id}/accounts", AuthAdminGetUserAccounts)
	r.Get("/auth/admin/accounts/{id}", AuthAdminGetAccountByID)
	r.Patch("/auth/admin/accounts/{id}", AuthAdminUpdateAccount)
	r.Delete("/auth/admin/accounts/{id}", AuthAdminDeleteAccount)

	// Admin - User State
	r.Get("/auth/admin/users/{user_id}/state", AuthAdminGetUserState)       //
	r.Post("/auth/admin/users/{user_id}/state", AuthAdminCreateUserState)   //
	r.Patch("/auth/admin/users/{user_id}/state", AuthAdminUpdateUserState)  //
	r.Delete("/auth/admin/users/{user_id}/state", AuthAdminDeleteUserState) //
	r.Get("/auth/admin/users/states/banned", AuthAdminGetBannedUserStates)
	r.Post("/auth/admin/users/{user_id}/ban", AuthAdminBanUser)
	r.Post("/auth/admin/users/{user_id}/unban", AuthAdminUnbanUser)

	// Admin - Session State
	r.Get("/auth/admin/sessions/{session_id}/state", AuthAdminGetSessionState)
	r.Post("/auth/admin/sessions/{session_id}/state", AuthAdminCreateSessionState)
	r.Patch("/auth/admin/sessions/{session_id}/state", AuthAdminUpdateSessionState)
	r.Delete("/auth/admin/sessions/{session_id}/state", AuthAdminDeleteSessionState)
	r.Post("/auth/admin/sessions/{session_id}/revoke", AuthAdminRevokeSession)
	r.Get("/auth/admin/sessions/states/revoked", AuthAdminGetRevokedSessionStates)
	r.Get("/auth/admin/users/{user_id}/sessions", AuthAdminGetUserSessions)

	// Admin - Impersonation
	r.Get("/auth/admin/impersonations", AuthAdminGetAllImpersonations)
	r.Get("/auth/admin/impersonations/{impersonation_id}", AuthAdminGetImpersonationByID) //
	r.Post("/auth/admin/impersonations", AuthAdminStartImpersonation)
	r.Post("/auth/admin/impersonations/{impersonation_id}/stop", AuthAdminStopImpersonation)

	// Access Control - Roles
	r.Post("/auth/access-control/roles", AuthAccessControlCreateRole)
	r.Get("/auth/access-control/roles", AuthAccessControlGetAllRoles)
	r.Get("/auth/access-control/roles/by-name/{role_name}", AuthAccessControlGetRoleByName) //
	r.Get("/auth/access-control/roles/{role_id}", AuthAccessControlGetRoleByID)             //
	r.Patch("/auth/access-control/roles/{role_id}", AuthAccessControlUpdateRole)
	r.Delete("/auth/access-control/roles/{role_id}", AuthAccessControlDeleteRole)

	// Access Control - Permissions
	r.Post("/auth/access-control/permissions", AuthAccessControlCreatePermission)
	r.Get("/auth/access-control/permissions", AuthAccessControlGetAllPermissions)
	r.Get("/auth/access-control/permissions/{permission_id}", AuthAccessControlGetPermissionByID) //
	r.Patch("/auth/access-control/permissions/{permission_id}", AuthAccessControlUpdatePermission)
	r.Delete("/auth/access-control/permissions/{permission_id}", AuthAccessControlDeletePermission)

	// Access Control - Role Permissions
	r.Post("/auth/access-control/roles/{role_id}/permissions", AuthAccessControlAddRolePermission)
	r.Get("/auth/access-control/roles/{role_id}/permissions", AuthAccessControlGetRolePermissions)
	r.Put("/auth/access-control/roles/{role_id}/permissions", AuthAccessControlReplaceRolePermissions) //
	r.Delete("/auth/access-control/roles/{role_id}/permissions/{permission_id}", AuthAccessControlRemoveRolePermission)

	// Access Control - User Access
	r.Get("/auth/access-control/users/{user_id}/roles", AuthAccessControlGetUserRoles)
	r.Post("/auth/access-control/users/{user_id}/roles", AuthAccessControlAssignUserRole)
	r.Put("/auth/access-control/users/{user_id}/roles", AuthAccessControlReplaceUserRoles)
	r.Delete("/auth/access-control/users/{user_id}/roles/{role_id}", AuthAccessControlRemoveUserRole)
	r.Get("/auth/access-control/users/{user_id}/permissions", AuthAccessControlGetUserPermissions)
	r.Post("/auth/access-control/users/{user_id}/permissions/check", AuthAccessControlCheckUserPermissions)
}

// ============================================================================
// Core Auth
// ============================================================================

// @Summary Get current authenticated user
// @Description Returns the currently authenticated user and active session.
// @Tags Auth
// @Produce json
// @Success 200 {object} models.GetMeResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/me [get]
func AuthGetMe(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Sign out
// @Description Signs out the current user. Optionally signs out all sessions.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.SignOutRequest false "Sign out request"
// @Success 200 {object} models.SignOutResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/sign-out [post]
func AuthSignOut(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Email-Password
// ============================================================================

// @Summary Sign up with email and password
// @Description Creates a new user account with email and password.
// @Tags Auth Email Password
// @Accept json
// @Produce json
// @Param request body eptypes.SignUpRequest true "Sign up request"
// @Success 201 {object} eptypes.SignUpResponse
// @Failure 400 {string} string "Bad request"
// @Router /auth/email-password/sign-up [post]
func AuthEmailPasswordSignUp(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Sign in with email and password
// @Description Authenticates a user with email and password.
// @Tags Auth Email Password
// @Accept json
// @Produce json
// @Param request body eptypes.SignInRequest true "Sign in request"
// @Success 200 {object} eptypes.SignInResponse
// @Failure 400 {string} string "Bad request"
// @Router /auth/email-password/sign-in [post]
func AuthEmailPasswordSignIn(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Verify email
// @Description Verifies a user's email using a token. Redirects if callback_url is provided.
// @Tags Auth Email Password
// @Produce json
// @Param token query string true "Verification token"
// @Param callback_url query string false "Callback URL"
// @Success 200 {object} map[string]any "Email verified successfully"
// @Failure 400 {string} string "Bad request"
// @Router /auth/email-password/verify-email [get]
func AuthEmailPasswordVerifyEmail(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Send email verification
// @Description Sends an email verification link to the specified email.
// @Tags Auth Email Password
// @Accept json
// @Produce json
// @Param request body eptypes.SendEmailVerificationRequest true "Send email verification request"
// @Success 200 {object} map[string]any "Email verification sent"
// @Failure 400 {string} string "Bad request"
// @Router /auth/email-password/send-email-verification [post]
func AuthEmailPasswordSendEmailVerification(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Request password reset
// @Description Requests a password reset for the specified email.
// @Tags Auth Email Password
// @Accept json
// @Produce json
// @Param request body eptypes.RequestPasswordResetRequest true "Request password reset request"
// @Success 200 {object} map[string]any "Password reset requested"
// @Failure 400 {string} string "Bad request"
// @Router /auth/email-password/request-password-reset [post]
func AuthEmailPasswordRequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Change password
// @Description Changes the user's password using a reset token.
// @Tags Auth Email Password
// @Accept json
// @Produce json
// @Param request body eptypes.ChangePasswordRequest true "Change password request"
// @Success 200 {object} eptypes.ChangePasswordResponse
// @Failure 400 {string} string "Bad request"
// @Router /auth/email-password/change-password [post]
func AuthEmailPasswordChangePassword(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Request email change
// @Description Requests an email change for the authenticated user.
// @Tags Auth Email Password
// @Accept json
// @Produce json
// @Param request body eptypes.RequestEmailChangeRequest true "Request email change request"
// @Success 200 {object} map[string]any "Email change requested"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/email-password/request-email-change [post]
func AuthEmailPasswordRequestEmailChange(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Admin - Users
// ============================================================================

// @Summary Create user
// @Description Creates a new user account (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param request body admtypes.CreateUserRequest true "Create user request"
// @Success 201 {object} admtypes.CreateUserResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users [post]
func AuthAdminCreateUser(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get all users
// @Description Retrieves a paginated list of all users (admin only).
// @Tags Auth Admin
// @Produce json
// @Param cursor query string false "Pagination cursor"
// @Param limit query int false "Pagination limit"
// @Success 200 {object} admtypes.UsersPage
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users [get]
func AuthAdminGetAllUsers(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get user by ID
// @Description Retrieves a single user by ID (admin only).
// @Tags Auth Admin
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} admtypes.GetUserByIDResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id} [get]
func AuthAdminGetUserByID(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Update user
// @Description Updates an existing user by ID (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body admtypes.UpdateUserRequest true "Update user request"
// @Success 200 {object} admtypes.UpdateUserResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id} [patch]
func AuthAdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Delete user
// @Description Deletes a user by ID (admin only).
// @Tags Auth Admin
// @Param user_id path string true "User ID"
// @Success 200 {object} admtypes.DeleteUserResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id} [delete]
func AuthAdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Admin - Accounts
// ============================================================================

// @Summary Create account
// @Description Creates an account for a user (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body admtypes.CreateAccountRequest true "Create account request"
// @Success 201 {object} admtypes.CreateAccountResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users/{user_id}/accounts [post]
func AuthAdminCreateAccount(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get user accounts
// @Description Retrieves all accounts for a user (admin only).
// @Tags Auth Admin
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} admtypes.UserAccountsResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id}/accounts [get]
func AuthAdminGetUserAccounts(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get account by ID
// @Description Retrieves a single account by ID (admin only).
// @Tags Auth Admin
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} admtypes.GetAccountByIDResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/accounts/{id} [get]
func AuthAdminGetAccountByID(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Update account
// @Description Updates an existing account by ID (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body admtypes.UpdateAccountRequest true "Update account request"
// @Success 200 {object} admtypes.UpdateAccountResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/accounts/{id} [patch]
func AuthAdminUpdateAccount(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Delete account
// @Description Deletes an account by ID (admin only).
// @Tags Auth Admin
// @Param id path string true "Account ID"
// @Success 200 {object} admtypes.DeleteAccountResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/accounts/{id} [delete]
func AuthAdminDeleteAccount(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Admin - User State
// ============================================================================

// @Summary Get user state
// @Description Retrieves the state (banned status) of a user (admin only).
// @Tags Auth Admin
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} admtypes.GetUserStateResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id}/state [get]
func AuthAdminGetUserState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Create user state
// @Description Creates a state record for a user (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body admtypes.CreateUserStateRequest true "Create user state request"
// @Success 201 {object} admtypes.UpsertUserStateResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users/{user_id}/state [post]
func AuthAdminCreateUserState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Update user state
// @Description Updates the state record for a user (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body admtypes.UpsertUserStateRequest true "Update user state request"
// @Success 200 {object} admtypes.UpsertUserStateResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users/{user_id}/state [patch]
func AuthAdminUpdateUserState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Delete user state
// @Description Deletes the state record for a user (admin only).
// @Tags Auth Admin
// @Param user_id path string true "User ID"
// @Success 200 {object} admtypes.DeleteUserStateResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id}/state [delete]
func AuthAdminDeleteUserState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get banned user states
// @Description Retrieves all banned user states (admin only).
// @Tags Auth Admin
// @Produce json
// @Success 200 {array} admtypes.AdminUserState
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users/states/banned [get]
func AuthAdminGetBannedUserStates(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Ban user
// @Description Bans a user by ID (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body admtypes.BanUserRequest true "Ban user request"
// @Success 200 {object} admtypes.BanUserResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users/{user_id}/ban [post]
func AuthAdminBanUser(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Unban user
// @Description Unbans a user by ID (admin only).
// @Tags Auth Admin
// @Param user_id path string true "User ID"
// @Success 200 {object} admtypes.UnbanUserResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/users/{user_id}/unban [post]
func AuthAdminUnbanUser(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Admin - Session State
// ============================================================================

// @Summary Get session state
// @Description Retrieves the state of a session (admin only).
// @Tags Auth Admin
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} admtypes.GetSessionStateResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/sessions/{session_id}/state [get]
func AuthAdminGetSessionState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Create session state
// @Description Creates a state record for a session (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Param request body admtypes.CreateSessionStateRequest true "Create session state request"
// @Success 201 {object} admtypes.UpsertSessionStateResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/sessions/{session_id}/state [post]
func AuthAdminCreateSessionState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Update session state
// @Description Updates the state record for a session (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Param request body admtypes.UpsertSessionStateRequest true "Update session state request"
// @Success 200 {object} admtypes.UpsertSessionStateResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/sessions/{session_id}/state [patch]
func AuthAdminUpdateSessionState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Delete session state
// @Description Deletes the state record for a session (admin only).
// @Tags Auth Admin
// @Param session_id path string true "Session ID"
// @Success 200 {object} admtypes.DeleteSessionStateResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/sessions/{session_id}/state [delete]
func AuthAdminDeleteSessionState(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Revoke session
// @Description Revokes a session by ID (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Param request body admtypes.RevokeSessionRequest true "Revoke session request"
// @Success 200 {object} admtypes.RevokeSessionResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/sessions/{session_id}/revoke [post]
func AuthAdminRevokeSession(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get revoked session states
// @Description Retrieves all revoked session states (admin only).
// @Tags Auth Admin
// @Produce json
// @Success 200 {array} admtypes.AdminSessionState
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/sessions/states/revoked [get]
func AuthAdminGetRevokedSessionStates(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get user sessions
// @Description Retrieves all sessions for a user (admin only).
// @Tags Auth Admin
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} admtypes.AdminUserSession
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/users/{user_id}/sessions [get]
func AuthAdminGetUserSessions(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Admin - Impersonation
// ============================================================================

// @Summary Get all impersonations
// @Description Retrieves all impersonation records (admin only).
// @Tags Auth Admin
// @Produce json
// @Success 200 {array} admtypes.Impersonation
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/impersonations [get]
func AuthAdminGetAllImpersonations(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get impersonation by ID
// @Description Retrieves a single impersonation record by ID (admin only).
// @Tags Auth Admin
// @Produce json
// @Param impersonation_id path string true "Impersonation ID"
// @Success 200 {object} admtypes.GetImpersonationByIDResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/impersonations/{impersonation_id} [get]
func AuthAdminGetImpersonationByID(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Start impersonation
// @Description Starts an impersonation session for a target user (admin only).
// @Tags Auth Admin
// @Accept json
// @Produce json
// @Param request body admtypes.StartImpersonationRequest true "Start impersonation request"
// @Success 201 {object} admtypes.StartImpersonationResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/admin/impersonations [post]
func AuthAdminStartImpersonation(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Stop impersonation
// @Description Stops an active impersonation session by ID.
// @Tags Auth Admin
// @Produce json
// @Param impersonation_id path string true "Impersonation ID"
// @Success 200 {object} admtypes.StopImpersonationResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/admin/impersonations/{impersonation_id}/stop [post]
func AuthAdminStopImpersonation(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Access Control - Roles
// ============================================================================

// @Summary Create role
// @Description Creates a new role (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param request body actypes.CreateRoleRequest true "Create role request"
// @Success 201 {object} actypes.CreateRoleResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/roles [post]
func AuthAccessControlCreateRole(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get all roles
// @Description Retrieves all roles (admin only).
// @Tags Auth Access Control
// @Produce json
// @Success 200 {array} actypes.Role
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/roles [get]
func AuthAccessControlGetAllRoles(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get role by name
// @Description Retrieves a single role by name (admin only).
// @Tags Auth Access Control
// @Produce json
// @Param role_name path string true "Role name"
// @Success 200 {object} actypes.Role
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/roles/by-name/{role_name} [get]
func AuthAccessControlGetRoleByName(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get role by ID
// @Description Retrieves a single role by ID with its permissions (admin only).
// @Tags Auth Access Control
// @Produce json
// @Param role_id path string true "Role ID"
// @Success 200 {object} actypes.RoleDetails
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/roles/{role_id} [get]
func AuthAccessControlGetRoleByID(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Update role
// @Description Updates an existing role by ID (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param request body actypes.UpdateRoleRequest true "Update role request"
// @Success 200 {object} actypes.UpdateRoleResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/roles/{role_id} [patch]
func AuthAccessControlUpdateRole(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Delete role
// @Description Deletes a role by ID (admin only).
// @Tags Auth Access Control
// @Param role_id path string true "Role ID"
// @Success 200 {object} actypes.DeleteRoleResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/roles/{role_id} [delete]
func AuthAccessControlDeleteRole(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Access Control - Permissions
// ============================================================================

// @Summary Create permission
// @Description Creates a new permission (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param request body actypes.CreatePermissionRequest true "Create permission request"
// @Success 201 {object} actypes.CreatePermissionResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/permissions [post]
func AuthAccessControlCreatePermission(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get all permissions
// @Description Retrieves all permissions (admin only).
// @Tags Auth Access Control
// @Produce json
// @Success 200 {array} actypes.Permission
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/permissions [get]
func AuthAccessControlGetAllPermissions(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get permission by ID
// @Description Retrieves a single permission by ID (admin only).
// @Tags Auth Access Control
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} actypes.Permission
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/permissions/{permission_id} [get]
func AuthAccessControlGetPermissionByID(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Update permission
// @Description Updates an existing permission by ID (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param permission_id path string true "Permission ID"
// @Param request body actypes.UpdatePermissionRequest true "Update permission request"
// @Success 200 {object} actypes.UpdatePermissionResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/permissions/{permission_id} [patch]
func AuthAccessControlUpdatePermission(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Delete permission
// @Description Deletes a permission by ID (admin only).
// @Tags Auth Access Control
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} actypes.DeletePermissionResponse
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Router /auth/access-control/permissions/{permission_id} [delete]
func AuthAccessControlDeletePermission(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Access Control - Role Permissions
// ============================================================================

// @Summary Add permission to role
// @Description Adds a permission to a role (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param request body actypes.AddRolePermissionRequest true "Add role permission request"
// @Success 200 {object} actypes.AddRolePermissionResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/roles/{role_id}/permissions [post]
func AuthAccessControlAddRolePermission(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get role permissions
// @Description Retrieves all permissions assigned to a role (admin only).
// @Tags Auth Access Control
// @Produce json
// @Param role_id path string true "Role ID"
// @Success 200 {array} actypes.UserPermissionInfo
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/roles/{role_id}/permissions [get]
func AuthAccessControlGetRolePermissions(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Replace role permissions
// @Description Replaces all permissions for a role (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param request body actypes.ReplaceRolePermissionsRequest true "Replace role permissions request"
// @Success 200 {object} actypes.ReplaceRolePermissionResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/roles/{role_id}/permissions [put]
func AuthAccessControlReplaceRolePermissions(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Remove permission from role
// @Description Removes a permission from a role (admin only).
// @Tags Auth Access Control
// @Param role_id path string true "Role ID"
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} actypes.RemoveRolePermissionResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/roles/{role_id}/permissions/{permission_id} [delete]
func AuthAccessControlRemoveRolePermission(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// ============================================================================
// Access Control - User Access
// ============================================================================

// @Summary Get user roles
// @Description Retrieves all roles assigned to a user (admin only).
// @Tags Auth Access Control
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} actypes.UserRoleInfo
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/users/{user_id}/roles [get]
func AuthAccessControlGetUserRoles(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Assign role to user
// @Description Assigns a role to a user (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body actypes.AssignUserRoleRequest true "Assign user role request"
// @Success 200 {object} actypes.AssignUserRoleResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/users/{user_id}/roles [post]
func AuthAccessControlAssignUserRole(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Replace user roles
// @Description Replaces all roles for a user (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body actypes.ReplaceUserRolesRequest true "Replace user roles request"
// @Success 200 {object} actypes.ReplaceUserRolesResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/users/{user_id}/roles [put]
func AuthAccessControlReplaceUserRoles(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Remove role from user
// @Description Removes a role from a user (admin only).
// @Tags Auth Access Control
// @Param user_id path string true "User ID"
// @Param role_id path string true "Role ID"
// @Success 200 {object} actypes.RemoveUserRoleResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/users/{user_id}/roles/{role_id} [delete]
func AuthAccessControlRemoveUserRole(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Get user permissions
// @Description Retrieves all permissions for a user (admin only).
// @Tags Auth Access Control
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} actypes.GetUserPermissionsResponse
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/users/{user_id}/permissions [get]
func AuthAccessControlGetUserPermissions(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// @Summary Check user permissions
// @Description Checks if a user has the specified permissions (admin only).
// @Tags Auth Access Control
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body actypes.CheckUserPermissionsRequest true "Check user permissions request"
// @Success 200 {object} actypes.CheckUserPermissionsResponse
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/access-control/users/{user_id}/permissions/check [post]
func AuthAccessControlCheckUserPermissions(w http.ResponseWriter, r *http.Request) {
	authProxy.ServeHTTP(w, r)
}

// _ensureImportsUsed keeps the type-package imports alive for swaggo parsing.
//
//nolint:unused
func _ensureImportsUsed() {
	_ = models.GetMeResponse{}
	_ = models.SignOutRequest{}
	_ = models.SignOutResponse{}
	_ = eptypes.SignUpRequest{}
	_ = eptypes.SignUpResponse{}
	_ = eptypes.SignInRequest{}
	_ = eptypes.SignInResponse{}
	_ = eptypes.SendEmailVerificationRequest{}
	_ = eptypes.RequestPasswordResetRequest{}
	_ = eptypes.ChangePasswordRequest{}
	_ = eptypes.ChangePasswordResponse{}
	_ = eptypes.RequestEmailChangeRequest{}
	_ = admtypes.CreateUserRequest{}
	_ = admtypes.CreateUserResponse{}
	_ = admtypes.UsersPage{}
	_ = admtypes.GetUserByIDResponse{}
	_ = admtypes.UpdateUserRequest{}
	_ = admtypes.UpdateUserResponse{}
	_ = admtypes.DeleteUserResponse{}
	_ = admtypes.CreateAccountRequest{}
	_ = admtypes.CreateAccountResponse{}
	_ = admtypes.UserAccountsResponse{}
	_ = admtypes.GetAccountByIDResponse{}
	_ = admtypes.UpdateAccountRequest{}
	_ = admtypes.UpdateAccountResponse{}
	_ = admtypes.DeleteAccountResponse{}
	_ = admtypes.GetUserStateResponse{}
	_ = admtypes.CreateUserStateRequest{}
	_ = admtypes.UpsertUserStateResponse{}
	_ = admtypes.UpsertUserStateRequest{}
	_ = admtypes.DeleteUserStateResponse{}
	_ = admtypes.AdminUserState{}
	_ = admtypes.BanUserRequest{}
	_ = admtypes.BanUserResponse{}
	_ = admtypes.UnbanUserResponse{}
	_ = admtypes.GetSessionStateResponse{}
	_ = admtypes.CreateSessionStateRequest{}
	_ = admtypes.UpsertSessionStateResponse{}
	_ = admtypes.UpsertSessionStateRequest{}
	_ = admtypes.DeleteSessionStateResponse{}
	_ = admtypes.AdminSessionState{}
	_ = admtypes.AdminUserSession{}
	_ = admtypes.RevokeSessionRequest{}
	_ = admtypes.RevokeSessionResponse{}
	_ = admtypes.Impersonation{}
	_ = admtypes.GetImpersonationByIDResponse{}
	_ = admtypes.StartImpersonationRequest{}
	_ = admtypes.StartImpersonationResponse{}
	_ = admtypes.StopImpersonationResponse{}
	_ = actypes.CreateRoleRequest{}
	_ = actypes.CreateRoleResponse{}
	_ = actypes.Role{}
	_ = actypes.UpdateRoleRequest{}
	_ = actypes.UpdateRoleResponse{}
	_ = actypes.DeleteRoleResponse{}
	_ = actypes.RoleDetails{}
	_ = actypes.CreatePermissionRequest{}
	_ = actypes.CreatePermissionResponse{}
	_ = actypes.Permission{}
	_ = actypes.UpdatePermissionRequest{}
	_ = actypes.UpdatePermissionResponse{}
	_ = actypes.DeletePermissionResponse{}
	_ = actypes.AddRolePermissionRequest{}
	_ = actypes.AddRolePermissionResponse{}
	_ = actypes.ReplaceRolePermissionsRequest{}
	_ = actypes.ReplaceRolePermissionResponse{}
	_ = actypes.RemoveRolePermissionResponse{}
	_ = actypes.UserPermissionInfo{}
	_ = actypes.UserRoleInfo{}
	_ = actypes.ReplaceUserRolesRequest{}
	_ = actypes.ReplaceUserRolesResponse{}
	_ = actypes.AssignUserRoleRequest{}
	_ = actypes.AssignUserRoleResponse{}
	_ = actypes.RemoveUserRoleResponse{}
	_ = actypes.GetUserPermissionsResponse{}
	_ = actypes.CheckUserPermissionsRequest{}
	_ = actypes.CheckUserPermissionsResponse{}
}
