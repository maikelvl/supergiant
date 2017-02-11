package digitalocean_test

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
	"github.com/supergiant/supergiant/pkg/core"
	"github.com/supergiant/supergiant/pkg/model"
	"github.com/supergiant/supergiant/pkg/provider/digitalocean"
	"github.com/supergiant/supergiant/test/fake_core"
	"github.com/supergiant/supergiant/test/fake_digitalocean_provider"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDigitalOceanProviderValidateAccount(t *testing.T) {
	Convey("DigitalOcean Provider ValidateAccount works correctly", t, func() {
		table := []struct {
			// Input
			cloudAccount *model.CloudAccount
			// Mocks
			// Expectations
			err error
		}{
			// A successful example
			{
				// Input
				cloudAccount: &model.CloudAccount{},
			},
		}

		for _, item := range table {

			c := &core.Core{
				DB:  new(fake_core.DB),
				Log: logrus.New(),
			}

			provider := &digitalocean.Provider{
				Core: c,
				Client: func(kube *model.Kube) *godo.Client {
					return &godo.Client{
						Droplets: &fake_digitalocean_provider.Droplets{
							ListFn: func(_ *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
								return nil, nil, nil
							},
						},
					}
				},
			}

			err := provider.ValidateAccount(item.cloudAccount)

			So(err, ShouldEqual, item.err)
		}
	})
}
