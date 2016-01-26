package peer

import (
  "net"
  "fmt"
)

/*
  Node must:
    ## primary api
    - receive requests
    - forward them on (with varying delay)
    - fulfill them if possible
      
      ## non-api requirements
      - track potential forwarding addresses
      - decide who to forward requests too
      
      ## variables & tracking
      - t: how many cached files
      - t: how many requests have I received
      - t: how many have I fulfilled a request
*/

type Connection struct {
  Address string
  Conn net.Conn
}

type Node struct {
  Address string
}

func (n *Node) Connect(addr string) error {
  conn, err := net.Dial("tcp", addr)
  if err != nil {
    return err
  }
  fmt.Fprintf(conn, "Hithere")
  fmt.Println("connection made")
  return nil
}

func (n *Node) Listen(start chan<- bool) {
  l, err := net.Listen("tcp", "localhost:0")
  if err != nil {
    return
  }
  n.Address = l.Addr().String()
  defer l.Close()
  start <- true
  for {
    conn, err := l.Accept()
    if err != nil {
      /* log error */
    }
    go n.HandleConnection(conn)
  }
}

func (n *Node) HandleConnection(conn net.Conn) {
  /* 
    either satisfy or forward request on
  */
  fmt.Fprintf(conn, "Hihi")
  fmt.Println("got a connection")
  conn.Close()
}

func main() {

  fmt.Println("HEHEH")
}
