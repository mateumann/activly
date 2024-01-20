package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	src := client.Host().Directory(".", dagger.HostDirectoryOpts{Exclude: []string{"build/", "ci/"}})

	golang := client.Container().From("golang:1.21")

	golang = golang.WithDirectory("/src", src).WithWorkdir("/src")

	path := "build/"

	golang = golang.WithExec([]string{"go", "build", "-o", path + "/activly", "./cmd/main.go"})

	// get reference to build output directory in container
	output := golang.Directory(path)

	// write contents of container build/ directory to the host
	out, err := output.Export(ctx, path)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
