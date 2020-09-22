package main

import (
  "context"
  "fmt"
  "flag"
  "regexp"

  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
)

func main() {

  dryptr := flag.Bool("d", false, "dry run")
  forceptr := flag.Bool("f", false, "force")
  flag.Parse()

  args := flag.Args()
  if len(args) < 1 {
    fmt.Println("missing argument - images regex")
    return
  }
  pattern := args[0]
  fmt.Printf("matching pattern: %s\n", pattern)

  if *dryptr {
    fmt.Println("dry run only")
  }

  ctx := context.Background()
  cli, err := client.NewEnvClient()
  if err != nil {
    panic(err)
  }
  defer cli.Close()
  lst, err := cli.ImageList(ctx, types.ImageListOptions {})
  if err != nil {
    panic(err)
  }
  for _, image := range(lst) {
    for _, tag := range(image.RepoTags) {
      matched, err := regexp.MatchString(pattern, tag)
      if err != nil {
        panic(err)
      }
      if matched {
        fmt.Printf("%s %s\n", tag, image.ID)
        if !*dryptr {
          options := types.ImageRemoveOptions{
            PruneChildren: true,
            Force: *forceptr,
          }
          _, err := cli.ImageRemove(ctx, image.ID, options)
          if err != nil {
            fmt.Println(err)
          }
        }
      }
    }
  }
}