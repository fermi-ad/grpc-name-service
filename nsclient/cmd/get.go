package cmd

import (
	"log"
	pb "nsclient/nameserver/proto"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [noun]",
	Short: "Gets a specified resource",
	Long: `Gets a specified resource from the system. Supports one of the following types:
	` + SupportedNsTypeString(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var functions = map[NsType]func(*ClientContext, *cobra.Command) error{
			Location: glocationFunc,
			Node:     gnodeFunc,
			Device:   gdeviceFunc,
			Channel:  gchannelFunc,
		}
		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}

		json_template, _ := cmd.Flags().GetBool("json-template")
		if json_template {
			printGetJsonTemplate(noun)
			return
		}
		json, _ := cmd.Flags().GetString("json")
		if json == "" {
			log.Fatalf("JSON is required")
			return
		}
		gfunc, ok := functions[noun]
		if !ok {
			log.Fatalf("Getting %s is not supported", args[0])
			return
		}
		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		err := gfunc(&cctx, cmd)
		if err != nil {
			log.Fatalf("Failed to get %s: %v", args[0], err)
		}
	},
}

func glocationFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.GetLocationRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.GetLocation(ctx, &request)
	if err == nil {
		log.Printf("Get location: %v", r.GetLocation())
	}
	return err
}

func gnodeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.GetNodeRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.GetNode(ctx, &request)
	if err == nil {
		log.Printf("Get node: %v", r.GetNode())
	}
	return err
}

func gdeviceFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.GetDeviceRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.GetDevice(ctx, &request)
	if err == nil {
		log.Printf("Get device: %v", r.GetDevice())
	}
	return err
}

func gchannelFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var request pb.GetChannelRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &request); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.GetChannel(ctx, &request)
	if err == nil {
		log.Printf("Get channel: %v", r.GetChannel())
	}
	return err
}

func printGetJsonTemplate(noun NsType) {
	switch noun {
	case Location:
		PrintProtoAsJSON(&pb.GetLocationRequest{
			Name: "LocationName",
		})
	case Node:
		PrintProtoAsJSON(&pb.GetNodeRequest{
			Hostname: "NodeHostname",
		})
	case Device:
		PrintProtoAsJSON(&pb.GetDeviceRequest{
			Name: "DeviceName",
		})
	case Channel:
		PrintProtoAsJSON(&pb.GetChannelRequest{
			Name: "ChannelName",
		})
	default:
		log.Fatalf("Unsupported noun for get JSON template: %v", noun)
	}
}

func init() {
	rootCmd.AddCommand(getCmd)
}
