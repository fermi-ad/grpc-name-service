package cmd

import (
	"fmt"

	"log"
	pb "nsclient/nameserver/proto"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [noun]",
	Short: "Deletes a specified resource",
	Long: `Deletes a specified resource from the system. Supports one of the following types:
	` + SupportedNsTypeString(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var dfunctions = map[NsType]func(*ClientContext, *cobra.Command) error{
			LocationType: deleteLocationType,
			Location:     deleteLocation,
			Node:         deleteNode,
			Device:       deleteDevice,
			Channel:      deleteChannel,
			Role:         deleteRole,
			ChannelAlarm: deleteChannelAlarm,
			AlarmType:    deleteAlarmType,
		}

		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}

		json_template, _ := cmd.Flags().GetBool("json-template")
		if json_template {
			printDeleteJsonTemplate(noun)
			return
		}

		json, _ := cmd.Flags().GetString("json")
		if json == "" {
			log.Fatalf("JSON is required")
			return
		}

		dfunc, ok := dfunctions[noun]
		if !ok {
			log.Fatalf("Deleting %s is not supported", args[0])
			return
		}

		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		err := dfunc(&cctx, cmd)
		if err != nil {
			log.Fatalf("Failed to delete %s: %v", args[0], err)
		}
		fmt.Printf("Deleted %s\n", args[0])
	},
}

func printDeleteJsonTemplate(noun NsType) {
	switch noun {
	case LocationType:
		PrintProtoAsJSON(&pb.DeleteLocationTypeRequest{
			Name: "LocationTypeName",
		})
	case Location:
		PrintProtoAsJSON(&pb.DeleteLocationRequest{
			Name: "LocationName",
		})
	case Node:
		PrintProtoAsJSON(&pb.DeleteNodeRequest{
			Hostname: "NodeHostname",
		})
	case Device:
		PrintProtoAsJSON(&pb.DeleteDeviceRequest{
			Name: "DeviceName",
		})
	case Channel:
		PrintProtoAsJSON(&pb.DeleteChannelRequest{
			Name: "ChannelName",
		})
	case Role:
		PrintProtoAsJSON(&pb.DeleteRoleRequest{
			Name: "RoleName",
		})
	case ChannelAlarm:
		PrintProtoAsJSON(&pb.DeleteChannelAlarmRequest{
			ChannelName: "ChannelName",
			AlarmType:   "AlarmType",
		})
	case AlarmType:
		PrintProtoAsJSON(&pb.DeleteAlarmTypeRequest{
			Name: "AlarmTypeName",
		})
	case ChannelTransform:
		PrintProtoAsJSON(&pb.DeleteChannelTransformRequest{
			ChannelName:   "ChannelName",
			TransformName: "TransformName",
		})
	case ChannelAccess:
		PrintProtoAsJSON(&pb.DeleteChannelAccessControlRequest{
			ChannelName: "ChannelName",
			Role:        "RoleName",
		})
	default:
		log.Fatalf("Unsupported noun for delete JSON template: %v", noun)
	}
}
func deleteLocationType(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var locationType pb.DeleteLocationTypeRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &locationType); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteLocationType(ctx, &locationType)
	return err
}

func deleteLocation(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var location pb.DeleteLocationRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &location); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteLocation(ctx, &location)
	return err
}

func deleteNode(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var node pb.DeleteNodeRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &node); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteNode(ctx, &node)
	return err
}

func deleteDevice(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var device pb.DeleteDeviceRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &device); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteDevice(ctx, &device)
	return err
}

func deleteChannel(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var channel pb.DeleteChannelRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channel); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteChannel(ctx, &channel)
	return err
}

func deleteRole(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var role pb.DeleteRoleRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &role); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteRole(ctx, &role)
	return err
}

func deleteAlarmType(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var alarmType pb.DeleteAlarmTypeRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &alarmType); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteAlarmType(ctx, &alarmType)
	return err
}

func deleteChannelAlarm(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var channelAlarm pb.DeleteChannelAlarmRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channelAlarm); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.DeleteChannelAlarm(ctx, &channelAlarm)
	return err
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
