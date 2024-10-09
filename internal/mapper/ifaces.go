package mapper

import (
	"fmt"
	"regexp"

	"github.com/ivankuchin/cisco-config-mapper/internal/cisco"
	"github.com/ivankuchin/cisco-config-mapper/internal/config"
)

func map_ifaces(ifaces []cisco.Iface, iface_mappings []config.Iface_mapping) ([]cisco.Iface, error) {
	var err error

	result := make([]cisco.Iface, 0, len(ifaces))

	for _, iface := range ifaces {
		dst_iface_mapping, err := get_mapping_by_from_iface(iface, iface_mappings)
		if err != nil {
			_, ok := err.(*ErrorMappingNotFound)
			if ok {
				continue
			}
			return result, err
		}

		err = iface.ReplaceName(dst_iface_mapping.From, dst_iface_mapping.To)
		if err != nil {
			return result, err
		}

		err = iface.RemoveLines(dst_iface_mapping.Remove_lines)
		if err != nil {
			return result, err
		}

		err = iface.RemovePrefixes(dst_iface_mapping.Remove_prefixes)
		if err != nil {
			return result, err
		}

		err = iface.PrependContent(dst_iface_mapping.Prepend)
		if err != nil {
			return result, err
		}

		err = iface.AppendContent(dst_iface_mapping.Append)
		if err != nil {
			return result, err
		}

		dst_iface_mapping.MarkVisited()
		result = append(result, iface)
	}

	for _, iface_mapping := range iface_mappings {
		if !iface_mapping.Visited {
			fmt.Println("ERROR: interface " + iface_mapping.From + " not found")
			err = fmt.Errorf("ERROR: required interfaces not found")
		}
	}

	return result, err
}

func get_mapping_by_from_iface(iface cisco.Iface, iface_mappings []config.Iface_mapping) (*config.Iface_mapping, error) {
	for j := range iface_mappings {
		re := regexp.MustCompile("\\b" + iface_mappings[j].From + "\\b")
		if re.MatchString(iface.Name) {
			return &iface_mappings[j], nil
		}
	}

	// This is not an actual error
	// code will hit this place every time when interface in src doesn't have a mapping
	// i.e. when interface is not required to be converted

	// I will return an error, just to show that no further processing required for this interface

	// fmt.Println("ERROR: interface mapping 'From' in config-file for " + iface.Name + " not found")
	return nil, &ErrorMappingNotFound{
		iface_name: iface.Name,
	}
}
