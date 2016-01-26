package peer

import (
  "testing"
  "fmt"
)

var (
  node_one Node
  node_two Node
)

func TestMain(m *testing.M) {
  started := make (chan bool)
  go node_one.Listen(started)
  go node_two.Listen(started)
  <-started
  <-started
  m.Run()
}

func TestNode(t *testing.T) {
  fmt.Printf("node_one running on %s\n", node_one.Address)
  fmt.Printf("node_two running on %s\n", node_two.Address)
  fmt.Println("All ok!")
}

func TestComm(t *testing.T) {
  node_one.Connect(node_two.Address)
  node_two.Connect(node_one.Address)
}
