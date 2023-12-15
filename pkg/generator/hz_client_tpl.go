/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import "github.com/cloudwego/cwgo/pkg/consts"

// related to service resolver
var (
	hzNilResolverFuncBody = "{\n\treturn\n}"

	hzAppendResolverFunc = `func GetResolverAddress() []string {
		e := os.Getenv("GO_HERTZ_RESOLVER_[[ToUpper .ServiceName]]")
		if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
	    return strings.Fields(e)
      }`

	hzCommonResolverImport = "github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"

	hzCommonResolverBody = `*ops = append(*ops, WithHertzClientMiddleware(sd.Discovery(r)))
	return nil`

	hzEtcdClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/etcd",
	}

	hzEtcdClient = `r, err := etcd.NewEtcdResolver(http.GetResolverAddress())
	if err != nil {
		return err
	}` + consts.LineBreak + hzCommonResolverBody

	hzNacosClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/nacos",
	}

	hzNacosClient = `r, err := nacos.NewDefaultNacosResolver()
	if err != nil {
		return err
	}` + consts.LineBreak + hzCommonResolverBody

	hzConsulClientImports = []string{
		hzCommonResolverImport,
		"github.com/hashicorp/consul/api",
		"github.com/hertz-contrib/registry/consul",
	}

	hzConsulClient = `consulConfig := api.DefaultConfig()
    consulConfig.Address = http.GetResolverAddress()[0]
    consulClient, err := api.NewClient(consulConfig)
    if err != nil {
        return err
    }
    
    r := consul.NewConsulResolver(consulClient)` + consts.LineBreak + hzCommonResolverBody

	hzEurekaClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/eureka",
	}

	hzEurekaClient = `r := eureka.NewEurekaResolver(http.GetResolverAddress())` +
		consts.LineBreak + hzCommonResolverBody

	hzPolarisClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/polaris",
	}

	hzPolarisClient = `r, err := polaris.NewPolarisResolver()
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonResolverBody

	hzServiceCombClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/servicecomb",
	}

	hzServiceCombClient = `r, err := servicecomb.NewDefaultSCResolver(http.GetResolverAddress())
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonResolverBody

	hzZKClientImports = []string{
		hzCommonResolverImport,
		"github.com/hertz-contrib/registry/zookeeper",
		"time",
	}

	hzZKClient = `r, err := zookeeper.NewZookeeperResolver(http.GetResolverAddress(), 40*time.Second)
    if err != nil {
        return err
    }` + consts.LineBreak + hzCommonResolverBody
)

var hzClientMVCTemplates = []Template{
	{
		Path:   consts.InitGo,
		Delims: [2]string{consts.LeftDelimiter, consts.RightDelimiter},
		UpdateBehavior: UpdateBehavior{
			AppendRender: map[string]interface{}{},
			ReplaceFunc: ReplaceFunc{
				ReplaceFuncName:   make([]string, 0, 5),
				ReplaceFuncImport: make([][]string, 0, 15),
				ReplaceFuncBody:   make([]string, 0, 5),
			},
		},
		Body: `package {{.InitOptsPackage}}
      import (
		{{range $key, $value := .GoFileImports}}
	    {{if eq $key "init.go"}}
	    {{range $k, $v := $value}}
        {{if ne $k ""}}"{{$k}}"{{end}}{{end}}{{end}}{{end}}
	  )

	  func initClientOpts(hostUrl string) (ops []Option, err error) {
		ops = append(ops, withHostUrl(hostUrl))
		
		if err = initResolver(&ops); err != nil {
		  panic(err)
		}

		return
	  }
	  
	  // If you do not use the service resolver function, do not edit this function.
	  // Otherwise, you can customize and modify it.
	  func initResolver(ops *[]Option) (err error) {
		{{if ne .ResolverName ""}}
		{{.ResolverBody}}
		{{else}}
		return
        {{end}}
	  }`,
	},

	{
		Path:   consts.DefaultHZClientDir + consts.Slash + consts.EnvGo,
		Delims: [2]string{"[[", "]]"},
		UpdateBehavior: UpdateBehavior{
			AppendRender: map[string]interface{}{},
		},
		CustomFunc: TemplateCustomFuncMap,
		Body: `// Code generated by cwgo generator. DO NOT EDIT.

	  package http
	  import (
		[[range $key, $value := .GoFileImports]]
	    [[if eq $key "env.go"]]
	    [[range $k, $v := $value]]
        [[if ne $k ""]]"[[$k]]"[[end]][[end]][[end]][[end]]
	  )

      [[if ne .ResolverName ""]]
      func GetResolverAddress() []string {
		e := os.Getenv("GO_HERTZ_RESOLVER_[[ToUpper .ServiceName]]")
	    if len(e) == 0 {
		  return []string{[[$lenSlice := len .ResolverAddress]][[range $key, $value := .ResolverAddress]]"[[$value]]"[[if eq $key (Sub $lenSlice 1)]][[else]], [[end]][[end]]}
	    }
		return strings.Fields(e)
      }
	  [[end]]`,
	},
}
