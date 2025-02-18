// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package castai // <-- Renamed from package xyz to package castai

import (
	"path"

	// Allow embedding bridge-metadata.json in the provider.
	_ "embed"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge/tokens"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"

	// Instead of referencing the xyz provider, reference the CAST AI provider code:
	castaiTF "github.com/castai/terraform-provider-castai/castai"

	// Updated import path to your local versioning package:
	"github.com/manav-castai/pulumi-castai/provider/pkg/version"
)

// Constants changed from xyz to castai
const (
	mainPkg = "castai" // The new Pulumi package name
	mainMod = "index"  // The primary module name
)

//go:embed cmd/pulumi-resource-castai/bridge-metadata.json
var metadata []byte

// Provider returns additional overlaid schema and metadata associated with the provider.
func Provider() tfbridge.ProviderInfo {
	prov := tfbridge.ProviderInfo{
		// Use shimv2.NewProvider(...) to wrap the Terraform provider from CAST AI:
		// P: shimv2.NewProvider(castaiTF.New(version.Version)()),
		P: shimv2.NewProvider(castaiTF.Provider(version.Version)),
		
		Name:         "castai", // Pulumi provider name
		Version:      version.Version,
		DisplayName:  "",
		Publisher:    "Pulumi", // or "CAST AI," or your own name/org
		LogoURL:      "",
		Description:  "A Pulumi package for creating and managing CAST AI resources.",
		Keywords:     []string{"castai", "category/cloud"},
		License:      "Apache-2.0",
		Homepage:     "https://www.pulumi.com",
		Repository:   "https://github.com/manav-castai/pulumi-castai",
		GitHubOrg:    "", // Optional if needed

		MetadataInfo: tfbridge.NewProviderMetadata(metadata),

		// Example config. You can customize or remove if CAST AI doesn't require region config.
		Config: map[string]*tfbridge.SchemaInfo{
			"region": {
				Type: "castai:region/region:Region",
			},
		},

		// If you have extra types for config or data, rename them to castai as well:
		ExtraTypes: map[string]schema.ComplexTypeSpec{
			"castai:region/region:Region": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Type: "string",
				},
				Enum: []schema.EnumValueSpec{
					{Name: "here", Value: "HERE"},
					{Name: "overThere", Value: "OVER_THERE"},
				},
			},
		},

		JavaScript: &tfbridge.JavaScriptInfo{
			RespectSchemaVersion: true,
		},
		Python: &tfbridge.PythonInfo{
			RespectSchemaVersion: true,
			PyProject: struct{ Enabled bool }{true},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: path.Join(
				"github.com/manav-castai/pulumi-castai/sdk/",
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				mainPkg,
			),
			GenerateResourceContainerTypes: true,
			GenerateExtraInputTypes:        true,
			RespectSchemaVersion:           true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RespectSchemaVersion: true,
			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
	}

	// MustComputeTokens will map resources/datasources from the upstream TF provider to Pulumi tokens.
	prov.MustComputeTokens(tokens.SingleModule("castai_", mainMod, tokens.MakeStandard(mainPkg)))

	prov.MustApplyAutoAliases()
	prov.SetAutonaming(255, "-")

	return prov
}