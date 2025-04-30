/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"log"
	pb "nsclient/nameserver/proto"
	"strconv"

	"github.com/spf13/cobra"
)

func createLocationType(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	location_type := pb.LocationType{Name: name}
	resp, err := client.CreateLocationType(ctx, &pb.CreateLocationTypeRequest{LocationType: &location_type})
	if err == nil {
		log.Printf("Response: %v", resp.GetLocationType())
	}
	return err
}

func createLocation(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	location := pb.Location{Name: name, LocationTypeName: "Rack"}
	resp, err := client.CreateLocation(ctx, &pb.CreateLocationRequest{Location: &location})
	if err == nil {
		log.Printf("Response: %v", resp.GetLocation())
	}
	return err
}

func createNode(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	location, _ := cmd.Flags().GetString("parent")
	if location == "" {
		return errors.New("Parent location is required")
	}
	node := pb.Node{Hostname: name, LocationName: location, IpAddress: "0.0.0.0"}
	resp, err := client.CreateNode(ctx, &pb.CreateNodeRequest{Node: &node})
	if err == nil {
		log.Printf("Response: %v", resp.GetNode())
	}
	return err
}

func createDevice(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	node, _ := cmd.Flags().GetString("parent")
	if node == "" {
		return errors.New("Parent node is required")
	}
	device := pb.Device{Name: name, NodeHostname: node}
	resp, err := client.CreateDevice(ctx, &pb.CreateDeviceRequest{Device: &device})
	if err == nil {
		log.Printf("Response: %v", resp.GetDevice())
	}
	return err
}

func createChannel(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	device, _ := cmd.Flags().GetString("parent")
	if device == "" {
		return errors.New("Parent device is required")
	}
	//channel_alarm := pb.ChannelAlarm{Type: "Major", TriggerCondition: "<0"}
	//channel_alarm2 := pb.ChannelAlarm{Type: "Minor", TriggerCondition: "<10"}
	//channel := pb.Channel{Name: name, DeviceName: device, Alarms: []*pb.ChannelAlarm{&channel_alarm, &channel_alarm2}}
	channel := pb.Channel{Name: name, DeviceName: device}
	resp, err := client.CreateChannel(ctx, &pb.CreateChannelRequest{Channel: &channel})
	if err == nil {
		log.Printf("Response: %v", resp.GetChannel())
	}
	return err
}

func createChannelAlarm(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	channel, _ := cmd.Flags().GetString("parent")
	if channel == "" {
		return errors.New("Parent channel is required")
	}
	trigger := ">= " + strconv.Itoa(randomNumber())
	channel_alarm := pb.ChannelAlarm{Type: name, TriggerCondition: trigger}
	resp, err := client.AddChannelAlarm(ctx, &pb.AddChannelAlarmRequest{ChannelName: channel, Alarm: &channel_alarm})
	if err == nil {
		log.Printf("Response: %v", resp.GetAlarm())
	}
	return err
}

func createRole(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	role := pb.Role{Name: name}
	resp, err := client.CreateRole(ctx, &pb.CreateRoleRequest{Role: &role})
	if err == nil {
		log.Printf("Response: %v", resp.GetRole())
	}
	return err
}

func createAlarmType(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	alarm_type := pb.AlarmType{Name: name}
	resp, err := client.CreateAlarmType(ctx, &pb.CreateAlarmTypeRequest{AlarmType: &alarm_type})
	if err == nil {
		log.Printf("Response: %v", resp.GetAlarmType())
	}
	return err
}

func createChannelAccess(cctx *ClientContext, name string, cmd *cobra.Command) error {
	//TODO
	return nil
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [noun]",
	Short: "creates a new object",
	Long: `Creates a new object of one of the following types:
	` + SupportedNsTypeString(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var createFunctions = map[NsType]func(*ClientContext, string, *cobra.Command) error{
			LocationType: createLocationType,
			Location:     createLocation,
			Node:         createNode,
			Device:       createDevice,
			Channel:      createChannel,
			Role:         createRole,
			ChannelAlarm: createChannelAlarm,
			AlarmType:    createAlarmType,
		}

		noun := StringToNsType(args[0])
		if noun == -1 {
			log.Fatalf("Unknown noun: %s", args[0])
			return
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Fatalf("Name is required")
			return
		}

		createFunc, ok := createFunctions[noun]
		if !ok {
			log.Fatalf("Listing %s is not supported", args[0])
			return
		}

		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		if err := createFunc(&cctx, name, cmd); err != nil {
			log.Fatalf("Failed to create %s: %v", args[0], err)
		} else {
			log.Printf("Created %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("parent", "", "Parent of the object to create")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
