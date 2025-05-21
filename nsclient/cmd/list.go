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

/*
* handlePagination handles the pagination of the list requests. Calls list function until all
* items for the request has been receivied. It returns early with an error if any of the requests fail.
*
* param callListFunc: The function that performs the list request and prints items in the list. Takes in
*                     pagination request and returns the pagination response and number of items received.
 */
func handlePagination(callListFunc func(*pb.PaginationRequest) (*pb.PaginationResponse, int, error)) error {
	var err error
	totalCount := -1
	var page uint32 = 1
	for {
		pagreq := &pb.PaginationRequest{
			Page: page,
		}
		page += 1
		pagres, nrecv, err := callListFunc(pagreq)
		if err == nil {
			if totalCount == -1 {
				totalCount = int(pagres.GetTotalCount())
			}
			totalCount -= nrecv
			if totalCount <= 0 {
				break
			}
		} else {
			break
		}
	}
	return err
}

func llocationTypeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.ListLocationTypes(ctx, &emptypb.Empty{})
	if err == nil {
		for _, locType := range r.GetLocationTypes() {
			PrettyPrintProto(locType)
		}
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

	err := handlePagination(func(pagreq *pb.PaginationRequest) (*pb.PaginationResponse, int, error) {
		request.Pagination = pagreq
		r, err := client.ListLocations(ctx, &request)
		if err == nil {
			for _, location := range r.GetLocations() {
				PrettyPrintProto(location)
			}
			nrecv := len(r.GetLocations())
			return r.GetPagination(), nrecv, nil
		}
		return nil, 0, err
	})
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

	err := handlePagination(func(pagreq *pb.PaginationRequest) (*pb.PaginationResponse, int, error) {
		request.Pagination = pagreq
		r, err := client.ListNodes(ctx, &request)
		if err == nil {
			for _, node := range r.GetNodes() {
				PrettyPrintProto(node)
			}
			nrecv := len(r.GetNodes())
			return r.GetPagination(), nrecv, nil
		}
		return nil, 0, err
	})
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

	err := handlePagination(func(pagreq *pb.PaginationRequest) (*pb.PaginationResponse, int, error) {
		request.Pagination = pagreq
		r, err := client.ListDevices(ctx, &request)
		if err == nil {
			for _, device := range r.GetDevices() {
				PrettyPrintProto(device)
			}
			nrecv := len(r.GetDevices())
			return r.GetPagination(), nrecv, nil
		}
		return nil, 0, err
	})
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

	err := handlePagination(func(pagreq *pb.PaginationRequest) (*pb.PaginationResponse, int, error) {
		request.Pagination = pagreq
		r, err := client.ListChannels(ctx, &request)
		if err == nil {
			for _, channel := range r.GetChannels() {
				PrettyPrintProto(channel)
			}
			nrecv := len(r.GetChannels())
			return r.GetPagination(), nrecv, nil
		}
		return nil, 0, err
	})

	return err
}

func lroleFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.ListRoles(ctx, &emptypb.Empty{})
	if err == nil {
		for _, role := range r.GetRoles() {
			PrettyPrintProto(role)
		}
	}
	return err
}

func lalarmTypeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.ListAlarmTypes(ctx, &emptypb.Empty{})
	if err == nil {
		for _, alarmType := range r.GetAlarmTypes() {
			PrettyPrintProto(alarmType)
		}
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
}
