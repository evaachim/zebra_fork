package show

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/spf13/cobra"
)

func showSw(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch switches\n")

	swName := args[0]
	sw := &network.Switch{} //nolint:exhaustruct,exhaustivestruct
	manySw := map[string]*network.Switch{}

	manySw[swName] = sw

	printSwitch(manySw)
}

func showNet(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch network\n")
}

func showLab(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch labs\n")

	labName := args[0]
	manyLabs := map[string]*dc.Lab{}

	lab := &dc.Lab{} //nolint:exhaustruct,exhaustivestruct

	manyLabs[labName] = lab

	printLab(manyLabs)
}

func showDC(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch data-centers\n")

	centName := args[0]
	manyCenters := map[string]*dc.Datacenter{}

	center := &dc.Datacenter{} //nolint:exhaustruct,exhaustivestruct
	manyCenters[centName] = center

	printDC(manyCenters)
}

func showServ(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch servers\n")

	srvName := args[0]
	manySrv := map[string]*compute.Server{}

	srv := &compute.Server{} //nolint:exhaustruct,exhaustivestruct

	manySrv[srvName] = srv

	printServer(manySrv)
}

func showESX(xmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch ESX servers\n")

	esxName := args[0]
	manyESX := map[string]*compute.ESX{}

	esx := &compute.ESX{} //nolint:exhaustruct,exhaustivestruct
	manyESX[esxName] = esx

	printESX(manyESX)
}

func showVC(xmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch V Centers\n")

	vcName := args[0]
	manyVC := map[string]*compute.VCenter{}

	vc := &compute.VCenter{} //nolint:exhaustruct,exhaustivestruct
	manyVC[vcName] = vc

	printVC(manyVC)
}

func showRack(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch racks\n")

	vcName := args[0]
	manyRacks := map[string]*dc.Rack{}

	rack := &dc.Rack{} //nolint:exhaustruct,exhaustivestruct
	manyRacks[vcName] = rack

	printRack(manyRacks)
}

func showReg(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch registrations\n")

	usrName := args[0]
	manyUsr := map[string]*auth.User{}

	usr := &auth.User{} //nolint:exhaustruct,exhaustivestruct
	manyUsr[usrName] = usr

	printUser(manyUsr)
}

func showUsr(cmd *cobra.Command, args []string) {
	fmt.Printf("\nfetch users\n")

	usrName := args[0]
	manyUsr := map[string]*auth.User{}

	usr := &auth.User{} //nolint:exhaustruct,exhaustivestruct
	manyUsr[usrName] = usr

	printUser(manyUsr)
}

func printSwitch(srv map[string]*network.Switch) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"id", "Management IP", "Credentials", "Serial Number", "Model", "Ports", "Labels"})

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
	data.AppendHeader(table.Row{"id", "Name", "ID", "Status", "Type", "Labels"})

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
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Type", "Address", "Labels"})

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
	data.AppendHeader(table.Row{"ID", "Name", "Board IP", "Status", "Type", "Model", "Credentials", "Labels"})

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
	data.AppendHeader(table.Row{"ID", "Name", "Server ID", "IP", "Status", "Type", "Credentials", "Labels"})

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
	data.AppendHeader(table.Row{"ID", "Name", "IP", "Status", "Type", "Credentials", "Labels"})

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

func printUser(users map[string]*auth.User) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Type", "Password Hash", "Role", "Priviledges", "Key", "Labels"})

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
}
