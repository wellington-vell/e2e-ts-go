package db

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Authula/authula"
	authulamodels "github.com/Authula/authula/models"
	accesscontrolplugin "github.com/Authula/authula/plugins/access-control"
	actypes "github.com/Authula/authula/plugins/access-control/types"
	adminplugin "github.com/Authula/authula/plugins/admin"
	admintypes "github.com/Authula/authula/plugins/admin/types"
	emailpasswordplugin "github.com/Authula/authula/plugins/email-password"
	"github.com/jaswdr/faker"
	"github.com/uptrace/bun"

	"server/internal/models"
)

func SeedDB(ctx context.Context, db *bun.DB, auth *authula.Auth) error {
	if err := truncateData(ctx, db, auth); err != nil {
		return fmt.Errorf("failed to truncate data: %w", err)
	}

	// Get plugins
	epPlugin, err := getEmailPasswordPlugin(auth)
	if err != nil {
		return err
	}

	adminPlugin, err := getAdminPlugin(auth)
	if err != nil {
		return err
	}

	acPlugin, err := getAccessControlPlugin(auth)
	if err != nil {
		return err
	}

	// Seed permissions
	permissions, err := seedPermissions(ctx, acPlugin)
	if err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	// Seed roles
	roles, err := seedRoles(ctx, acPlugin)
	if err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	// Assign permissions to roles
	if err := seedRolePermissions(ctx, db, roles, permissions); err != nil {
		return fmt.Errorf("failed to seed role permissions: %w", err)
	}

	// Create fixed admin user
	adminUser, err := createAdminUser(ctx, adminPlugin)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// Assign admin role to admin user
	if err := acPlugin.Api.AssignRoleToUser(ctx, adminUser.ID, actypes.AssignUserRoleRequest{
		RoleID: roles["admin"].ID,
	}, nil); err != nil {
		return fmt.Errorf("failed to assign admin role: %w", err)
	}

	// Create fake users
	f := faker.New()
	usersCount := 50
	userIDs := make([]string, 0, usersCount)
	for i := range usersCount {
		name := f.Person().Name()
		email := f.Internet().Email()
		password := "password123"

		result, err := epPlugin.Api.SignUp(ctx, name, email, password, nil, json.RawMessage(`{}`), nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to seed user %d (%s): %w", i, email, err)
		}
		userIDs = append(userIDs, result.User.ID)
	}

	// Assign roles to fake users: first 10% manager, rest user
	managerCount := max(usersCount/10, 1)
	for i, userID := range userIDs {
		roleName := "user"
		if i < managerCount {
			roleName = "manager"
		}
		if err := acPlugin.Api.AssignRoleToUser(ctx, userID, actypes.AssignUserRoleRequest{
			RoleID: roles[roleName].ID,
		}, nil); err != nil {
			return fmt.Errorf("failed to assign %s role to user %s: %w", roleName, userID, err)
		}
	}

	// Seed todos
	statuses := []models.TodoStatus{
		models.TodoStatusBacklog,
		models.TodoStatusTodo,
		models.TodoStatusInProgress,
		models.TodoStatusDone,
		models.TodoStatusCanceled,
	}

	priorities := []models.TodoPriority{
		models.TodoPriorityLow,
		models.TodoPriorityMedium,
		models.TodoPriorityHigh,
	}

	labels := []*string{
		nil,
		strPtr(models.TodoLabelBug),
		strPtr(models.TodoLabelFeature),
		strPtr(models.TodoLabelDoc),
	}

	// Include admin user in the pool so todos can be assigned to them too
	allUserIDs := append([]string{adminUser.ID}, userIDs...)

	seedCount := 5000
	for i := range seedCount {
		status := statuses[rand.Intn(len(statuses))]
		priority := priorities[rand.Intn(len(priorities))]
		label := labels[rand.Intn(len(labels))]

		var dueDate *models.Date
		if rand.Float64() > 0.3 {
			d := time.Now().AddDate(0, 0, rand.Intn(90)-30)
			dueDate = &models.Date{Time: d}
		}

		var completedAt *time.Time
		if status == models.TodoStatusDone {
			t := time.Now().Add(-time.Duration(rand.Intn(720)) * time.Hour)
			completedAt = &t
		}

		progress := 0
		switch status {
		case models.TodoStatusDone:
			progress = 100
		case models.TodoStatusInProgress:
			progress = rand.Intn(100)
		}

		userID := allUserIDs[rand.Intn(len(allUserIDs))]

		todo := models.Todo{
			Text:           f.Lorem().Sentence(rand.Intn(8) + 3),
			Status:         status,
			Label:          label,
			Priority:       priority,
			EstimatedHours: float64(rand.Intn(40)),
			ActualHours:    float64(rand.Intn(40)),
			Progress:       progress,
			Cost:           float64(rand.Intn(100000)) / 100,
			DueDate:        dueDate,
			CompletedAt:    completedAt,
			UserID:         &userID,
		}

		_, err := db.NewInsert().Model(&todo).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to seed todo %d: %w", i, err)
		}
	}

	return nil
}

func truncateData(ctx context.Context, db *bun.DB, auth *authula.Auth) error {
	// Truncate access-control tables first due to FK constraints
	if _, err := db.ExecContext(ctx, "DELETE FROM access_control_user_roles"); err != nil {
		return fmt.Errorf("failed to truncate access_control_user_roles: %w", err)
	}
	if _, err := db.ExecContext(ctx, "DELETE FROM access_control_role_permissions"); err != nil {
		return fmt.Errorf("failed to truncate access_control_role_permissions: %w", err)
	}
	if _, err := db.ExecContext(ctx, "DELETE FROM access_control_roles"); err != nil {
		return fmt.Errorf("failed to truncate access_control_roles: %w", err)
	}
	if _, err := db.ExecContext(ctx, "DELETE FROM access_control_permissions"); err != nil {
		return fmt.Errorf("failed to truncate access_control_permissions: %w", err)
	}

	// Truncate todos
	_, err := db.NewDelete().Model((*models.Todo)(nil)).Where("1=1").Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to truncate todos: %w", err)
	}

	// Delete users via admin API
	adminPlugin, err := getAdminPlugin(auth)
	if err != nil {
		return err
	}

	cursor := (*string)(nil)
	for {
		page, err := adminPlugin.Api.GetAllUsers(ctx, cursor, 100)
		if err != nil {
			return fmt.Errorf("failed to get users for deletion: %w", err)
		}

		for _, user := range page.Users {
			if err := adminPlugin.Api.DeleteUser(ctx, user.ID); err != nil {
				return fmt.Errorf("failed to delete user %s: %w", user.ID, err)
			}
		}

		if page.NextCursor == nil || len(page.Users) == 0 {
			break
		}
		cursor = page.NextCursor
	}

	return nil
}

func getEmailPasswordPlugin(auth *authula.Auth) (*emailpasswordplugin.EmailPasswordPlugin, error) {
	plugin := auth.PluginRegistry.GetPlugin("email_password")
	if plugin == nil {
		return nil, fmt.Errorf("email-password plugin not found")
	}
	epPlugin, ok := plugin.(*emailpasswordplugin.EmailPasswordPlugin)
	if !ok {
		return nil, fmt.Errorf("email-password plugin has unexpected type")
	}
	return epPlugin, nil
}

func getAdminPlugin(auth *authula.Auth) (*adminplugin.AdminPlugin, error) {
	plugin := auth.PluginRegistry.GetPlugin("admin")
	if plugin == nil {
		return nil, fmt.Errorf("admin plugin not found")
	}
	adminPlugin, ok := plugin.(*adminplugin.AdminPlugin)
	if !ok {
		return nil, fmt.Errorf("admin plugin has unexpected type")
	}
	return adminPlugin, nil
}

func getAccessControlPlugin(auth *authula.Auth) (*accesscontrolplugin.AccessControlPlugin, error) {
	plugin := auth.PluginRegistry.GetPlugin("access_control")
	if plugin == nil {
		return nil, fmt.Errorf("access-control plugin not found")
	}
	acPlugin, ok := plugin.(*accesscontrolplugin.AccessControlPlugin)
	if !ok {
		return nil, fmt.Errorf("access-control plugin has unexpected type")
	}
	return acPlugin, nil
}

func seedPermissions(ctx context.Context, acPlugin *accesscontrolplugin.AccessControlPlugin) (map[string]*actypes.Permission, error) {
	defs := []struct {
		key         string
		description string
		isSystem    bool
	}{
		{"users:read", "View users", true},
		{"users:manage", "Create, update, and delete users", true},
		{"roles:read", "View roles", true},
		{"roles:manage", "Create, update, and delete roles", true},
		{"permissions:read", "View permissions", true},
		{"permissions:manage", "Create, update, and delete permissions", true},
		{"sessions:read", "View sessions", true},
		{"sessions:manage", "Revoke and manage sessions", true},
		{"impersonation:read", "View impersonations", true},
		{"impersonation:manage", "Start and stop impersonations", true},
		{"todos:read", "View todos", false},
		{"todos:create", "Create todos", false},
		{"todos:update", "Update todos", false},
		{"todos:delete", "Delete todos", false},
	}

	permissions := make(map[string]*actypes.Permission, len(defs))
	for _, def := range defs {
		desc := def.description
		perm, err := acPlugin.Api.CreatePermission(ctx, actypes.CreatePermissionRequest{
			Key:         def.key,
			Description: &desc,
			IsSystem:    def.isSystem,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create permission %s: %w", def.key, err)
		}
		permissions[def.key] = perm
	}

	return permissions, nil
}

func seedRoles(ctx context.Context, acPlugin *accesscontrolplugin.AccessControlPlugin) (map[string]*actypes.Role, error) {
	defs := []struct {
		name        string
		description string
		weight      int
		isSystem    bool
	}{
		{"admin", "Administrator with full access", 100, true},
		{"manager", "Manager with elevated access", 50, false},
		{"user", "Standard user", 10, true},
	}

	roles := make(map[string]*actypes.Role, len(defs))
	for _, def := range defs {
		desc := def.description
		w := def.weight
		role, err := acPlugin.Api.CreateRole(ctx, actypes.CreateRoleRequest{
			Name:        def.name,
			Description: &desc,
			Weight:      &w,
			IsSystem:    def.isSystem,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create role %s: %w", def.name, err)
		}
		roles[def.name] = role
	}

	return roles, nil
}

func seedRolePermissions(ctx context.Context, db *bun.DB, roles map[string]*actypes.Role, permissions map[string]*actypes.Permission) error {
	mapping := map[string][]string{
		"admin": {
			"users:read", "users:manage",
			"roles:read", "roles:manage",
			"permissions:read", "permissions:manage",
			"sessions:read", "sessions:manage",
			"impersonation:read", "impersonation:manage",
			"todos:read", "todos:create", "todos:update", "todos:delete",
		},
		"manager": {
			"todos:read", "todos:create", "todos:update", "todos:delete",
			"users:read",
			"roles:read",
			"permissions:read",
			"sessions:read",
		},
		"user": {
			"todos:read", "todos:create", "todos:update",
		},
	}

	// Insert directly into the join table to bypass the service-layer guard
	// that rejects writes to system roles and system permissions. Seeding is
	// the trusted bootstrap path those guards are designed to protect at runtime.
	now := time.Now().UTC()
	for roleName, permKeys := range mapping {
		role := roles[roleName]
		for _, key := range permKeys {
			rp := &actypes.RolePermission{
				RoleID:       role.ID,
				PermissionID: permissions[key].ID,
				GrantedAt:    now,
			}
			if _, err := db.NewInsert().Model(rp).Exec(ctx); err != nil {
				return fmt.Errorf("failed to assign permission %s to role %s: %w", key, roleName, err)
			}
		}
	}

	return nil
}

func createAdminUser(ctx context.Context, adminPlugin *adminplugin.AdminPlugin) (*authulamodels.User, error) {
	adminEmail := os.Getenv("SEED_ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}
	adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	verified := true
	adminUser, err := adminPlugin.Api.CreateUser(ctx, admintypes.CreateUserRequest{
		Name:          "Admin",
		Email:         adminEmail,
		EmailVerified: &verified,
		Metadata:      json.RawMessage(`{}`),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	_, err = adminPlugin.Api.CreateAccount(ctx, adminUser.ID, admintypes.CreateAccountRequest{
		ProviderID: "email",
		AccountID:  adminEmail,
		Password:   &adminPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create admin account: %w", err)
	}

	return adminUser, nil
}

func strPtr(s string) *string {
	return &s
}
