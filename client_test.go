package chef

import (
	"fmt"
  "os"
  "testing"
)

func TestClient(t *testing.T) {
	c := Config{
		Name:    "yabuta",
		Key:     "/Users/ryany/.chef/yabuta-triggit.pem",
		BaseURL: "https://nj-chef-server.triggit.com:4443",
	}
	client, err := c.Client()
	if err != nil {
		fmt.Println("Issue setting up client:", err)
		os.Exit(1)
	}
  _, err = client.Nodes.List()
  if err != nil {
    fmt.Println("Can't output node list:", err)
    os.Exit(1)
  }
}
