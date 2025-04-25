/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	pb "nsclient/nameserver/proto"
	"strconv"

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
		var functions = map[NsType]func(*ClientContext, string, *cobra.Command) error{

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
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Fatalf("Name is required")
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

		err := ufunc(&cctx, name, cmd)
		if err != nil {
			log.Fatalf("Failed to update %s: %v", args[0], err)
		}
		fmt.Printf("Updated %s: %s\n", args[0], name)
	},
}

func ulocationTypeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	locationType := pb.LocationType{Name: name, Description: &desc}
	r, err := client.UpdateLocationType(ctx, &pb.UpdateLocationTypeRequest{LocationType: &locationType})
	if err == nil {
		log.Printf("Updated location type: %v", r.GetLocationType())
	}
	return err
}

func ulocationFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	location := pb.Location{Name: name, Description: &desc}
	r, err := client.UpdateLocation(ctx, &pb.UpdateLocationRequest{Location: &location})
	if err == nil {
		log.Printf("Updated location: %v", r.GetLocation())
	}
	return err
}

func unodeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	node := pb.Node{Hostname: name, Description: &desc}
	r, err := client.UpdateNode(ctx, &pb.UpdateNodeRequest{Node: &node})
	if err == nil {
		log.Printf("Updated node: %v", r.GetNode())
	}
	return err
}

func udeviceFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	device := pb.Device{Name: name, Description: &desc}
	r, err := client.UpdateDevice(ctx, &pb.UpdateDeviceRequest{Device: &device})
	if err == nil {
		log.Printf("Updated device: %v", r.GetDevice())
	}
	return err
}
func uchannelFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	channel := pb.Channel{Name: name, Description: &desc}
	r, err := client.UpdateChannel(ctx, &pb.UpdateChannelRequest{Channel: &channel})
	if err == nil {
		log.Printf("Updated channel: %v", r.GetChannel())
	}
	return err
}

func uroleFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	role := pb.Role{Name: name, Description: &desc}
	r, err := client.UpdateRole(ctx, &pb.UpdateRoleRequest{Role: &role})
	if err == nil {
		log.Printf("Updated role: %v", r.GetRole())
	}
	return err
}

func ualarmTypeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	desc := "Updated at " + currentTime()
	alarmType := pb.AlarmType{Name: name, Description: &desc}
	r, err := client.UpdateAlarmType(ctx, &pb.UpdateAlarmTypeRequest{AlarmType: &alarmType})
	if err == nil {
		log.Printf("Updated alarm type: %v", r.GetAlarmType())
	}
	return err
}

func uchannelAlarmFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
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
		log.Printf("Updated channel alarm: %v", resp.GetAlarm())
	}
	return err
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().String("parent", "", "Parent of the object to create")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
