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
			AlarmType:    lalarmTypeFunc,
		}
		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}

		json_template, _ := cmd.Flags().GetBool("json-template")
		if json_template {
			printListJsonTemplate(noun)
			return
		}

		jsonstr, _ := cmd.Flags().GetString("json")
		if jsonstr == "" {
			cmd.Flags().Set("json", "{}")
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

	var request pb.ListLocationsRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.ListLocations(ctx, &request)
	if err == nil {
		log.Printf("List locations: %v", r.GetLocations())
	}
	return err
}

func lnodeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.ListNodesRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.ListNodes(ctx, &request)
	if err == nil {
		log.Printf("List nodes: %v", r.GetNodes())
	}
	return err
}

func ldeviceFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.ListDevicesRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.ListDevices(ctx, &request)
	if err == nil {
		log.Printf("List devices: %v", r.GetDevices())
	}
	return err
}

func lchannelFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.ListChannelsRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.ListChannels(ctx, &request)
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

func printListJsonTemplate(noun NsType) {
	pagination := pb.PaginationRequest{
		PageSize: 10,
		Page:     1,
	}
	optional := "optional"
	switch noun {
	case LocationType:
	case Role:
	case AlarmType:
		PrintProtoAsJSON(&emptypb.Empty{})
	case Location:
		PrintProtoAsJSON(&pb.ListLocationsRequest{
			LocationTypeName:   &optional,
			ParentLocationName: &optional,
			Pagination:         &pagination,
		})
	case Node:
		PrintProtoAsJSON(&pb.ListNodesRequest{
			Hostname:     &optional,
			IpAddress:    &optional,
			LocationName: &optional,
			Pagination:   &pagination,
		})
	case Device:
		PrintProtoAsJSON(&pb.ListDevicesRequest{
			Name:         &optional,
			NodeHostname: &optional,
			Pagination:   &pagination,
		})
	case Channel:
		PrintProtoAsJSON(&pb.ListChannelsRequest{
			Name:       &optional,
			DeviceName: &optional,
			Pagination: &pagination,
		})
	default:
		log.Fatalf("Unsupported noun for list JSON template: %v", noun)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
