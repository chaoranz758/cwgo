/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package api

import (
	"fmt"
	"github.com/bytedance/gopkg/util/gopool"
	registryconfig "github.com/cloudwego/cwgo/platform/server/shared/config/internal/registry"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/registry"
	"github.com/cloudwego/hertz/pkg/app/server"
	http2config "github.com/hertz-contrib/http2/config"
	http2factory "github.com/hertz-contrib/http2/factory"
	"github.com/hertz-contrib/pprof"
	"time"
)

type ConfigManager struct {
	config                Config
	RegistryConfigManager registryconfig.IRegistryConfigManager
	BuiltinRegistry       *registry.BuiltinRegistry
	Server                *server.Hertz
	ServiceId             string
	ServiceName           string
}

func NewConfigManager(config Config, registryConfig registryconfig.Config, serviceId string) *ConfigManager {
	var registryConfigManager registryconfig.IRegistryConfigManager
	var r *registry.BuiltinRegistry
	var err error

	switch registryConfig.Type {
	case consts.RegistryTypeBuiltin:
		r = registry.NewBuiltinRegistry()
		registryConfigManager, err = registryconfig.NewBuiltinRegistryConfigManager(registryConfig.Builtin, r)

	case consts.RegistryTypeConsul:
		registryConfigManager, err = registryconfig.NewConsulRegistryConfigManager(registryConfig.Consul)
		if err != nil {
			panic(fmt.Sprintf("initialize registry failed, err: %v", err))
		}

	default:

	}

	hertzServer := server.New(
		server.WithHostPorts(fmt.Sprintf("%s:%d", config.Host, config.Port)),
		server.WithKeepAliveTimeout(1*time.Minute),
		server.WithReadTimeout(3*time.Minute),
		server.WithIdleTimeout(3*time.Minute),
		server.WithMaxRequestBodySize(1<<20*4), // 4M
		server.WithHandleMethodNotAllowed(true),
		server.WithExitWaitTime(5*time.Second),
		server.WithBasePath("/api"),
		server.WithMaxKeepBodySize(1<<20*4),
		server.WithKeepAlive(true),
		server.WithH2C(true),
		server.WithReadBufferSize(1<<10*4),
	)

	gopool.SetCap(10000) // max connections

	// user http2
	hertzServer.AddProtocol("h2",
		http2factory.NewServerFactory(
			http2config.WithReadTimeout(1*time.Minute),
			http2config.WithDisableKeepAlive(false),
		),
	)

	// register pprof service
	pprof.Register(hertzServer)

	return &ConfigManager{
		config:                config,
		Server:                hertzServer,
		RegistryConfigManager: registryConfigManager,
		BuiltinRegistry:       r,
		ServiceId:             serviceId,
		ServiceName:           fmt.Sprintf("%s-%s-%s", "cwgo", consts.ServerTypeAgent, serviceId),
	}
}
