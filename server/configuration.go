package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
	"gitlab.com/w1572/workchat-plugin-nagios/go-nagios/nagios"
)

// configuration captures the plugin's external configuration as exposed in the
// Workchat server configuration, as well as values computed from the
// configuration. Any public fields will be deserialized from the Workchat
// server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and
// the plugin configuration can change at any time, access to the configuration
// must be synchronized. The strategy used in this plugin is to guard a pointer
// to the configuration, and clone the entire struct whenever it changes. You
// may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to
// rewrite Clone as a deep copy appropriate for your types.
type configuration struct {
	NagiosURL              string
	Token                  string
	InitialLogsLimit       int
	InitialLogsStartTime   int
	InitialReportFrequency int
}

func (c *configuration) isValid() error {
	if len(c.NagiosURL) == 0 {
		return errors.New("the Nagios URL must be set")
	}

	if c.InitialLogsLimit <= 0 {
		return errors.New("initial logs limit must be a positive integer")
	}

	if c.InitialLogsStartTime <= 0 {
		return errors.New("initial logs start time must be a positive integer")
	}

	if c.InitialReportFrequency <= 0 {
		return errors.New("initial report frequency must be a positive integer")
	}

	return nil
}

// Clone shallow copies the configuration. Your implementation may require a
// deep copy if your configuration has reference types.
func (c *configuration) Clone() *configuration {
	var clone = *c
	return &clone
}

// getConfiguration retrieves the active configuration under lock, making it
// safe to use concurrently. The active configuration may change underneath the
// client of this method, but the struct returned by this API call is considered
// immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as
// sync.Mutex is not reentrant. In particular, avoid using the plugin API
// entirely, as this may in turn trigger a hook back into the plugin. If that
// hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing
// configuration. This almost certainly means that the configuration was
// modified without being cloned and may result in an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will
		// optimize the allocation for same to point at the same memory address,
		// breaking the check above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

// OnConfigurationChange is invoked when configuration changes may have been
// made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Workchat server
	// configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	p.setConfiguration(configuration)

	config := p.getConfiguration()

	if err := config.isValid(); err != nil {
		return err
	}

	c, err := nagios.NewClient(http.DefaultClient, config.NagiosURL)
	if err != nil {
		return fmt.Errorf("nagios.NewClient: %w", err)
	}

	p.client = c

	if err := p.storeInitialKV(); err != nil {
		return fmt.Errorf("p.storeInitialKV: %w", err)
	}

	p.API.LogInfo("Reloaded configuration")

	return nil
}
