/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
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
		var dfunctions = map[NsType]func(*ClientContext, string, *cobra.Command) error{
			LocationType: dlocationTypeFunc,
			Location:     dlocationFunc,
			Node:         dnodeFunc,
			Device:       ddeviceFunc,
			Channel:      dchannelFunc,
			Role:         droleFunc,
			ChannelAlarm: dchAlarmFunc,
			// 	ChannelAccess : createChannelAccess,
			AlarmType: dalarmTypeFunc,
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

		dfunc, ok := dfunctions[noun]
		if !ok {
			log.Fatalf("Deleting %s is not supported", args[0])
			return
		}

		cctx := CreateClientContext(cmd)
		defer cctx.conn.Close()
		defer cctx.cancel()

		err := dfunc(&cctx, name, cmd)
		if err != nil {
			log.Fatalf("Failed to delete %s: %v", args[0], err)
		}
		fmt.Printf("Deleted %s: %s\n", args[0], name)
	},
}

func dlocationTypeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	_, err := client.DeleteLocationType(ctx, &pb.DeleteLocationTypeRequest{Name: name})
	return err
}

func dlocationFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	_, err := client.DeleteLocation(ctx, &pb.DeleteLocationRequest{Name: name})
	return err
}

func dnodeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	_, err := client.DeleteNode(ctx, &pb.DeleteNodeRequest{Hostname: name})
	return err
}

func ddeviceFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	_, err := client.DeleteDevice(ctx, &pb.DeleteDeviceRequest{Name: name})
	return err
}

func dchannelFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	_, err := client.DeleteChannel(ctx, &pb.DeleteChannelRequest{Name: name})
	return err
}

func dchAlarmFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	channel, _ := cmd.Flags().GetString("parent")
	if channel == "" {
		return errors.New("Parent (channel) is required for ChannelAlarm")
	}
	_, err := client.DeleteChannelAlarm(ctx, &pb.DeleteChannelAlarmRequest{AlarmType: name, ChannelName: channel})
	return err
}

func droleFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	_, err := client.DeleteRole(ctx, &pb.DeleteRoleRequest{Name: name})
	return err
}
func dalarmTypeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	_, err := client.DeleteAlarmType(ctx, &pb.DeleteAlarmTypeRequest{Name: name})
	return err
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deleteCmd.Flags().String("parent", "", "Parent of object to delete. Needed for Channel* objects")
}
