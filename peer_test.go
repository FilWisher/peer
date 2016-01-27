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
  <-node_two.Ready
  node_two.Connect(node_one.Address)
  <-node_one.Ready
  node_two.Request(node_one.Address, "Hell\n")
  node_one.Request(node_two.Address, "GOodbyep\n")
  msg1 := <-node_one.In
  msg2 := <-node_two.In
  if node_one.Requests < 1 || node_two.Requests < 1 {
    t.Errorf("Less than one request received")
  }
  fmt.Printf("2 requests - 1: %s\n2: %s\n", msg1, msg2)
}
