# go-discovery
Muticast Discovery for Golang. This library is compatible with the Java impelemtation found here https://github.com/andy-trimble/discovery.

See the example for usage. Here's a basic example:

```Golang
d := discovery.Discovery{}
err := d.Start("server")
if err != nil {
	log.Fatal(err)
}
defer d.Shutdown()

go func() {
	for {
		actor := <-d.Discovered
		log.Printf("Discovered %+v", actor)
	}
}()

for {
	err := <-d.Err
    log.Printf("%+v", err)
}
```

The configuration options are specified by environment variables. The variables that can be set are as follows:

```
DISCOVERY_INTERFACE - The network interface to utilize - default: eth0
DISCOVERY_GROUP - The multicast group - default: 230.1.1.1
DISCOVERY_PORT - The multicast port - default: 8989
DISCOVERY_ANNOUNCE_COUNT - The number of times to announce oneself - default: 5
DISCOVERY_ANNOUNCH_WAIT - The duration between subsequent announcements - default: 500 ms
```
