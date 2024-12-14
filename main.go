package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		provider, err := createProvider(ctx)
		if err != nil {
			return err
		}

		err = createVPC(ctx, provider)
		if err != nil {
			return err
		}

		return nil
	})
}
