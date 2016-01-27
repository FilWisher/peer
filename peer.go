package peer

import (
  "net"
  "fmt"
  "errors"
  "bufio"
)

/*
  Node must:
    ## primary api
    - [ ] connect
      - add connection to Connections
      - start go routine to listen for incoming
    - [ ] listen
      - listen for incoming connections
      - start go routine to handle
    - [ ] request
      - send message to connection in Connections
    - [ ] receive requests
    - [ ] forward them on (with varying delay)
    - [ ] fulfill them if possible
      
      ## non-api requirements
      - track potential forwarding addresses
      - decide who to forward requests too
      
      ## variables & tracking
      - t: how many cached files
      - t: how many requests have I received
      - t: how many have I fulfilled a request
*/

/* track connections (graph edges) */
type Connection struct {
  Address string
  Conn net.Conn
}

/* data for node */
type Node struct {
  Address string
  Requests int
  Connections map[string] Connection
  In chan string
  Ready chan bool
}

/* check connection for incoming requests */
func (n *Node) CheckIncoming(conn net.Conn, in chan string) {
  scanner := bufio.NewScanner(conn)
  for scanner.Scan() {
    if err := scanner.Err(); err != nil {
      /* TODO: log error */
      continue
    }
    n.Requests += 1
    in <- scanner.Text()
  }
}

/* add specified node to connections */
func (n *Node) Connect(addr string) {
  conn, err := net.Dial("tcp", addr)
  if err != nil {
    return
  }
  n.Connections[addr] = Connection{
    addr,
    conn,
  }
  go n.CheckIncoming(conn, n.In)
}

/* make request to specified node */
func (n *Node) Request(addr, content string) error {
  connection, ok := n.Connections[addr]
  if !ok {
    return errors.New("Address not known")
  }
  fmt.Fprintf(connection.Conn, content)
  return nil
}

/* listen for incoming requests */
func (n *Node) Listen(start chan<- bool) {
  l, err := net.Listen("tcp", "localhost:0")
  if err != nil {
    return
  }
  n.Address = l.Addr().String()
  n.Connections = make(map[string]Connection)
  n.Ready = make(chan bool, 5)
  n.In = make(chan string, 5)
  defer l.Close()
  n.Requests = 0
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
  n.Connections[conn.RemoteAddr().String()] = Connection{
    conn.RemoteAddr().String(),
    conn,
  }
  fmt.Println("got a connection")
  n.Ready <- true
  go n.CheckIncoming(conn, n.In)
}
