ao
==

Await order and execute it.

```
[Terminal A]
$ ao await
1..5
1
2
3
4
5
(keep waiting for requests)

[Terminal B]
$ ao order echo '1..5'
(request to execute "echo '1..5'")

$ ao order seq 5
(request to execute "seq 5")
```

Usage
-----

```
$ ao <command> [...]
await order and execute it.

commands:
  ao a|await [-port=<port>]                  # await order
  ao o|order [-port=<port>] <cmd> [<arg(s)>] # order to execute command
  ao h|help                                  # display usage
  ao v|version                               # display version
```

Installation
------------

### go get

```
go get github.com/nil2nekoni/ao
```

Commands
--------

### a, await

Await order and execute it.

By default, `ao await` accepts requests on port 60080.
To change the port, specify it with `-port`.

```
$ ao a
(accepts requests on port 60080)

$ ao a -port 50180
(accepts requests on port 50180)
```

### o, order

Order to execute command.

By default, `ao order` send a request to port 60080.
To change the port, specify it with `-port`.

```
$ ao o seq 5
(send a request to execute "seq 5" to 60080 port)

$ ao o -port 50180 seq 5
(send o request to execute "seq 5" to 50180 port)
```

### h, help

Display usage.

### v, version

Display version information.

Specification
-------------

### Request

`ao` exchange requests with JSON.
The structure is as follows:

```
{
  "cmd": ["<command>", "<arg-1>", "<arg-2>", ... "<arg-n>"]
}
```

You can also send a request by using curl as follows:

```
$ curl -H 'Content-Type: application/json' -d '{"cmd":["ls"]}' localhost:60080
```

License
-------

MIT License

Author
------

nil2 <nil2@nil2.org>
