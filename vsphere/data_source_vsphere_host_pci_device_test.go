// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vsphere

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/testhelper"
)

func TestAccDataSourceVSphereHostPciDevice_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccDataSourceVSphereHostPciDevicePreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereHostPciDeviceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(
						"data.vsphere_host_pci_device.device",
						"pci_devices.#",
						func(value string) error {
							valueInt, err := strconv.Atoi(value)
							if err != nil {
								return err
							}

							if valueInt <= 0 {
								return fmt.Errorf("number of PCI devices should be greater than 0")
							}
							return nil
						},
					),
				),
			},
			{
				Config: testAccDataSourceVSphereHostPciDeviceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.vsphere_host_pci_device.device",
						"pci_devices.0.name",
						regexp.MustCompile("(.*?)"),
					),
				),
			},
		},
	})
}

func testAccDataSourceVSphereHostPciDevicePreCheck(t *testing.T) {
	if os.Getenv("TF_VAR_VSPHERE_DATACENTER") == "" {
		t.Skip("set TF_VAR_VSPHERE_DATACENTER to run vsphere_host_pci_device acceptance tests")
	}
	if os.Getenv("TF_VAR_VSPHERE_ESXI1") == "" {
		t.Skip("set TF_VAR_VSPHERE_ESXI1 to run vsphere_host_pci_device acceptance tests")
	}
}

func testAccDataSourceVSphereHostPciDeviceConfig() string {
	return fmt.Sprintf(`
%s

data "vsphere_host" "host" {
  name          = "%s"
  datacenter_id = "${data.vsphere_datacenter.rootdc1.id}"
}

data "vsphere_host_pci_device" "device" {
  host_id    = "${data.vsphere_host.host.id}"
  name_regex = ""
}
`, testhelper.CombineConfigs(testhelper.ConfigDataRootDC1(), testhelper.ConfigDataRootPortGroup1()), os.Getenv("TF_VAR_VSPHERE_ESXI1"))
}
