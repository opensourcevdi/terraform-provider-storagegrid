package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	metadata := provider.MetadataResponse{}
	(&Provider{}).Metadata(context.Background(), provider.MetadataRequest{}, &metadata)
	providerserver.Serve(context.Background(),
		func() provider.Provider { return &Provider{} },
		providerserver.ServeOpts{Address: "invalid.invalid/invalid/" + metadata.TypeName})
}
