package ironic

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePortV1() *schema.Resource {
	return &schema.Resource{
		Description: "Ports may be specified as part of the node resource, or as a separate ironic_port_v1 declaration.",
		Create:      resourcePortV1Create,
		Read:        resourcePortV1Read,
		Update:      resourcePortV1Update,
		Delete:      resourcePortV1Delete,

		Schema: map[string]*schema.Schema{
			"node_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_group_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_link_connection": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"pxe_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"physical_network": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extra": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"is_smart_nic": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePortV1Create(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	opts := portSchemaToCreateOpts(d)
	result, err := ports.Create(client, opts).Extract()
	if err != nil {
		return err
	}
	d.SetId(result.UUID)

	return resourcePortV1Read(d, meta)
}

func resourcePortV1Read(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	port, err := ports.Get(client, d.Id()).Extract()
	if err != nil {
		return err
	}

	err = d.Set("address", port.Address)
	if err != nil {
		return err
	}
	err = d.Set("node_uuid", port.NodeUUID)
	if err != nil {
		return err
	}
	err = d.Set("port_group_uuid", port.PortGroupUUID)
	if err != nil {
		return err
	}
	err = d.Set("local_link_connection", port.LocalLinkConnection)
	if err != nil {
		return err
	}
	err = d.Set("pxe_enabled", port.PXEEnabled)
	if err != nil {
		return err
	}
	err = d.Set("physical_network", port.PhysicalNetwork)
	if err != nil {
		return err
	}
	err = d.Set("extra", port.Extra)
	if err != nil {
		return err
	}
	return d.Set("is_smart_nic", port.IsSmartNIC)
}

func resourcePortV1Update(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}
	d.Partial(true)

	stringFields := []string{
		"pxe_enabled",
		"address",
		"physical_network",
		"extra",
		"is_smart_nic",
		"node_uuid",
		"port_group_uuid",
	}

	localLinkConnectionFields := []string{
		"port_id",
		"switch_id",
		"switch_info",
	}

	for _, field := range stringFields {
		if d.HasChange(field) {
			opts := ports.UpdateOpts{
				ports.UpdateOperation{
					Op:    ports.ReplaceOp,
					Path:  fmt.Sprintf("/%s", field),
					Value: d.Get(field).(string),
				},
			}

			_, err := ports.Update(client, d.Id(), opts).Extract()
			if err != nil {
				return err
			}
		}
	}
	if d.HasChange("local_link_connection") {
		localLinkConnectionInfo := d.Get("local_link_connection").(map[string]interface{})
		for _, field := range localLinkConnectionFields {
			opts := ports.UpdateOpts{
				ports.UpdateOperation{
					Op:    ports.ReplaceOp,
					Path:  fmt.Sprintf("/local_link_connection/%s", field),
					Value: localLinkConnectionInfo[field],
				},
			}
			_, err := ports.Update(client, d.Id(), opts).Extract()
			if err != nil {
				return err
			}
		}

	}
	d.Partial(false)

	return resourcePortV1Read(d, meta)

}

func resourcePortV1Delete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}
	err = ports.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return err
	}

	return nil

}

func portSchemaToCreateOpts(d *schema.ResourceData) *ports.CreateOpts {
	pxeEnabled := d.Get("pxe_enabled").(bool)
	// isSmartNic := d.Get("is_smart_nic").(bool)

	opts := ports.CreateOpts{
		NodeUUID:            d.Get("node_uuid").(string),
		Address:             d.Get("address").(string),
		PortGroupUUID:       d.Get("port_group_uuid").(string),
		LocalLinkConnection: d.Get("local_link_connection").(map[string]interface{}),
		PXEEnabled:          &pxeEnabled,
		PhysicalNetwork:     d.Get("physical_network").(string),
		Extra:               d.Get("extra").(map[string]interface{}),
		// IsSmartNIC: &isSmartNic,
	}

	return &opts
}
