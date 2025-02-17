// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.
//go:build linux
// +build linux

package modules

import (
	"fmt"

	"github.com/DataDog/datadog-agent/cmd/system-probe/api/module"
	"github.com/DataDog/datadog-agent/cmd/system-probe/config"
	"github.com/DataDog/datadog-agent/pkg/ebpf"
	sconfig "github.com/DataDog/datadog-agent/pkg/security/config"
	secmodule "github.com/DataDog/datadog-agent/pkg/security/module"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	// DefaultRuntimePoliciesDir is the default policies directory used by the runtime security module
	DefaultRuntimePoliciesDir = "/etc/datadog-agent/runtime-security.d"
)

// SecurityRuntime - Security runtime Factory
var SecurityRuntime = module.Factory{
	Name:             config.SecurityRuntimeModule,
	ConfigNamespaces: []string{"runtime_security_config"},
	Fn: func(agentConfig *config.Config) (module.Module, error) {
		config, err := sconfig.NewConfig(agentConfig)
		if err != nil {
			return nil, fmt.Errorf("invalid security runtime module configuration: %w", err)
		}

		m, err := secmodule.NewModule(config)
		if err == ebpf.ErrNotImplemented {
			log.Info("Datadog runtime security agent is only supported on Linux")
			return nil, module.ErrNotEnabled
		}
		return m, err
	},
}
