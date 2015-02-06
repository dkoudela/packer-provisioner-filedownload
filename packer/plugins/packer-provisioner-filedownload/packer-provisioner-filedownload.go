package main

import (
  "github.com/mitchellh/packer/packer/plugin"
  "github.com/dkoudela/packer-provisioner-filedownload/packer/provisioner/filedownload"
)

func main() {
  server, err := plugin.Server()
  if err != nil {
    panic(err)
  }
  server.RegisterProvisioner(new(filedownload.Provisioner))
  server.Serve()
}
