package chef

import (
	"bytes"
	"encoding/json"
	"fmt"
	chefGo "github.com/go-chef/chef"
	"io"
	"os"
	"testing"
)

// JSONReader handles arbitrary types and synthesizes a streaming encoder for them.
func JSONReader(v interface{}) (r io.Reader, err error) {
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(v)
	r = bytes.NewReader(buf.Bytes())
	return
}

func TestClient(t *testing.T) {
	c := Config{
		Name:    "yabuta",
		Key:     "/Users/ryany/.chef/yabuta-triggit.pem",
		BaseURL: "https://nj-chef-server.triggit.com:4443",
	}
	_, err := c.Client()
	// client, err := c.Client()
	if err != nil {
		fmt.Println("Issue setting up client:", err)
		os.Exit(1)
	}
	node := chefGo.NewNode("test.local")
	node.AutomaticAttributes = map[string]interface{}{}
	node.NormalAttributes = map[string]interface{}{}
	node.DefaultAttributes = map[string]interface{}{}
	node.OverrideAttributes = map[string]interface{}{}
	node.RunList = []string{}
	// nr, err := client.Nodes.Post(node)
	// fmt.Println(nr)
	r, err := JSONReader(node)
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	fmt.Println(b.String())
	if err != nil {
		fmt.Println("Can't output node list:", err)
		os.Exit(1)
	}
}
