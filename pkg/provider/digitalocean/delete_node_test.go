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

func TestDigitalOceanProviderDeleteNode(t *testing.T) {
    Convey("DigitalOcean Provider DeleteNode works correctly", t, func() {
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
                    ProviderID: "1",
                    Kube: &model.Kube{
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
                            // Delete
                            DeleteFn: func(_ int) (*godo.Response, error) {
                                return nil, nil
                            },
                        },
                    }
                },
            }

            action := &core.Action{Status: new(model.ActionStatus)}
            err := provider.DeleteNode(item.node, action)

            So(err, ShouldEqual, item.err)
        }
    })
}
