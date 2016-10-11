# chefviz

## Description

Chefviz creates the dot files of recipes dependency-graph for graphviz.


## Usage

```
chefviz [--rootdir /path/to/chef-directory] cookbook::recipe
```

The rootdir option can specify both absolute and relative path.

For example,,,

```
$ chefviz --rootdir /path/to/sample-chef-repo nginx::default > test.dot
$ cat test.dot
digraph G {
	"nginx::default"->"nginx::ngx_lua_module";
	"nginx::ngx_lua_module"->"nginx::lua";
	"nginx::ngx_lua_module"->"nginx::ngx_devel_module";
	"nginx::default";
	"nginx::lua";
	"nginx::ngx_devel_module";
	"nginx::ngx_lua_module";

}

# Then you can get the graph by dot command(graphviz).
$ dot -Tpng test.dot -o test.png
```

![top-page](https://github.com/tom--bo/chefviz/raw/image/nginx-sample.png)

## Install

To install, use `go get`:

```bash
$ go get github.com/tom--bo/chefviz
```

## Contribution

1. Fork ([https://github.com/tom--bo/chefviz/fork](https://github.com/tom--bo/chefviz/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[tom--bo](https://github.com/tom--bo)
