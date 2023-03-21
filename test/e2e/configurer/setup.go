package configurer

import (
	"fmt"
)

type setupFn func(configurer Configurer) error

func baseSetup(configurer Configurer) error {
	return configurer.RunValidators()
}

func withIBC(setupHandler setupFn) setupFn {
	return func(configurer Configurer) error {
		if err := setupHandler(configurer); err != nil {
			return err
		}

		return configurer.RunIBC()
	}
}

func withUpgrade(setupHandler setupFn) setupFn {
	return func(configurer Configurer) error {
		if err := setupHandler(configurer); err != nil {
			return err
		}

		upgradeConfigurer, ok := configurer.(*UpgradeConfigurer)
		if !ok {
			return fmt.Errorf("to run with upgrade, %v must be set during initialization", &UpgradeConfigurer{})
		}

		if err := upgradeConfigurer.CreatePreUpgradeState(); err != nil {
			return err
		}

		return upgradeConfigurer.RunUpgrade()
	}
}
