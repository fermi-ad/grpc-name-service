/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	pb "nsclient/nameserver/proto"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [noun]",
	Short: "Lists a specific type of resource",
	Long: `Lists a specific type of resource from the system. Supports one of the following types:
	` + SupportedNsTypeString(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var lfunctions = map[NsType]func(*ClientContext, *cobra.Command) error{
			LocationType: llocationTypeFunc,
			Location:     llocationFunc,
			Node:         lnodeFunc,
			Device:       ldeviceFunc,
			Channel:      lchannelFunc,
			Role:         lroleFunc,
			// 	ChannelAccess : createChannelAccess,
			AlarmType: lalarmTypeFunc,
		}
		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}

		lfunc, ok := lfunctions[noun]
		if !ok {
			log.Fatalf("Listing %s is not supported", args[0])
			return
		}
		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		err := lfunc(&cctx, cmd)
		if err != nil {
			log.Fatalf("Failed to list %s: %v", args[0], err)
		}
	},
}

func llocationTypeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.ListLocationTypes(ctx, &emptypb.Empty{})
	if err == nil {
		log.Printf("List location types: %v", r.GetLocationTypes())
	}
	return err
}

func llocationFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	qrequest := pb.ListLocationsRequest{}
	name, _ := cmd.Flags().GetString("name")
	if name != "" {
		qrequest.Name = &name
	}
	parent, _ := cmd.Flags().GetString("parent")
	if parent != "" {
		qrequest.ParentLocationName = &parent
	}

	r, err := client.ListLocations(ctx, &qrequest)
	if err == nil {
		log.Printf("List locations: %v", r.GetLocations())
	}
	return err
}

func lnodeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	qrequest := pb.ListNodesRequest{}
	name, _ := cmd.Flags().GetString("name")
	if name != "" {
		qrequest.Hostname = &name
	}
	parent, _ := cmd.Flags().GetString("parent")
	if parent != "" {
		qrequest.LocationName = &parent
	}
	r, err := client.ListNodes(ctx, &qrequest)
	if err == nil {
		log.Printf("List nodes: %v", r.GetNodes())
	}
	return err
}

func ldeviceFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	qrequest := pb.ListDevicesRequest{}
	name, _ := cmd.Flags().GetString("name")
	if name != "" {
		qrequest.Name = &name
	}
	parent, _ := cmd.Flags().GetString("parent")
	if parent != "" {
		qrequest.NodeHostname = &parent
	}
	r, err := client.ListDevices(ctx, &qrequest)
	if err == nil {
		log.Printf("List devices: %v", r.GetDevices())
	}
	return err
}

func lchannelFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	qrequest := pb.ListChannelsRequest{}
	name, _ := cmd.Flags().GetString("name")
	if name != "" {
		qrequest.Name = &name
	}
	parent, _ := cmd.Flags().GetString("parent")
	if parent != "" {
		qrequest.DeviceName = &parent
	}
	r, err := client.ListChannels(ctx, &qrequest)
	if err == nil {
		log.Printf("List channels: %v", r.GetChannels())
	}
	return err
}
func lroleFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.ListRoles(ctx, &emptypb.Empty{})
	if err == nil {
		log.Printf("List roles: %v", r.GetRoles())
	}
	return err
}

func lalarmTypeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.ListAlarmTypes(ctx, &emptypb.Empty{})
	if err == nil {
		log.Printf("List alarm types: %v", r.GetAlarmTypes())
	}
	return err
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().String("parent", "", "Filter by parent of the objects to list")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
