// tiny-chat
package main

import (
  "fmt"
  "os"
  "net"
  "time"
  "flag"
  "bufio"
)

const defaultListenPort = "127.0.0.1:3000"
const defaultRemotePort = "127.0.0.1:3001"

func main() {
  // command parameters
  selfUrl := flag.String("local", defaultListenPort, "Local listen URL")
  remoteUrl := flag.String("remote", defaultRemotePort, "Remote destination URL")
  flag.Parse()
  fmt.Println("selfUrl=", *selfUrl)
  fmt.Println("remoteUrl=", *remoteUrl)

  // channel create
  // cinData  := make(chan []byte)

  // accept remote access
  rin, err := net.Listen("tcp", *selfUrl)
  defer rin.Close()
  if err != nil {
    fmt.Println("can't open listen port")
    return
  }
  go func(conn net.Listener) {
    c, _ := rin.Accept()
    defer c.Close()
    buffer := make([]byte, 1024)
    for {
      n,err := c.Read(buffer)
      fmt.Println("read chars = ", n)
      if err != nil {
        fmt.Println("Read error: ", err)
        return
      }
      fmt.Println(string(buffer[:n]))
    }
  }(rin)

  // connect remote
  var remoteConn net.Conn
  maxRetries := 3
  retryDelay := 5 * time.Second
  for retries := 0; retries < maxRetries; retries++ {
    var err error
    conn, err := net.Dial("tcp", *remoteUrl)
    if err == nil {
      fmt.Println("Connection successful!")
      remoteConn = conn
      break
    }

    fmt.Printf("Connection failed (attempt %d): %v\n", retries+1, err)
    time.Sleep(retryDelay)
  }
  defer remoteConn.Close()

  // process console input
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Println("input = ", line)
    remoteConn.Write([]byte(line))
    fmt.Println("message = ", string(line))
  }
}
