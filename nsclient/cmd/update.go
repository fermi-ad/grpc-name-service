/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	pb "nsclient/nameserver/proto"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [noun]",
	Short: "Updates a specified resource",
	Long: `Updates a specified resource in the system. Supports one of the following types:
	` + SupportedNsTypeString(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var functions = map[NsType]func(*ClientContext, *cobra.Command) error{

			LocationType: ulocationTypeFunc,
			Location:     ulocationFunc,
			Node:         unodeFunc,
			Device:       udeviceFunc,
			Channel:      uchannelFunc,
			Role:         uroleFunc,
			// 	ChannelAccess : createChannelAccess,
			ChannelAlarm: uchannelAlarmFunc,
			AlarmType:    ualarmTypeFunc,
		}
		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}
		json_template, _ := cmd.Flags().GetBool("json-template")
		if json_template {
			printUpdateJsonTemplate(noun)
			return
		}
		json, _ := cmd.Flags().GetString("json")
		if json == "" {
			log.Fatalf("JSON is required")
			return
		}

		ufunc, ok := functions[noun]
		if !ok {
			log.Fatalf("Updating %s is not supported", args[0])
			return
		}
		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		err := ufunc(&cctx, cmd)
		if err != nil {
			log.Fatalf("Failed to update %s: %v", args[0], err)
		}
		fmt.Printf("Updated %s\n", args[0])
	},
}

func ulocationTypeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var locationType pb.LocationType
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &locationType); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateLocationType(ctx, &pb.UpdateLocationTypeRequest{LocationType: &locationType})
	if err == nil {
		log.Printf("Updated location type: %v", r.GetLocationType())
	}
	return err
}

func ulocationFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var location pb.Location
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &location); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: &location})
	if err == nil {
		log.Printf("Updated location: %v", r.GetLocation())
	}
	return err
}

func unodeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var node pb.Node
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &node); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateNode(ctx, &pb.UpdateNodeRequest{Node: &node})
	if err == nil {
		log.Printf("Updated node: %v", r.GetNode())
	}
	return err
}

func udeviceFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var device pb.Device
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &device); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateDevice(ctx, &pb.UpdateDeviceRequest{Device: &device})
	if err == nil {
		log.Printf("Updated device: %v", r.GetDevice())
	}
	return err
}
func uchannelFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var channel pb.Channel
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channel); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateChannel(ctx, &pb.UpdateChannelRequest{Channel: &channel})
	if err == nil {
		log.Printf("Updated channel: %v", r.GetChannel())
	}
	return err
}

func uroleFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var role pb.Role
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &role); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateRole(ctx, &pb.UpdateRoleRequest{Role: &role})
	if err == nil {
		log.Printf("Updated role: %v", r.GetRole())
	}
	return err
}

func ualarmTypeFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var alarmType pb.AlarmType
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &alarmType); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateAlarmType(ctx, &pb.UpdateAlarmTypeRequest{AlarmType: &alarmType})
	if err == nil {
		log.Printf("Updated alarm type: %v", r.GetAlarmType())
	}
	return err
}

func uchannelAlarmFunc(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var channelAlarm pb.UpdateChannelAlarmRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channelAlarm); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	r, err := client.UpdateChannelAlarm(ctx, &channelAlarm)
	if err == nil {
		log.Printf("Updated channel alarm: %v", r.GetAlarm())
	}
	return err
}

func printUpdateJsonTemplate(noun NsType) {
	optional := "Description"
	switch noun {
	case LocationType:
		PrintProtoAsJSON(
			&pb.LocationType{
				Name:        "LocationTypeName",
				Description: &optional,
			})
	case Location:
		PrintProtoAsJSON(&pb.Location{
			Name:               "LocationName",
			Description:        &optional,
			LocationTypeName:   "LocationTypeName",
			ParentLocationName: &optional,
		})
	case Node:
		PrintProtoAsJSON(&pb.Node{
			Hostname:     "NodeHostname",
			Description:  &optional,
			IpAddress:    "IPAddress",
			LocationName: "LocationName",
		})
	case Device:
		PrintProtoAsJSON(&pb.Device{
			Name:         "DeviceName",
			Description:  &optional,
			NodeHostname: "NodeHostname",
		})
	case Channel:
		PrintProtoAsJSON(&pb.Channel{
			Name:        "ChannelName",
			Description: &optional,
			DeviceName:  "DeviceName",
		})
	case Role:
		PrintProtoAsJSON(&pb.Role{
			Name:        "RoleName",
			Description: &optional,
		})
	case AlarmType:
		PrintProtoAsJSON(&pb.AlarmType{
			Name:        "AlarmTypeName",
			Description: &optional,
		})
	case ChannelAlarm:
		PrintProtoAsJSON(&pb.UpdateChannelAlarmRequest{
			ChannelName: "ChannelName",
			Alarm: &pb.ChannelAlarm{
				Type:             "AlarmType",
				TriggerCondition: "TriggerCondition",
			},
		})
	case ChannelTransform:
		PrintProtoAsJSON(&pb.UpdateChannelTransformRequest{
			ChannelName: "ChannelName",
			Transform: &pb.ChannelTransform{
				Name:        "TransformName",
				Transform:   "TransformData",
				Description: "Description",
			},
		})
	case ChannelAccess:
		PrintProtoAsJSON(&pb.UpdateChannelAccessControlRequest{
			ChannelName: "ChannelName",
			Accesscontrol: &pb.ChannelAccessControl{
				Role:  "RoleName",
				Read:  new(bool), // Set to true or false
				Write: new(bool), // Set to true or false
			},
		})
	default:
		log.Fatalf("Unsupported noun for update JSON template: %v", noun)
	}
}

func init() {
	rootCmd.AddCommand(updateCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
