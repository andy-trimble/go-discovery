# go-discovery
Muticast Discovery for Golang

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
