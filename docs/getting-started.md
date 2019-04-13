# Getting Started

Set the `KONG_HOST` environment variable on your machine, KONG_HOST uses the address of your real gateway.

```
export KONG_HOST= http://xx.xx.xx.xx:8001
```

View the features supported by the command line tool.

```
kongctl
NAME:
   kongctl - kong(0.14.0) api gateway command line tool.
   https://docs.konghq.com/0.14.x/admin-api

USAGE:
   kongctl [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     certificate  The kong certificate object.
     consumer     The kong consumer object.
     plugin       The kong plugin object.
     route        The kong route object.
     service      The kong service object.
     snis         The kong sni object.
     target       The kong target object.
     upstream     The kong upstream object.
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --auth value   basic authoritarian for api gateway [$KONG_AUTH]
   --host value   api gateway(kong) server address [$KONG_HOST]
   --help, -h     show help
   --version, -v  print the version
```

### Service object

Service entities, as the name implies, are abstractions of each of your own upstream services. Examples of Services would be a data transformation microservice, a billing API, etc.

```
kongctl service
NAME:
   kongctl service - The kong service object.

USAGE:
   kongctl service command [command options] [arguments...]

COMMANDS:
     create  create service object
     get     retrieve service object
     delete  delete service object
     list    list all services object
     routes  list routes associated to a service

OPTIONS:
   --help, -h  show help
```

Create a service resource object:

```
kongctl service create --help
NAME:
   kongctl service create - create service object

USAGE:
   kongctl service create [command options] [arguments...]

OPTIONS:
   --id value               the service id
   --name value             the service name
   --retries value          the number of retries to execute upon failure to proxy (default: 5)
   --procotol value         the protocol used to communicate with the upstream (default: "http")
   --host value             the host of the upstream server
   --port value             the upstream server port (default: 80)
   --path value             the path to be used requests to the upstream
   --connect_timeout value  the timeout in milliseconds for establishing a connection to the upstream server (default: 60000)
   --write_timeout value    the timeout in milliseconds between two successive write operations for transmitting a request to the upstream server (default: 60000)
   --read_timeout value     the timeout in milliseconds between two successive read operations for transmitting a request to the upstream server (default: 60000)
   --url value              shorthand attribute to set protocol, host, port and path at once
```

example:

```
kongctl service create --name=web --url=http://xx.xx.xx.xx:9999

response body:
{
	"host": "xx.xxx.xxx.xx",
	"created_at": 1553907158,
	"connect_timeout": 60000,
	"id": "495fec4e-8fdf-42cd-b4e9-604b3776cfed",
	"protocol": "http",
	"name": "web",
	"read_timeout": 60000,
	"port": 9999,
	"path": null,
	"updated_at": 1553907158,
	"retries": 5,
	"write_timeout": 60000
}
```


### Route Object

The Route entities defines rules to match client requests. Each Route is associated with a Service, and a Service may have multiple Routes associated to it. Every request matching a given Route will be proxied to its associated Service.

```
kongctl route
NAME:
   kongctl route - The kong route object.

USAGE:
   kongctl route command [command options] [arguments...]

COMMANDS:
     create  create route object
     get     retrieve route object
     delete  delete route object
     list    list all routes object

OPTIONS:
   --help, -h  show help
```

Create a route resource object:

```
kongctl route create --h
NAME:
   kongctl route create - create route object

USAGE:
   kongctl route create [command options] [arguments...]

OPTIONS:
   --protocols value       A list of the protocols this route should allow
   --methods value         A list of HTTP methods that match this Route
   --hosts value           A list of domain names that match this route
   --paths value           A list of paths that match this route
   --regex_priority value  Determines the relative order of this Route against others when evaluating regex paths (default: 0)
   --strip_path            When matching a route via one of the paths, strip the matching prefix from the upstream request URL
   --preserve_host         When matching a route via one of the hosts domain names, use the request Host header in the upstream request headers
   --service_id value      The service id this route is associated to
   --snis value            A list of SNIs that match this route when using stream routing
   --sources value         A list of IP sources of incoming connections that match this route when using stream routing
   --destinations value    A list of IP destinations of incoming connections that match this route when using stream routing
```


example:

```
kongctl route create --paths=/v1/model --service_id=495fec4e-8fdf-42cd-b4e9-604b3776cfed --protocols=http
{
	"created_at": 1553907462,
	"strip_path": true,
	"hosts": null,
	"preserve_host": false,
	"regex_priority": 0,
	"updated_at": 1553907462,
	"paths": [
		"\/v1\/model"
	],
	"service": {
		"id": "495fec4e-8fdf-42cd-b4e9-604b3776cfed"
	},
	"methods": null,
	"protocols": [
		"http"
	],
	"id": "22ea22b4-97d0-4210-aa8e-8e32677969cc"
}
```


### Consumer Object

The Consumer object represents a consumer - or a user - of a Service. You can either rely on Kong as the primary datastore, or you can map the consumer list with your database to keep consistency between Kong and your existing primary datastore.


```
kongctl consumer
NAME:
   kongctl consumer - The kong consumer object.

USAGE:
   kongctl consumer command [command options] [arguments...]

COMMANDS:
     create  create consumer object
     list    list all consumers object
     get     retrieve consumer object
     delete  delete consumer object

OPTIONS:
   --help, -h  show help

```

### Upstream Object

The upstream object represents a virtual hostname and can be used to loadbalance incoming requests over multiple services (targets). So for example an upstream named service.v1.xyz for a Service object whose host is service.v1.xyz. Requests for this Service would be proxied to the targets defined within the upstream.

```
kongctl upstream --h
NAME:
   kongctl upstream - The kong upstream object.

USAGE:
   kongctl upstream command [command options] [arguments...]

COMMANDS:
     create  create upstream object
     get     get upstream object
     list    list all upstream object
     delete  delete upstream object

OPTIONS:
   --help, -h  show help
```

example:
```
kongctl upstream create --name="web"

response body:
{
	"healthchecks": {
		"active": {
			"unhealthy": {
				"http_statuses": [
					429,
					404,
					500,
					501,
					502,
					503,
					504,
					505
				],
				"tcp_failures": 0,
				"timeouts": 0,
				"http_failures": 0,
				"interval": 0
			},
			"http_path": "\/",
			"healthy": {
				"http_statuses": [
					200,
					302
				],
				"interval": 0,
				"successes": 0
			},
			"timeout": 1,
			"concurrency": 10
		},
		"passive": {
			"unhealthy": {
				"http_failures": 0,
				"http_statuses": [
					429,
					500,
					503
				],
				"tcp_failures": 0,
				"timeouts": 0
			},
			"healthy": {
				"http_statuses": [
					200,
					201,
					202,
					203,
					204,
					205,
					206,
					207,
					208,
					226,
					300,
					301,
					302,
					303,
					304,
					305,
					306,
					307,
					308
				],
				"successes": 0
			}
		}
	},
	"created_at": 1553932148565,
	"hash_on": "none",
	"id": "cf9454f1-ee9b-4444-98cc-78ee467b674f",
	"hash_on_cookie_path": "\/",
	"name": "web",
	"hash_fallback": "none",
	"slots": 1000
}
```


### Target Object

A target is an ip address/hostname with a port that identifies an instance of a backend service. Every upstream can have many targets, and the targets can be dynamically added. Changes are effectuated on the fly.


```
kongctl target --h
NAME:
   kongctl target - The kong target object.

USAGE:
   kongctl target command [command options] [arguments...]

COMMANDS:
     create  Create target object
     list    Lists all targets currently active on the upstream’s load balancing wheel
     delete  Disable a target in the load balancer

OPTIONS:
   --help, -h  show help
```

example:

```
kongctl target create --upstream_id=c3cdc9e6-135e-4f6e-9249-791a7ccc9270 --target=xx.xxx.xxx.xx:8080

response body:
{
	"created_at": 1553934934921,
	"weight": 100,
	"upstream_id": "c3cdc9e6-135e-4f6e-9249-791a7ccc9270",
	"target": "xx.xxx.xxx.xx:8080",
	"id": "badf8d0d-d9c4-432d-b784-ce38fbbbff7b"
}
```

### Plugin Object

A Plugin entity represents a plugin configuration that will be executed during the HTTP request/response lifecycle. It is how you can add functionalities to Services that run behind Kong, like Authentication or Rate Limiting for example. You can find more information about how to install and what values each plugin takes by visiting the Kong Hub.


```
kongctl plugin
NAME:
   kongctl plugin - The kong plugin object.

USAGE:
   kongctl plugin command [command options] [arguments...]

COMMANDS:
     avalible_plugins  current support plugins object
     create            create a plugin object
     get               retrieve a plugin object
     list              list all plugins object
     delete            delete a plugin object

OPTIONS:
   --help, -h  show help
```

```
kongctl plugin create --h
NAME:
   kongctl plugin create - create a plugin object

USAGE:
   kongctl plugin create command [command options] [arguments...]

COMMANDS:
     basic-auth  create basic-auth plugin
     statsd      log metrics for a service, route to a StatsD server

OPTIONS:
   --help, -h  show help
```





