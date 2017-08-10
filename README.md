# nats cli
Go Application that combines nats sub and pub

#Download
```
go get github.com/starkandwayne/nats
```

#Actions

## Pub 
Publish a message on a subject

## Sub
Subscribe to a subject and await messages

## Request
Publish a message and await reply


#Usage
```
  nats pub [-s server] [--ssl]  <subject> <msg> 
      or
  nats sub [-s server] [--ssl] [-t] [-r] <subject> 
      or
  nats req [-s server] [--ssl] [-t] [-r] [-w] <subject> <msg>
```
