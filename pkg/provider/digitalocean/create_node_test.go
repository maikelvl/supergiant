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

func TestDigitalOceanProviderCreateNode(t *testing.T) {
    Convey("DigitalOcean Provider CreateNode works correctly", t, func() {
        table := []struct {
            // Input
            node *model.Node
            // Mocks
            // Expectations
            err error
        }{
            // A successful example
            {
                // Input
                node: &model.Node{
                    Kube: &model.Kube{
                        CloudAccount: &model.CloudAccount{
                            Credentials: map[string]string{"token": "my-special-token"},
                        },
                        DigitalOceanConfig: &model.DOKubeConfig{},
                    },
                },
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
                            // Create
                            CreateFn: func(_ *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
                                return &godo.Droplet{
                                    ID: 1,
                                }, nil, nil
                            },
                            // Get
                            GetFn: func(int) (*godo.Droplet, *godo.Response, error) {
                                return &godo.Droplet{
                                    ID: 1,
                                    Networks: &godo.Networks{
                                        V4: []godo.NetworkV4{
                                            {
                                                Type:      "public",
                                                IPAddress: "99.99.99.99",
                                            },
                                            {
                                                Type:      "private",
                                                IPAddress: "10.0.0.99",
                                            },
                                        },
                                    },
                                }, nil, nil
                            },
                        },
                        Tags: &fake_digitalocean_provider.Tags{},
                    }
                },
            }

            action := &core.Action{Status: new(model.ActionStatus)}
            err := provider.CreateNode(item.node, action)

            So(err, ShouldEqual, item.err)
        }
    })
}
