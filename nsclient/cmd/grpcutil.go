package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	pb "nsclient/nameserver/proto"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"crypto/tls"
	"crypto/x509"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type NsType int64

const (
	LocationType NsType = iota
	Location
	Node
	Device
	Channel
	Role
	ChannelAccess
	AlarmType
	ChannelAlarm
	ChannelTransform
)

var nsTypeNames = map[NsType]string{
	LocationType:     "location_type",
	Location:         "location",
	Node:             "node",
	Device:           "device",
	Channel:          "channel",
	Role:             "role",
	ChannelAccess:    "channel_access",
	AlarmType:        "alarm_type",
	ChannelAlarm:     "channel_alarm",
	ChannelTransform: "channel_transform",
}

func StringToNsType(s string) NsType {
	for k, v := range nsTypeNames {
		if s == v {
			return k
		}
	}
	return -1
}

func SupportedNsTypeString() string {
	var supportedTypes []string
	for _, v := range nsTypeNames {
		supportedTypes = append(supportedTypes, v)
	}
	return "Supported types: " + strings.Join(supportedTypes, ", ")
}

// NsTypeToString converts an NsType to its string representation
func NsTypeToString(t NsType) string {
	switch t {
	case LocationType:
		return "locationtype"
	case Location:
		return "location"
	case Node:
		return "node"
	case Device:
		return "device"
	case Channel:
		return "channel"
	case Role:
		return "role"
	case AlarmType:
		return "alarm_type"
	case ChannelAccess:
		return "channel_access"
	case ChannelAlarm:
		return "channel_alarm"
	default:
		return "unknown"
	}
}

type PasswordGrantClient struct {
	tokenURL     string
	clientID     string
	clientSecret string
	username     string
	password     string
	token        string
	tokenExpiry  time.Time
}

// NewPasswordGrantClient creates a new client for password grant
func NewPasswordGrantClient(cmd *cobra.Command) *PasswordGrantClient {

	keycloakURL := viper.GetString("keycloak-url")
	if keycloakURL == "" {
		log.Fatalf("Keycloak URL is required")
	}

	realm := viper.GetString("auth.realm")
	if realm == "" {
		log.Fatalf("Realm is required")
	}

	clientID := viper.GetString("auth.client-id")
	if clientID == "" {
		log.Fatalf("Client ID is required")
	}
	clientSecret := viper.GetString("auth.client-secret")
	if clientSecret == "" {
		log.Fatalf("Client secret is required")
	}

	username := viper.GetString("user")
	password := viper.GetString("password")
	if username == "" || password == "" {
		log.Fatalf("Username and password are required")
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, realm)

	return &PasswordGrantClient{
		tokenURL:     tokenURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		username:     username,
		password:     password,
	}
}

// GetToken gets a new token or returns the existing one if it's still valid
func (c *PasswordGrantClient) GetToken(ctx context.Context) (string, error) {
	// Check if we have a valid token
	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	// Prepare the form data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("username", c.username)
	data.Set("password", c.password)

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", c.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get token: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Store the token
	c.token = result.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn-30) * time.Second)

	return c.token, nil
}

func DebugToken(token string) {
	// Parse the token without verification (just for debugging)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Fatalf("invalid token format")
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatalf("failed to decode token payload: %w", err)
	}

	// Pretty print the JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, payload, "", "  "); err != nil {
		log.Fatalf("failed to format token payload: %w", err)
	}

	fmt.Println("JWT Token payload:")
	fmt.Println(prettyJSON.String())
}

func createGrpcConn(cmd *cobra.Command) *grpc.ClientConn {
	addr, _ := cmd.Flags().GetString("addr")
	noTls, _ := cmd.Flags().GetBool("no-tls")

	var creds credentials.TransportCredentials

	if addr == "" {
		addr = "localhost"
		if noTls {
			addr += ":8080"
		} else {
			addr += ":8443"
		}
	}
	if noTls {
		creds = insecure.NewCredentials()
	} else {
		tlsPath := viper.GetString("ssl-cert")
		if tlsPath == "" {
			log.Fatalf("SSL certificate path is required")
		}
		cacert, err := os.ReadFile(tlsPath)
		if err != nil {
			log.Fatalf("Failed to load CA cert: %v", err)
		}
		certpool := x509.NewCertPool()
		if !certpool.AppendCertsFromPEM(cacert) {
			log.Fatalf("Failed to add CA cert to pool")
		}
		tlscfg := &tls.Config{RootCAs: certpool, ServerName: "Nameserver"}

		creds = credentials.NewTLS(tlscfg)
	}
	log.Printf("Connecting to %s", addr)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

type ClientContext struct {
	conn    *grpc.ClientConn
	client  pb.NameServiceClient
	context context.Context
	cancel  context.CancelFunc
}

// Given the command-line arguments "cmd", authenicates user,
// creates a Context object with authorization Bearer token attached, and
// creates a gRPC connection to the server.
func CreateClientContext(cmd *cobra.Command) ClientContext {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	noAuth, _ := cmd.Flags().GetBool("no-authnz")
	verbose, _ := cmd.Flags().GetBool("verbose")
	if !noAuth {
		authClient := NewPasswordGrantClient(cmd)
		token, err := authClient.GetToken(ctx)
		if verbose {
			DebugToken(token)
		}
		if err != nil {
			log.Fatalf("Failed to get token: %v", err)
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	}
	conn := createGrpcConn(cmd)
	c := pb.NewNameServiceClient(conn)
	return ClientContext{conn, c, ctx, cancel}
}

// decodes JSON string to a proto buffer message object
func DecodeJSONToProto(jsonStr string, protoMsg proto.Message) error {
	err := protojson.Unmarshal([]byte(jsonStr), protoMsg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// encode proto buffer message object to JSON string
func PrintProtoAsJSON(protoMsg proto.Message) {
	jsonData, err := protojson.Marshal(protoMsg)
	if err != nil {
		log.Fatalf("Failed to marshal proto message to JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

func currentTime() string {
	return time.Now().Format(time.RFC3339)
}

func randomNumber() int {
	return int(time.Now().UnixNano() % 1000)
}

func PrettyPrintProto(message proto.Message) {
	// Use protojson to marshal the proto message into JSON with indentation.
	marshaler := protojson.MarshalOptions{
		Multiline: true, // Enables pretty-printing with newlines and indentation.
		Indent:    "  ", // Sets the indentation level.
	}
	jsonBytes, err := marshaler.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal proto message: %v", err)
		return
	}

	// Print the JSON string.
	fmt.Println(string(jsonBytes))
}
