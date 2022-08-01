package main

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

		if _, e := client.Get(path, nil, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	// printUser(manyUsr)

	data := printUser(manyUsr)

	fmt.Println(data.Render())

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

	manyUsr := map[string]*auth.User{}

	if len(args) == 0 {
		if _, e := client.Get("registrations", nil, manyUsr); e != nil {
			return e
		}
	} else {
		path := fmt.Sprintf("registrations/%s", args[0])

		usr := &auth.User{} //nolint:exhaustruct,exhaustivestruct

		if _, e := client.Get(path, nil, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	// printUser(manyUsr)

	data := printUser(manyUsr)

	fmt.Println(data.Render())

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

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	vlans := map[string]*network.VLANPool{}

	if len(args) == 0 {
		if _, e := client.Get("vlans", nil, vlans); e != nil {
			return e
		}
	} else {
		netName := args[0]

		vlan := &network.VLANPool{} //nolint:exhaustruct,exhaustivestruct
		vlans[netName] = vlan
	}

	// printNets(vlans)

	data := printNets(vlans)

	fmt.Println(data.Render())

	return nil
}

func ShowSw(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch switches\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)

	if e != nil {
		return e
	}

	swName := args[0]
	sw := &network.Switch{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, sw); e != nil {
		return e
	}

	manySw := map[string]*network.Switch{}

	manySw[swName] = sw

	// printSwitch(manySw)
	data := printSwitch(manySw)

	fmt.Println(data.Render())

	return nil
}

func ShowIP(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch vlan\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	IPName := args[0]
	addr := &network.IPAddressPool{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, addr); e != nil {
		return e
	}

	pools := map[string]*network.IPAddressPool{}

	pools[IPName] = addr

	// printIP(pools)

	data := printIP(pools)

	fmt.Println(data.Render())

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

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	centName := args[0]
	center := &dc.Datacenter{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, center); e != nil {
		return e
	}

	manyCenters := map[string]*dc.Datacenter{}

	manyCenters[centName] = center

	// printDC(manyCenters)

	data := printDC(manyCenters)

	fmt.Println(data.Render())

	return nil
}

func ShowLab(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch labs\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	labName := args[0]

	lab := &dc.Lab{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, lab); e != nil {
		return e
	}

	manyLabs := map[string]*dc.Lab{}

	manyLabs[labName] = lab

	// printLab(manyLabs)

	data := printLab(manyLabs)

	fmt.Println(data.Render())

	return nil
}

func ShowRack(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch racks\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	vcName := args[0]
	rack := &dc.Rack{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, rack); e != nil {
		return e
	}

	manyRacks := map[string]*dc.Rack{}

	manyRacks[vcName] = rack

	// printRack(manyRacks)

	data := printRack(manyRacks)

	fmt.Println(data.Render())

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

	client, e := NewClient(config)

	if e != nil {
		return e
	}

	fmt.Println(configFile)

	srvName := args[0]

	srv := &compute.Server{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, srv); e != nil {
		return e
	}

	manySrv := map[string]*compute.Server{}

	manySrv[srvName] = srv

	// printServer(manySrv)

	data := printServer(manySrv)

	fmt.Println(data.Render())

	return nil
}

func ShowESX(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch ESX servers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)
	if e != nil {
		return e
	}

	esxName := args[0]
	esx := &compute.ESX{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, esx); e != nil {
		return e
	}

	manyESX := map[string]*compute.ESX{}
	manyESX[esxName] = esx

	// printESX(manyESX)

	data := printESX(manyESX)

	fmt.Println(data.Render())

	return nil
}

func ShowVC(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch V Centers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)

	if e != nil {
		return e
	}

	vcName := args[0]
	vc := &compute.VCenter{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, vc); e != nil {
		return e
	}

	manyVC := map[string]*compute.VCenter{}
	manyVC[vcName] = vc

	// printVC(manyVC)

	data := printVC(manyVC)

	fmt.Println(data.Render())

	return nil
}

func ShowVM(cmd *cobra.Command, args []string) error {
	fmt.Printf("\nfetch V Centers\n")

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	client, e := NewClient(config)

	if e != nil {
		return e
	}

	vcName := args[0]
	vm := &compute.VM{} //nolint:exhaustruct,exhaustivestruct

	if _, e := client.Get("", nil, vm); e != nil {
		return e
	}

	manyVM := map[string]*compute.VM{}

	manyVM[vcName] = vm

	// printVM(manyVM)

	data := printVM(manyVM)

	fmt.Println(data.Render())

	return nil
}

func printSwitch(srv map[string]*network.Switch) table.Writer {
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

	//	fmt.Println(data.Render())

	return data
}

func printLab(labs map[string]*dc.Lab) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printDC(dcs map[string]*dc.Datacenter) table.Writer {
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

	// fmt.Println(data.Render())
	return data
}

func printServer(servers map[string]*compute.Server) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printESX(manyEsx map[string]*compute.ESX) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printVC(manyVC map[string]*compute.VCenter) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printVM(manyVM map[string]*compute.VM) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"ID", "Name", "IP", "Type", "Credentials",
		"ESXID", "VCenter ID", "Management IP", "Labels",
	})

	for piece, vm := range manyVM {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vm.NamedResource.ID),

			fmt.Sprintf(vm.NamedResource.Name),
			fmt.Sprintf(vm.NamedResource.Type),
			fmt.Sprintf("%s", vm.Credentials.Keys),
			fmt.Sprintf(vm.ESXID),
			fmt.Sprintf(vm.VCenterID),
			vm.ManagementIP.String(),
			fmt.Sprintf("%s", vm.NamedResource.Labels),
		})
	}

	// fmt.Println(data.Render())

	return data
}

func printRack(racks map[string]*dc.Rack) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printUser(users map[string]*auth.User) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printNets(vlans map[string]*network.VLANPool) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}

func printIP(vlans map[string]*network.IPAddressPool) table.Writer {
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

	// fmt.Println(data.Render())

	return data
}
