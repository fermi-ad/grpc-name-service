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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [noun]",
	Short: "Gets a specified resource",
	Long: `Gets a specified resource from the system. Supports one of the following types:
	` + SupportedNsTypeString(),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var functions = map[NsType]func(*ClientContext, string, *cobra.Command) error{
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
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Fatalf("Name is required")
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

		err := gfunc(&cctx, name, cmd)
		if err != nil {
			log.Fatalf("Failed to get %s: %v", args[0], err)
		}
		fmt.Printf("Getd %s: %s\n", args[0], name)
	},
}

func glocationFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.GetLocation(ctx, &pb.GetLocationRequest{Name: name})
	if err == nil {
		log.Printf("Getd location: %v", r.GetLocation())
	}
	return err
}

func gnodeFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context
	r, err := client.GetNode(ctx, &pb.GetNodeRequest{Hostname: name})
	if err == nil {
		log.Printf("Getd node: %v", r.GetNode())
	}
	return err
}

func gdeviceFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.GetDevice(ctx, &pb.GetDeviceRequest{Name: name})
	if err == nil {
		log.Printf("Getd device: %v", r.GetDevice())
	}
	return err
}
func gchannelFunc(cctx *ClientContext, name string, cmd *cobra.Command) error {
	client := cctx.client
	ctx := cctx.context

	r, err := client.GetChannel(ctx, &pb.GetChannelRequest{Name: name})
	if err == nil {
		log.Printf("Getd channel: %v", r.GetChannel())
	}
	return err
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
