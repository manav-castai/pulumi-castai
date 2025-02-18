package main

import (
    _ "embed"

    "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"

    castaiProvider "github.com/manav-castai/pulumi-castai/provider"
    "github.com/manav-castai/pulumi-castai/provider/pkg/version"
)

//go:embed schema.json
var pulumiSchema []byte

func main() {
    tfbridge.Main("castai", version.Version, castaiProvider.Provider(), pulumiSchema)
}