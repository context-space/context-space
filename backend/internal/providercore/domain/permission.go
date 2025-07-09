package domain

type Permission struct {
	Identifier  string   `json:"identifier"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OAuthScopes []string `json:"oauth_scopes"`
}

func NewPermission(identifier, name, description string, oauthScopes []string) *Permission {
	return &Permission{
		Identifier:  identifier,
		Name:        name,
		Description: description,
		OAuthScopes: oauthScopes,
	}
}

type PermissionSet map[string]Permission

func NewPermissionSet(permissions []Permission) PermissionSet {
	permissionsMap := make(map[string]Permission)
	for _, permission := range permissions {
		permissionsMap[permission.Identifier] = permission
	}
	return permissionsMap
}

func (p PermissionSet) RequiredOAuthScopes(permissionList []Permission) []string {
	uniqueScopes := make(map[string]bool)
	for _, permission := range permissionList {
		for _, scope := range p[permission.Identifier].OAuthScopes {
			uniqueScopes[scope] = true
		}
	}

	scopes := make([]string, 0, len(uniqueScopes))
	for scope := range uniqueScopes {
		scopes = append(scopes, scope)
	}
	return scopes
}

func (p PermissionSet) RequiredOAuthScopesByIdentifiers(permissionIdentifiers []string) []string {
	requiredPermissions := make([]Permission, 0, len(permissionIdentifiers))
	for _, identifier := range permissionIdentifiers {
		requiredPermissions = append(requiredPermissions, p[identifier])
	}
	return p.RequiredOAuthScopes(requiredPermissions)
}

func (p PermissionSet) CheckMissingPermissions(requiredPermissionList []Permission, authorizedScopes []string) (bool, []Permission) {
	authorizedScopeMap := make(map[string]bool)

	// Create a map of provided scopes for O(1) lookup
	for _, scope := range authorizedScopes {
		authorizedScopeMap[scope] = true
	}

	// Track missing permissions using a map to ensure uniqueness
	missingPermissionsMap := make(map[string]Permission)
	allScopesPresent := true

	// Check each permission's required scopes
	for _, permission := range requiredPermissionList {
		permissionScopes := p[permission.Identifier].OAuthScopes
		hasAllScopes := true

		// Check if all scopes for this permission are present
		for _, requiredScope := range permissionScopes {
			if !authorizedScopeMap[requiredScope] {
				hasAllScopes = false
				allScopesPresent = false
				break
			}
		}

		// If any scope is missing for this permission, add it to missing list
		if !hasAllScopes {
			missingPermissionsMap[permission.Identifier] = permission
		}
	}

	// Convert map to slice for return value
	missingPermissions := make([]Permission, 0, len(missingPermissionsMap))
	for _, permission := range missingPermissionsMap {
		missingPermissions = append(missingPermissions, permission)
	}

	return allScopesPresent, missingPermissions
}

func (p PermissionSet) CheckMissingPermissionsByIdentifiers(requiredPermissionIdentifiers []string, authorizedScopes []string) (bool, []string) {
	requiredPermissions := make([]Permission, 0, len(requiredPermissionIdentifiers))
	for _, identifier := range requiredPermissionIdentifiers {
		requiredPermissions = append(requiredPermissions, p[identifier])
	}

	allScopesPresent, missingPermissions := p.CheckMissingPermissions(requiredPermissions, authorizedScopes)

	missingPermissionsIdentifiers := make([]string, 0, len(missingPermissions))
	for _, permission := range missingPermissions {
		missingPermissionsIdentifiers = append(missingPermissionsIdentifiers, permission.Identifier)
	}
	return allScopesPresent, missingPermissionsIdentifiers
}

// GetPermissionIdentifiersFromScopes returns all permission identifiers that are associated with any of the provided scopes
func (p PermissionSet) GetPermissionIdentifiersFromScopes(scopes []string) []string {
	// Create a map for O(1) lookup of input scopes
	scopeMap := make(map[string]bool)
	for _, scope := range scopes {
		scopeMap[scope] = true
	}

	// Track matched permission identifiers using a map to ensure uniqueness
	matchedIdentifiers := make(map[string]bool)

	// Check each permission to see if any of its scopes match the input scopes
	for identifier, permission := range p {
		for _, permScope := range permission.OAuthScopes {
			if scopeMap[permScope] {
				matchedIdentifiers[identifier] = true
				break
			}
		}
	}

	// Convert map to slice for return value
	identifiers := make([]string, 0, len(matchedIdentifiers))
	for identifier := range matchedIdentifiers {
		identifiers = append(identifiers, identifier)
	}

	return identifiers
}
