package config

import (
	"fmt"
	"path"

	"github.com/spf13/viper"

	sdkioerrors "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/server/config"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// DefaultGRPCAddress is the default address the gRPC server binds to.
	DefaultGRPCAddress = "0.0.0.0:9900"
)

// Config defines the server's top level configuration. It includes the default app config
// from the SDK as well as the TLS configuration.
type Config struct {
	config.Config
	TLS     TLSConfig     `mapstructure:"tls"`
	Supply  SupplyConfig  `mapstructure:"supply"`
	Metrics MetricsConfig `mapstructure:"metrics"`
}

// TLSConfig defines the certificate and matching private key for the server.
type TLSConfig struct {
	// CertificatePath the file path for the certificate .pem file
	CertificatePath string `mapstructure:"certificate-path"`
	// KeyPath the file path for the key .pem file
	KeyPath string `mapstructure:"key-path"`
}

type SupplyConfig struct {
	// Enabled determines whether we expose the supply endpoint.
	Enabled bool `mapstructure:"enabled"`
}

type MetricsConfig struct {
	URL string `mapstructure:"url"`
}

// AppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func AppConfig(denom string) (customAppTemplate string, customAppConfig interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := config.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// By default, we set the min gas prices to 0.
	if denom != "" {
		srvCfg.MinGasPrices = "0" + denom
	}

	customAppConfig = Config{
		Config:  *srvCfg,
		TLS:     *DefaultTLSConfig(),
		Supply:  *DefaultSupplyConfig(),
		Metrics: *DefaultMetricsConfig(),
	}

	customAppTemplate = config.DefaultConfigTemplate + DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config:  *config.DefaultConfig(),
		TLS:     *DefaultTLSConfig(),
		Supply:  *DefaultSupplyConfig(),
		Metrics: *DefaultMetricsConfig(),
	}
}

// DefaultConfig returns server's default configuration.
func DefaultSupplyConfig() *SupplyConfig {
	return &SupplyConfig{
		Enabled: false,
	}
}

func DefaultMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		URL: "",
	}
}

// DefaultTLSConfig returns the default TLS configuration.
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		CertificatePath: "",
		KeyPath:         "",
	}
}

// Validate returns an error if the TLS certificate and key file extensions are invalid.
func (c TLSConfig) Validate() error {
	certExt := path.Ext(c.CertificatePath)

	if c.CertificatePath != "" && certExt != ".pem" {
		return fmt.Errorf("invalid extension %s for certificate path %s, expected '.pem'", certExt, c.CertificatePath)
	}

	keyExt := path.Ext(c.KeyPath)

	if c.KeyPath != "" && keyExt != ".pem" {
		return fmt.Errorf("invalid extension %s for key path %s, expected '.pem'", keyExt, c.KeyPath)
	}

	return nil
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) Config {
	cfg, _ := config.GetConfig(v)

	return Config{
		Config: cfg,
		TLS: TLSConfig{
			CertificatePath: v.GetString("tls.certificate-path"),
			KeyPath:         v.GetString("tls.key-path"),
		},
		Supply: SupplyConfig{
			Enabled: v.GetBool("supply.enabled"),
		},
		Metrics: MetricsConfig{
			URL: v.GetString("metrics.url"),
		},
	}
}

// ParseConfig retrieves the default environment configuration for the
// application.
func ParseConfig(v *viper.Viper) (*Config, error) {
	conf := DefaultConfig()
	err := v.Unmarshal(conf)

	return conf, err
}

// ValidateBasic returns an error any of the application configuration fields are invalid.
func (c Config) ValidateBasic() error {
	if err := c.TLS.Validate(); err != nil {
		return sdkioerrors.Wrapf(sdkerrors.ErrAppConfig, "invalid tls config value: %s", err.Error())
	}

	return c.Config.ValidateBasic()
}
