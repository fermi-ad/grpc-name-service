/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	pb "nsclient/nameserver/proto"

	"github.com/spf13/cobra"
)

func createLocationType(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var location_type pb.LocationType

	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &location_type); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.CreateLocationType(ctx, &pb.CreateLocationTypeRequest{LocationType: &location_type})
	return err
}

func createLocation(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var location pb.Location
	jsonstr := cmd.Flag("json").Value.String()

	if err := DecodeJSONToProto(jsonstr, &location); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	_, err := client.CreateLocation(ctx, &pb.CreateLocationRequest{Location: &location})
	return err
}

func createNode(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var node pb.Node
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &node); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.CreateNode(ctx, &pb.CreateNodeRequest{Node: &node})
	return err
}

func createDevice(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var device pb.Device
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &device); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.CreateDevice(ctx, &pb.CreateDeviceRequest{Device: &device})
	return err
}

func createChannel(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	var channel pb.Channel
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channel); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	_, err := client.CreateChannel(ctx, &pb.CreateChannelRequest{Channel: &channel})
	return err
}

func createChannelAlarm(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var channel_alarm pb.AddChannelAlarmRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channel_alarm); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	_, err := client.AddChannelAlarm(ctx, &channel_alarm)
	return err
}

func createRole(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var role pb.Role
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &role); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	_, err := client.CreateRole(ctx, &pb.CreateRoleRequest{Role: &role})
	return err
}

func createAlarmType(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var alarm_type pb.AlarmType
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &alarm_type); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	_, err := client.CreateAlarmType(ctx, &pb.CreateAlarmTypeRequest{AlarmType: &alarm_type})
	return err
}

func createChannelAccess(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	var channel_access pb.AddChannelAccessControlRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channel_access); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}
	_, err := client.AddChannelAccessControl(ctx, &channel_access)
	return err
}

func createChannelTransform(cctx *ClientContext, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	var channel_transform pb.AddChannelTransformRequest
	jsonstr := cmd.Flag("json").Value.String()
	if err := DecodeJSONToProto(jsonstr, &channel_transform); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	_, err := client.AddChannelTransform(ctx, &channel_transform)
	return err
}

func printCreateJsonTemplate(noun NsType) {
	desc := "(optional)Description"
	switch noun {
	case LocationType:
		PrintProtoAsJSON(&pb.LocationType{
			Name:        "LocationType",
			Description: &desc,
		})
	case Location:
		parent := "(optional)ParentLocation"
		PrintProtoAsJSON(&pb.Location{
			Name:               "Location",
			Description:        &desc,
			LocationTypeName:   "LocationType",
			ParentLocationName: &parent,
		})
	case Node:
		PrintProtoAsJSON(&pb.Node{
			Hostname:     "Hostname",
			Description:  &desc,
			IpAddress:    "0.0.0.0",
			LocationName: "Location",
		})

	case Device:
		PrintProtoAsJSON(&pb.Device{
			Name:         "Device",
			Description:  &desc,
			NodeHostname: "NodeHostname",
		})
	case Channel:
		metadata := "(optional)Metadata"
		flag := true
		PrintProtoAsJSON(&pb.Channel{
			Name:        "Channel",
			Description: &desc,
			DeviceName:  "Device",
			Metadata:    &metadata,

			Alarms: []*pb.ChannelAlarm{
				{
					Type:             "AlarmType",
					TriggerCondition: "TriggerCondition",
				},
			},
			Transforms: []*pb.ChannelTransform{
				{
					Name:        "TransformName",
					Transform:   "TransformData",
					Description: "Description",
				},
			},
			Accesscontrols: []*pb.ChannelAccessControl{
				{
					Role:  "Role",
					Read:  &flag,
					Write: &flag,
				},
			},
		})
	case ChannelAccess:
		flag := true
		PrintProtoAsJSON(&pb.AddChannelAccessControlRequest{
			ChannelName: "Channel",
			Accesscontrol: &pb.ChannelAccessControl{
				Role:  "Role",
				Read:  &flag,
				Write: &flag,
			}})

	case Role:
		PrintProtoAsJSON(&pb.Role{
			Name:        "Role",
			Description: &desc,
		})
	case AlarmType:
		PrintProtoAsJSON(&pb.AlarmType{
			Name:        "AlarmType",
			Description: &desc,
		})
	case ChannelAlarm:
		PrintProtoAsJSON(&pb.AddChannelAlarmRequest{
			ChannelName: "Channel",
			Alarm: &pb.ChannelAlarm{
				Type:             "AlarmType",
				TriggerCondition: "TriggerCondition",
			}})
	case ChannelTransform:
		PrintProtoAsJSON(&pb.AddChannelTransformRequest{
			ChannelName: "Channel",
			Transform: &pb.ChannelTransform{
				Name:        "TransformName",
				Transform:   "TransformData",
				Description: "Description",
			}})
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [noun]",
	Short: "creates a new resource",
	Long:  `Creates a new resource of one of the following types:` + SupportedNsTypeString(),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var createFunctions = map[NsType]func(*ClientContext, *cobra.Command) error{
			LocationType:     createLocationType,
			Location:         createLocation,
			Node:             createNode,
			Device:           createDevice,
			Channel:          createChannel,
			Role:             createRole,
			ChannelAlarm:     createChannelAlarm,
			AlarmType:        createAlarmType,
			ChannelAccess:    createChannelAccess,
			ChannelTransform: createChannelTransform,
		}

		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}

		create_template, _ := cmd.Flags().GetBool("json-template")
		if create_template {
			printCreateJsonTemplate(noun)
			return
		}

		createFunc, ok := createFunctions[noun]
		if !ok {
			log.Fatalf("Listing %s is not supported", args[0])
			return
		}

		json := cmd.Flag("json").Value.String()
		if json == "" {
			log.Fatalf("JSON string is required")
		}
		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		if err := createFunc(&cctx, cmd); err != nil {
			log.Fatalf("Failed to create %s: %v", args[0], err)
		} else {
			log.Printf("Created %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
