package cisco

import (
	"fmt"
	"regexp"

	textmanipulations "github.com/ivankuchin/cisco-config-mapper/internal/text-manipulations"
)

type Iface struct {
	Name    string
	Content []string
}

type Cisco struct {
	Before_ifaces []string
	Ifaces        []Iface
	After_ifaces  []string
}

func (iface *Iface) ReplaceName(from string, to string) error {
	re := regexp.MustCompile("\\b" + from + "\\b")
	if re.MatchString(iface.Name) {

		// actual iface name replacement / prepending / appending
		iface.Name = re.ReplaceAllString(iface.Name, to)
		return nil
	}

	err := fmt.Errorf("ERROR: interface %s not found", from)
	fmt.Printf("%v", err)

	return err
}

func (iface *Iface) PrependContent(str string) error {
	if str != "" {
		iface.Content = append([]string{str}, iface.Content...)
		return nil
	}

	return nil
}

func (iface *Iface) AppendContent(str string) error {
	if str != "" {
		iface.Content = append(iface.Content, str)
		return nil
	}

	return nil
}

func (iface *Iface) RemoveLines(lines []string) error {
	iface.Content = textmanipulations.Remove_lines(iface.Content, lines)
	return nil
}

func (iface *Iface) RemovePrefixes(prefixes []string) error {
	iface.Content = textmanipulations.Remove_prefixes(iface.Content, prefixes)
	return nil
}
