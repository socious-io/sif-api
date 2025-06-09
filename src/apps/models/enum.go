package models

import (
	"database/sql/driver"
	"fmt"
)

type (
	IdentityType              string
	OauthConnectedProviders   string
	WalletENV                 string
	ProjectStatus             string
	KybVerificationStatusType string
	DonationStatus            string
	OrganizationStatus        string
	PaymentType               string
	ProjectCategory           string
)

const (
	IdentityTypeUsers         IdentityType = "users"
	IdentityTypeOrganizations IdentityType = "organizations"

	OauthConnectedProvidersSociousID OauthConnectedProviders = "SOCIOUS_ID"

	WalletCardanoENV WalletENV = "CARDANO"

	ProjectStatusDraft  ProjectStatus = "DRAFT"
	ProjectStatusExpire ProjectStatus = "EXPIRE"
	ProjectStatusActive ProjectStatus = "ACTIVE"

	KYBStatusPending  KybVerificationStatusType = "PENDING"
	KYBStatusApproved KybVerificationStatusType = "APPROVED"
	KYBStatusRejected KybVerificationStatusType = "REJECTED"

	DonationStatusPending  DonationStatus = "PENDING"
	DonationStatusApproved DonationStatus = "APPROVED"
	DonationStatusRejected DonationStatus = "REJECTED"
	DonationStatusReleased DonationStatus = "RELEASED"

	OrganizationStatusActive    OrganizationStatus = "ACTIVE"
	OrganizationStatusNotActive OrganizationStatus = "NOT_ACTIVE"
	OrganizationStatusPending   OrganizationStatus = "PENDING"

	Fiat   PaymentType = "FIAT"
	Crypto PaymentType = "CRYPTO"

	ProjectCategoryOpenInnovation  ProjectCategory = "OPEN_INNOVATION"
	ProjectCategoryWomenLeaders    ProjectCategory = "WOMEN_LEADERS"
	ProjectCategoryEmergingMarkets ProjectCategory = "EMERGING_MARKETS"
)

// ------------------------------------------------------

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
func (o *KybVerificationStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*o = KybVerificationStatusType(string(v))
	case string:
		*o = KybVerificationStatusType(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (o KybVerificationStatusType) Value() (driver.Value, error) {
	return string(o), nil
}

// ----------------------------------------------------------

func (o *OrganizationStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*o = OrganizationStatus(string(v))
	case string:
		*o = OrganizationStatus(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (o OrganizationStatus) Value() (driver.Value, error) {
	return string(o), nil
}

// ----------------------------------------------------------

func (p *PaymentType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*p = PaymentType(string(v))
	case string:
		*p = PaymentType(v)
	default:
		return fmt.Errorf("failed to scan operator type: %v", value)
	}
	return nil
}

func (p PaymentType) Value() (driver.Value, error) {
	return string(p), nil
}

func (pc *ProjectCategory) Scan(value interface{}) error {
	return scanEnum(value, (*string)(pc))
}

func (pc ProjectCategory) Value() (driver.Value, error) {
	return string(pc), nil
}

// ----------------------------------------------------------

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
