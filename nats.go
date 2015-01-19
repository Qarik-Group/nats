package main

import(
  "os"
  "log"
  "strings"
  "runtime"

  "github.com/apcera/nats"
  "github.com/codegangsta/cli"
  )

  var message = "Usage: nats sub or nats pub"
  var subMessage = "Usage: nats sub [-s server] [--ssl] [-t] <subject> \n"
  var pubMessage = "Usage: nats pub [-s server] [--ssl] [-t] <subject> <msg> \n"
  var index = 0

func usage() {
  log.Fatalf(message)
}

func printMsg(m *nats.Msg, i int){
  index += 1
  log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main(){
  app := cli.NewApp()
  app.Name = "nats"
  app.Usage = "Nats Pub and Sub - Go Client"
  app.Version = "1.0"
  app.Action = func(c *cli.Context) {
    cli.ShowAppHelp(c)
  }
  app.Commands = []cli.Command{
    {
      Name:       "pub",
      ShortName:  "p",
      Usage:      pubMessage,
      Flags:  []cli.Flag{
        cli.StringFlag{Name:   "s", Value:  nats.DefaultURL, Usage: "The nats server URLs (separated by comma)"},
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

        log.Printf("Published [%s] : '%s'\n", subj, msg)
      },
    },
    {
      Name:       "sub",
      ShortName:  "s",
      Usage:      subMessage,
      Flags:  []cli.Flag{
          cli.StringFlag{Name:   "s", Value:  nats.DefaultURL, Usage: "The nats server URLs (separated by comma)"},
          cli.BoolFlag{Name:   "ssl", Usage:  "Use Secure Connection"},
          cli.BoolFlag{Name:   "t",Usage:  "Display timestamps"},
        },
      Action:  func(c *cli.Context){
        var urls = c.String("s")
        var showTime = c.Bool("t")
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

        subj, i := args[0], 0

        nc.Subscribe(subj, func(msg *nats.Msg) {
          i += 1
          printMsg(msg, i)
          })

          log.Printf("Listening on [%s]\n", subj)
          if showTime {
            log.SetFlags(log.LstdFlags)
          }

          runtime.Goexit()
      },
    },
  }
  app.Run(os.Args)
}
