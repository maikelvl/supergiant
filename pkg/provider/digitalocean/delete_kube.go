package digitalocean

import (
    "strings"

    "github.com/supergiant/supergiant/pkg/core"
    "github.com/supergiant/supergiant/pkg/model"
)

// DeleteKube deletes a DO kubernetes cluster.
func (p *Provider) DeleteKube(m *model.Kube, action *core.Action) error {
    // New Client
    client := p.Client(m)
    // Step procedure
    procedure := &core.Procedure{
        Core:   p.Core,
        Name:   "Delete Kube",
        Model:  m,
        Action: action,
    }

    procedure.AddStep("deleting master", func() error {
        if m.DigitalOceanConfig.MasterID == 0 {
            return nil
        }
        if _, err := client.Droplets.Delete(m.DigitalOceanConfig.MasterID); err != nil && !strings.Contains(err.Error(), "404") {
            return err
        }
        m.DigitalOceanConfig.MasterID = 0
        return nil
    })

    return procedure.Run()
}
