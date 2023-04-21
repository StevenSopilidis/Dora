package vault

import (
	"context"
	"testing"

	"github.com/stevensopilidis/dora/registry"
	m "github.com/stevensopilidis/dora/vault/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setUpTest(t *testing.T) (
	vault *Vault,
	teardown func(),
) {
	t.Helper()
	vault = InitializeVault(&InitializeVaultConfig{
		Host: "localhost",
		Port: 5432,
		Db:   "doradb",
		User: "dora",
		Pass: "dora123",
	})
	teardown = func() {
		CloseVault(vault)
	}
	return vault, teardown
}

func TestVault(t *testing.T) {
	v, teardown := setUpTest(t)
	defer teardown()
	for scenario, fn := range map[string]func(
		t *testing.T,
		v *Vault,
	){
		"append entry to vault":           testAppendToVault,
		"get all entries from vault":      testGetFromVault,
		"remove entrie from vault":        testRemoveEntryFromVault,
		"remove invalid entry from vault": testRemoveInvalidEntryFromVault,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, v)
		})
	}
}

func testAppendToVault(t *testing.T, v *Vault) {
	service := &m.ServiceModel{
		Service: registry.Service{
			Addr: "127.0.0.1",
			Port: 80,
		},
	}
	err := v.AddService(context.Background(), service)
	require.NoError(t, err)
}

func testGetFromVault(t *testing.T, v *Vault) {
	//clear database
	result := v.db.Exec("DELETE FROM service_models WHERE 1=1")
	require.NoError(t, result.Error)
	service := &m.ServiceModel{
		Service: registry.Service{
			Addr: "127.0.0.1",
			Port: 80,
		},
	}
	err := v.AddService(context.Background(), service)
	require.NoError(t, err)
	service = &m.ServiceModel{
		Service: registry.Service{
			Addr: "127.5.0.1",
			Port: 60,
		},
	}
	err = v.AddService(context.Background(), service)
	require.NoError(t, err)
	err, services := v.GetServices(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, len(services))
}

func testRemoveEntryFromVault(t *testing.T, v *Vault) {
	service := &m.ServiceModel{
		Service: registry.Service{
			Addr: "127.0.0.1",
			Port: 80,
		},
	}
	err := v.AddService(context.Background(), service)
	require.NoError(t, err)
	err = v.RemoveService(context.Background(), service.ID)
	require.NoError(t, err)
}

func testRemoveInvalidEntryFromVault(t *testing.T, v *Vault) {
	service := &m.ServiceModel{
		Service: registry.Service{
			Addr: "127.0.0.1",
			Port: 80,
		},
	}
	err := v.RemoveService(context.Background(), service.ID)
	require.Equal(t, err, gorm.ErrRecordNotFound)
}
