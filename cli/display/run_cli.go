package display

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/spf13/cobra"
)

// user info.
func ShowUsr(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch users\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	// usrName := args[0]

	manyUsr := map[string]*auth.User{}

	if len(args) == 0 {
		if _, e := client.Get("users", nil, manyUsr); e != nil {
			return e
		}
	} else {
		path := fmt.Sprintf("users/%s", args[0])

		usr := &auth.User{} //nolint:exhaustruct,exhaustivestruct

		if _, e := client.Get(path, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	printUser(manyUsr)

	return nil
}

func ShowReg(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch registrations\n")

	/*
		usrName := args[0]
		manyUsr := map[string]*auth.User{}

		usr := &auth.User{} //nolint:exhaustruct,exhaustivestruct
		manyUsr[usrName] = usr

		printUser(manyUsr)

		return nil*/

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	// usrName := args[0]

	manyUsr := map[string]*auth.User{}

	if len(args) == 0 {
		if _, e := client.Get("registrations", nil, manyUsr); e != nil {
			return e
		}
	} else {
		path := fmt.Sprintf("registrations/%s", args[0])

		usr := &auth.User{} //nolint:exhaustruct,exhaustivestruct

		if _, e := client.Get(path, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	printUser(manyUsr)

	return nil
}

// network resources.
func ShowVlan(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch vlan\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	vlans := map[string]*network.VLANPool{}

	if len(args) == 0 {
		if _, e := config.Get("registrations", nil, vlans); e != nil {
			return e
		}
	} else {
		netName := args[0]

		vlan := &network.VLANPool{} //nolint:exhaustruct,exhaustivestruct
		vlans[netName] = vlan
	}

	printNets(vlans)

	return nil
}

func ShowSw(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch switches\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	swName := args[0]
	sw := &network.Switch{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, sw); e != nil {
		return e
	}

	manySw := map[string]*network.Switch{}

	manySw[swName] = sw

	printSwitch(manySw)

	return nil
}

func ShowIP(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch vlan\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	IPName := args[0]
	addr := &network.IPAddressPool{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, addr); e != nil {
		return e
	}

	pools := map[string]*network.IPAddressPool{}

	pools[IPName] = addr

	printIP(pools)

	return nil
}

// datacenter.
func ShowDC(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch data-centers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	centName := args[0]
	center := &dc.Datacenter{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, center); e != nil {
		return e
	}

	manyCenters := map[string]*dc.Datacenter{}

	manyCenters[centName] = center

	printDC(manyCenters)

	return nil
}

func ShowLab(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch labs\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	labName := args[0]

	lab := &dc.Lab{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, lab); e != nil {
		return e
	}

	manyLabs := map[string]*dc.Lab{}

	manyLabs[labName] = lab

	printLab(manyLabs)

	return nil
}

func ShowRack(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch racks\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	vcName := args[0]
	rack := &dc.Rack{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, rack); e != nil {
		return e
	}

	manyRacks := map[string]*dc.Rack{}

	manyRacks[vcName] = rack

	printRack(manyRacks)

	return nil
}

// server.
func ShowServ(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch servers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	srvName := args[0]
	srv := &compute.Server{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, srv); e != nil {
		return e
	}

	manySrv := map[string]*compute.Server{}

	manySrv[srvName] = srv

	printServer(manySrv)

	return nil
}

func ShowESX(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch ESX servers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	esxName := args[0]
	esx := &compute.ESX{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, esx); e != nil {
		return e
	}

	manyESX := map[string]*compute.ESX{}
	manyESX[esxName] = esx

	printESX(manyESX)

	return nil
}

func ShowVC(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch V Centers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	vcName := args[0]
	vc := &compute.VCenter{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, vc); e != nil {
		return e
	}

	manyVC := map[string]*compute.VCenter{}
	manyVC[vcName] = vc

	printVC(manyVC)

	return nil
}

func ShowVM(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch V Centers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)
	if e != nil {
		return e
	}

	vcName := args[0]
	vm := &compute.VCenter{} //nolint:exhaustruct,exhaustivestruct

	if _, e := config.Get("", nil, vm); e != nil {
		return e
	}

	manyVM := map[string]*compute.VCenter{}

	manyVM[vcName] = vm

	printVC(manyVM)

	return nil
}

func printSwitch(srv map[string]*network.Switch) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Management IP", "Credentials", "Serial Number", "Model", "Ports", "Labels"})

	for piece, sw := range srv {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(sw.ID),
			sw.ManagementIP.String(),

			fmt.Sprintf("%s", sw.Credentials.Keys),
			fmt.Sprintf(sw.SerialNumber),
			fmt.Sprintf(sw.Model),

			fmt.Sprintf("%010d", sw.NumPorts),
			fmt.Sprintf("%s", sw.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printLab(labs map[string]*dc.Lab) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Type", "Labels"})

	for piece, lb := range labs {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(lb.NamedResource.ID),
			fmt.Sprintf(lb.NamedResource.Name),

			fmt.Sprintf(lb.NamedResource.Type),
			fmt.Sprintf("%s", lb.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printDC(dcs map[string]*dc.Datacenter) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Type", "Address", "Labels"})

	for piece, dc := range dcs {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(dc.NamedResource.ID),
			fmt.Sprintf(dc.NamedResource.Name),

			fmt.Sprintf(dc.NamedResource.Type),
			fmt.Sprintf(dc.Address),
			fmt.Sprintf("%s", dc.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printServer(servers map[string]*compute.Server) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Board IP", "Type", "Model", "Credentials", "Labels"})

	for piece, s := range servers {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(s.NamedResource.ID),
			fmt.Sprintf(s.NamedResource.Name),
			s.BoardIP.String(),

			fmt.Sprintf(s.NamedResource.Type),
			fmt.Sprintf(s.Model),

			fmt.Sprintf("%s", s.Credentials.Keys),
			fmt.Sprintf("%s", s.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printESX(manyEsx map[string]*compute.ESX) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Server ID", "IP", "Type", "Credentials", "Labels"})

	for piece, esx := range manyEsx {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(esx.NamedResource.ID),
			fmt.Sprintf(esx.NamedResource.Name),

			fmt.Sprintf(esx.ServerID),
			esx.IP.String(),

			fmt.Sprintf(esx.NamedResource.Type),
			fmt.Sprintf("%s", esx.Credentials.Keys),
			fmt.Sprintf("%s", esx.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printVC(manyVC map[string]*compute.VCenter) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "IP", "Type", "Credentials", "Labels"})

	for piece, vc := range manyVC {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vc.NamedResource.ID),

			fmt.Sprintf(vc.NamedResource.Name),
			vc.IP.String(),
			fmt.Sprintf(vc.NamedResource.Type),
			fmt.Sprintf("%s", vc.Credentials.Keys),
			fmt.Sprintf("%s", vc.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printRack(racks map[string]*dc.Rack) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Type", "Row", "Labels"})

	for piece, rack := range racks {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(rack.NamedResource.ID),
			fmt.Sprintf(rack.NamedResource.Name),

			fmt.Sprintf(rack.NamedResource.Type),
			fmt.Sprintf(rack.Row),
			fmt.Sprintf("%s", rack.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printUser(users map[string]*auth.User) error {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Type", "Password Hash", "Role", "Priviledges", "Labels"})

	for piece, user := range users {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(user.NamedResource.ID),
			fmt.Sprintf(user.NamedResource.Name),

			fmt.Sprintf(user.NamedResource.Type),
			fmt.Sprintf(user.PasswordHash),
			fmt.Sprintf(user.Role.Name),

			fmt.Sprintf("%s", user.Role.Privileges),
			fmt.Sprintf("%s", user.NamedResource.Labels),
		})
	}

	fmt.Println(data.Render())
	return nil
}

func printNets(vlans map[string]*network.VLANPool) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Status", "Type", "Range Start", "Range End", "Labels"})

	for piece, vlan := range vlans {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vlan.ID),
			fmt.Sprintf(vlan.Status.UsedBy),
			fmt.Sprintf(vlan.Type),
			fmt.Sprintf("%010d", vlan.RangeStart),
			fmt.Sprintf("%010d", vlan.RangeEnd),
			fmt.Sprintf("%s", vlan.Labels),
		})
	}

	fmt.Println(data.Render())
}

func printIP(vlans map[string]*network.IPAddressPool) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Status", "Type", "Subnets", "Labels"})

	for piece, pool := range vlans {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(pool.ID),
			fmt.Sprintf(pool.Status.UsedBy),
			fmt.Sprintf(pool.Type),
			fmt.Sprintf("%s", pool.Subnets),
			fmt.Sprintf("%s", pool.Labels),
		})
	}

	fmt.Println(data.Render())
}
