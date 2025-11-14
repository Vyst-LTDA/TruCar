package schemas

type OrganizationUpdate struct {
	Name   string `json:"name"`
	Sector string `json:"sector"`
}

type OrganizationFuelIntegrationUpdate struct {
	FuelProviderName string  `json:"fuel_provider_name" binding:"required"`
	APIKey           *string `json:"api_key"`
	APISecret        *string `json:"api_secret"`
}

type OrganizationFuelIntegrationPublic struct {
	FuelProviderName string `json:"fuel_provider_name"`
	IsAPIKeySet      bool   `json:"is_api_key_set"`
	IsAPISecretSet   bool   `json:"is_api_secret_set"`
}
