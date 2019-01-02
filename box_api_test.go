package vagrant_go

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoxApiList(t *testing.T) {
	t.Parallel()

	t.Run(
		"with no vagrant boxes available, it returns empty slice",
		func(t *testing.T) {
			boxAPI := &boxAPI{
				client: emptyTestClient(t),
			}

			boxes, err := boxAPI.List()

			require.NoError(t, err)
			assert.Empty(t, boxes)
		},
	)

	t.Run(
		"with 1 vagrant box available, it returns slice of 1 box",
		func(t *testing.T) {
			boxCommandRunFunc := func(cmd string, args ...string) ([]byte, error) {
				output := `
1546015529,,ui,info,my-debian (libvirt%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-debian
1546015529,,box-provider,libvirt
1546015529,,box-version,1.2.3
`
				return []byte(output), nil
			}

			boxAPI := &boxAPI{
				client: testClient(
					t,
					boxCommandRunFunc,
					emptyLookPathFunc,
				),
			}

			boxes, err := boxAPI.List()

			require.NoError(t, err)
			assert.Len(t, boxes, 1)

			assert.Equal(t, "my-debian", boxes[0].Name)
			assert.Equal(t, "libvirt", boxes[0].Provider)
			assert.Equal(t, "1.2.3", boxes[0].Version)
		},
	)

	t.Run(
		"with 3 vagrant boxes available, it returns slice of 3 boxes",
		func(t *testing.T) {
			boxCommandRunFunc := func(cmd string, args ...string) ([]byte, error) {
				output := `
1546015529,,ui,info,my-debian (libvirt%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-debian
1546015529,,box-provider,libvirt
1546015529,,box-version,1.2.3
1546015529,,ui,info,my-vbox-debian (virtualbox%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-vbox-debian
1546015529,,box-provider,virtualbox
1546015529,,box-version,1.2.4
1546015529,,ui,info,my-vmware-debian (vmware%!(VAGRANT_COMMA) 0)
1546015529,,box-name,my-vmware-debian
1546015529,,box-provider,vmware
1546015529,,box-version,1.2.5
`
				return []byte(output), nil
			}

			boxAPI := &boxAPI{
				client: testClient(
					t,
					boxCommandRunFunc,
					emptyLookPathFunc,
				),
			}

			boxes, err := boxAPI.List()

			require.NoError(t, err)
			assert.Len(t, boxes, 3)

			assert.Equal(t, "my-debian", boxes[0].Name)
			assert.Equal(t, "libvirt", boxes[0].Provider)
			assert.Equal(t, "1.2.3", boxes[0].Version)

			assert.Equal(t, "my-vbox-debian", boxes[1].Name)
			assert.Equal(t, "virtualbox", boxes[1].Provider)
			assert.Equal(t, "1.2.4", boxes[1].Version)

			assert.Equal(t, "my-vmware-debian", boxes[2].Name)
			assert.Equal(t, "vmware", boxes[2].Provider)
			assert.Equal(t, "1.2.5", boxes[2].Version)
		},
	)
}
