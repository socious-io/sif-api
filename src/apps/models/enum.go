package models

import (
	"database/sql/driver"
	"fmt"
)

type (
	IdentityType            string
	OauthConnectedProviders string
	WalletENV               string
	ProjectStatus           string
)

const (
	IdentityTypeUsers         IdentityType = "users"
	IdentityTypeOrganizations IdentityType = "organizations"

	OauthConnectedProvidersSociousID OauthConnectedProviders = "SOCIOUS_ID"

	WalletCardanoENV WalletENV = "CARDANO"

	ProjectStatusDraft  ProjectStatus = "DRAFT"
	ProjectStatusExpire ProjectStatus = "EXPIRE"
	ProjectStatusActive ProjectStatus = "ACTIVE"
)

// ------------------------------------------------------

func (v *IdentityType) Scan(value interface{}) error {
	return scanEnum(value, (*string)(v))
}

func (v IdentityType) Value() (driver.Value, error) {
	return string(v), nil
}

// ------------------------------------------------------

func (v *OauthConnectedProviders) Scan(value interface{}) error {
	return scanEnum(value, (*string)(v))
}

func (v OauthConnectedProviders) Value() (driver.Value, error) {
	return string(v), nil
}

// ------------------------------------------------------

func (v *WalletENV) Scan(value interface{}) error {
	return scanEnum(value, (*string)(v))
}

func (v WalletENV) Value() (driver.Value, error) {
	return string(v), nil
}

// ------------------------------------------------------

func (ps *ProjectStatus) Scan(value interface{}) error {
	return scanEnum(value, (*string)(ps))
}

func (ps ProjectStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// ------------------------------------------------------

// scanEnum is a helper function that converts an interface{} value to a string
// to support database scanning. It handles both byte slices and string values.
func scanEnum(value interface{}, target interface{}) error {
	switch v := value.(type) {
	case []byte:
		*target.(*string) = string(v) // Convert byte slice to string.
	case string:
		*target.(*string) = v // Assign string value.
	default:
		return fmt.Errorf("failed to scan type: %v", value) // Error on unsupported type.
	}
	return nil
}
