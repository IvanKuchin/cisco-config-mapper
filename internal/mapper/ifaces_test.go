package mapper

import (
	"fmt"
	"testing"

	"github.com/ivankuchin/cisco-config-mapper/internal/cisco"
	"github.com/ivankuchin/cisco-config-mapper/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestMapIfaces(t *testing.T) {
	tests := []struct {
		name           string
		ifaces         []cisco.Iface
		ifaceMappings  []config.Iface_mapping
		expectedResult []cisco.Iface
		expectedError  error
	}{
		{
			name: "Simple mapping",
			ifaces: []cisco.Iface{
				{
					Name: "GigabitEthernet0/0",
					Content: []string{
						" description uplink",
						" no shutdown",
					},
				},
			},
			ifaceMappings: []config.Iface_mapping{
				{
					From:            "GigabitEthernet0/0",
					To:              "GigabitEthernet1/0",
					Remove_lines:    []string{" no shutdown"},
					Remove_prefixes: []string{" desc"},
					Prepend:         " ip address 192.168.168.1/24",
					Append:          " ip hsrp version 2",
				},
			},
			expectedResult: []cisco.Iface{
				{
					Name: "GigabitEthernet1/0",
					Content: []string{
						" ip address 192.168.168.1/24",
						" ip hsrp version 2",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Subiface mapping",
			ifaces: []cisco.Iface{
				{
					Name: "GigabitEthernet0/0.100",
					Content: []string{
						" description uplink",
						" no shutdown",
					},
				},
			},
			ifaceMappings: []config.Iface_mapping{
				{
					From:    "GigabitEthernet0/0",
					To:      "GigabitEthernet1/0",
					Prepend: " ip address 192.168.168.1/24",
					Append:  " ip hsrp version 2",
				},
			},
			expectedResult: []cisco.Iface{
				{
					Name: "GigabitEthernet1/0.100",
					Content: []string{
						" ip address 192.168.168.1/24",
						" description uplink",
						" no shutdown",
						" ip hsrp version 2",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Mapping not found",
			ifaces: []cisco.Iface{
				{Name: "GigabitEthernet0/1"},
			},
			ifaceMappings: []config.Iface_mapping{
				{
					From: "GigabitEthernet0/0",
					To:   "GigabitEthernet1/0",
				},
			},
			expectedResult: []cisco.Iface{},
			expectedError:  fmt.Errorf("ERROR: required interfaces not found"),
		},
		{
			name: "Error in ReplaceName",
			ifaces: []cisco.Iface{
				{Name: "GigabitEthernet1/1"},
			},
			ifaceMappings: []config.Iface_mapping{
				{
					From: "GigabitEthernet0/0",
					To:   "InvalidName",
				},
			},
			expectedResult: []cisco.Iface{},
			expectedError:  fmt.Errorf("ERROR: required interfaces not found"),
			// expectedError: &ErrorMappingNotFound{
			// 	iface_name: "GigabitEthernet1/1",
			// },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := map_ifaces(tt.ifaces, tt.ifaceMappings)
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
