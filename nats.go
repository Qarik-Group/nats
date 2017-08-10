package main

import(
  "os"
  "log"
  "fmt"
  "strings"
  "runtime"

  "github.com/nats-io/nats"
  "github.com/codegangsta/cli"
  "time"
)

  var message = "Usage: nats sub or nats pub"
  var subMessage = "Usage: nats sub [-s server url] [--ssl] [-t] <subject> \n"
  var pubMessage = "Usage: nats pub [-s server url] [--ssl] [-t] <subject> <msg> \n"
  var reqMessage = "Usage: nats req [-s server url] [--ssl] [-t] [-w] <subject> <msg> \n"
  var index = 0

func usage() {
  log.Fatalf(message)
}

func printMsg(m *nats.Msg, i int){
  index += 1
  fmt.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func printTimeMsg(m *nats.Msg, i int){
  index += 1
  log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func printRaw(m *nats.Msg){
  fmt.Printf("%s\n", string(m.Data))
}

func main(){
  app := cli.NewApp()
  app.Name = "nats"
  app.Usage = "Nats Pub and Sub - Go Client"
  app.Version = "1.0.2"
  app.Action = func(c *cli.Context) error {
    cli.ShowAppHelp(c)
    return nil
  }
  app.Commands = []cli.Command{
    {
      Name:       "pub",
      ShortName:  "p",
      Usage:      pubMessage,
      Flags:  []cli.Flag{
        cli.StringFlag{Name:   "s", Value:  nats.DefaultURL, Usage: "The nats server URLs (separated by comma).\n\tServer Url must be in following format: nats://nats_user:nats_password@host:port or nats://host:port"},
        cli.BoolFlag{Name:   "ssl", Usage:  "Use Secure Connection"},
        },
      Action: func(c *cli.Context){
        var urls = c.String("s")
        var ssl = c.Bool("ssl")

        args := c.Args()
        if len(args) < 1 {
          cli.ShowAppHelp(c)
          os.Exit(1)
        }

        opts := nats.DefaultOptions
        opts.Servers = strings.Split(urls, ",")
        for i, s := range opts.Servers {
          opts.Servers[i] = strings.Trim(s, " ")
        }

        opts.Secure = ssl

        nc, err := opts.Connect()
        if err != nil {
          log.Fatalf("Can't connect: %v\n", err)
        }

        subj, msg := args[0], []byte(args[1])

        nc.Publish(subj, msg)
        nc.Close()

        fmt.Printf("Published [%s] : '%s'\n", subj, msg)
      },
    },
    {
      Name:       "sub",
      ShortName:  "s",
      Usage:      subMessage,
      Flags:  []cli.Flag{
          cli.StringFlag{Name:   "s", Value:  nats.DefaultURL, Usage: "The nats server URLs (separated by comma).\n\tServer Url must be in following format: nats://nats_user:nats_password@host:port or nats://host:port"},
          cli.BoolFlag{Name:   "ssl", Usage:  "Use Secure Connection"},
          cli.BoolFlag{Name:   "t",Usage:  "Display timestamps"},
          cli.BoolFlag{Name:   "r",Usage:  "Display raw output"},
        },
      Action:  func(c *cli.Context){
        var urls = c.String("s")
        var ssl = c.Bool("ssl")
        var showtime = c.Bool("t")
        var rawoutput = c.Bool("r")

        args := c.Args()
        if len(args) < 1 {
          cli.ShowAppHelp(c)
          os.Exit(1)
        }

        opts := nats.DefaultOptions
        opts.Servers = strings.Split(urls, ",")
        for i, s := range opts.Servers {
          opts.Servers[i] = strings.Trim(s, " ")
        }
        opts.Secure = ssl

        nc, err := opts.Connect()
        if err != nil {
          log.Fatalf("Can't connect: %v\n", err)
        }

        subj, i := args[0], 0

        if(rawoutput){
          nc.Subscribe(subj, func(msg *nats.Msg) {
            printRaw(msg)
            })

        }else{
          nc.Subscribe(subj, func(msg *nats.Msg) {
            i += 1
            if(showtime){
              printTimeMsg(msg, i)
            }else{
              printMsg(msg, i)
            }
            })

          fmt.Printf("Listening on [%s]\n", subj)
        }
        runtime.Goexit()
      },
    },
    {
      Name:       "req",
      ShortName:  "r",
      Usage:      reqMessage,
      Flags:  []cli.Flag{
        cli.StringFlag{Name:   "s", Value:  nats.DefaultURL, Usage: "The nats server URLs (separated by comma).\n\tServer Url must be in following format: nats://nats_user:nats_password@host:port or nats://host:port"},
        cli.BoolFlag{Name:   "ssl", Usage:  "Use Secure Connection"},
        cli.DurationFlag{Name: "w", Value: 2 * time.Second, Usage: "Duration to wait for response, eg '1s'"},
        cli.BoolFlag{Name:   "t",Usage:  "Display timestamps"},
        cli.BoolFlag{Name:   "r",Usage:  "Display raw output"},
      },
      Action: func(c *cli.Context){
        var urls = c.String("s")
        var ssl = c.Bool("ssl")
        var timeout = c.Duration("w")
        var showtime = c.Bool("t")
        var rawoutput = c.Bool("r")

        args := c.Args()
        if len(args) < 1 {
          cli.ShowAppHelp(c)
          os.Exit(1)
        }

        opts := nats.DefaultOptions
        opts.Servers = strings.Split(urls, ",")
        for i, s := range opts.Servers {
          opts.Servers[i] = strings.Trim(s, " ")
        }

        opts.Secure = ssl

        nc, err := opts.Connect()
        if err != nil {
          log.Fatalf("Can't connect: %v\n", err)
        }

        subj, msg := args[0], []byte(args[1])

        resp, err := nc.Request(subj, msg, timeout)
        if err != nil {
          log.Fatal("Error during request:", err)
        }
        fmt.Printf("Published [%s] : '%s'\n", subj, msg)

        if(rawoutput){
          printRaw(resp)
        }else{
          if(showtime){
            printTimeMsg(resp, 1)
          }else{
            printMsg(resp, 1)
          }
        }

        nc.Close()
      },
    },
  }
  app.Run(os.Args)
}
