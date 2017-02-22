package digitalocean

import (
    "bytes"
    "strings"
    "text/template"
    "time"

    "github.com/digitalocean/godo"
    "github.com/supergiant/supergiant/bindata"
    "github.com/supergiant/supergiant/pkg/core"
    "github.com/supergiant/supergiant/pkg/model"
    "github.com/supergiant/supergiant/pkg/util"
)

// CreateKube creates a new DO kubernetes cluster.
func (p *Provider) CreateKube(m *model.Kube, action *core.Action) error {
    procedure := &core.Procedure{
        Core:   p.Core,
        Name:   "Create Kube",
        Model:  m,
        Action: action,
    }

    client := p.Client(m)

    procedure.AddStep("creating global tags for Kube", func() error {
        // These are created once, and then attached by name to created resource
        globalTags := []string{
            "Kubernetes-Cluster",
            m.Name,
            m.Name + "-master",
            m.Name + "-minion",
        }
        for _, tag := range globalTags {
            createInput := &godo.TagCreateRequest{
                Name: tag,
            }
            if _, _, err := client.Tags.Create(createInput); err != nil {
                // TODO
                p.Core.Log.Warnf("Failed to create Digital Ocean tag '%s': %s", tag, err)
            }
        }
        return nil
    })

    procedure.AddStep("creating master", func() error {
        if m.MasterPublicIP != "" {
            return nil
        }

        // Build template
        masterUserdataTemplate, err := bindata.Asset("config/providers/digitalocean/master.yaml")
        if err != nil {
            return err
        }
        masterTemplate, err := template.New("master_template").Parse(string(masterUserdataTemplate))
        if err != nil {
            return err
        }
        var masterUserdata bytes.Buffer
        if err = masterTemplate.Execute(&masterUserdata, m); err != nil {
            return err
        }

        dropletRequest := &godo.DropletCreateRequest{
            Name:              m.Name + "-master-" + strings.ToLower(util.RandomString(5)),
            Region:            m.DigitalOceanConfig.Region,
            Size:              m.MasterNodeSize,
            PrivateNetworking: true,
            UserData:          string(masterUserdata.Bytes()),
            SSHKeys: []godo.DropletCreateSSHKey{
                {
                    Fingerprint: m.DigitalOceanConfig.SSHKeyFingerprint,
                },
            },
            Image: godo.DropletCreateImage{
                Slug: "coreos-stable",
            },
        }
        tags := []string{"Kubernetes-Cluster", m.Name, dropletRequest.Name}

        masterDroplet, publicIP, privateIP, err := p.createDroplet(client, action, dropletRequest, tags)
        if err != nil {
            return err
        }

        m.DigitalOceanConfig.MasterID = masterDroplet.ID
        m.MasterPublicIP = publicIP
        m.MasterPrivateIP = privateIP
        return nil
    })

    procedure.AddStep("building Kubernetes minion", func() error {
        // Load Nodes to see if we've already created a minion
        // TODO -- I think we can get rid of a lot of this do-unless behavior if we
        // modify Procedure to save progess on Action (which is easy to implement).
        if err := p.Core.DB.Find(&m.Nodes, "kube_name = ?", m.Name); err != nil {
            return err
        }
        if len(m.Nodes) > 0 {
            return nil
        }

        node := &model.Node{
            KubeName: m.Name,
            Kube:     m,
            Size:     m.NodeSizes[0],
        }
        return p.Core.Nodes.Create(node)
    })

    // TODO repeated in provider_aws.go
    procedure.AddStep("waiting for Kubernetes", func() error {
        return action.CancellableWaitFor("Kubernetes API and first minion", 20*time.Minute, 3*time.Second, func() (bool, error) {
            k8s := p.Core.K8S(m)
            k8sNodes, err := k8s.ListNodes("")
            if err != nil {
                return false, nil
            }
            return len(k8sNodes) > 0, nil
        })
    })

    return procedure.Run()
}
