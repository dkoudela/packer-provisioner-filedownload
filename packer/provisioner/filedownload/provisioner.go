package filedownload

import (
  "github.com/mitchellh/packer/common"
  "github.com/mitchellh/packer/packer"
  "os"
  "log"
  "fmt"
  "errors"
)

type config struct {
  common.PackerConfig `mapstructure:",squash"`

  // The remote path of the file to download.
  Source string

  // The local path where the local file will be downloaded to.
  Destination string

  tpl *packer.ConfigTemplate
}

type Provisioner struct {
  config config
}

func (p *Provisioner) Prepare(raws ...interface{}) error {
  md, err := common.DecodeConfig(&p.config, raws...)
    if err != nil {
    return err
  }

  p.config.tpl, err = packer.NewConfigTemplate()
  if err != nil {
    return err
  }
  p.config.tpl.UserVars = p.config.PackerUserVars

  // Accumulate any errors
  errs := common.CheckUnusedConfig(md)

  templates := map[string]*string{
    "source": &p.config.Source,
    "destination": &p.config.Destination,
  }

  for n, ptr := range templates {
    var err error
    *ptr, err = p.config.tpl.Process(*ptr, nil)
    if err != nil {
      errs = packer.MultiErrorAppend(
        errs, fmt.Errorf("Error processing %s: %s", n, err))
    }
  }

  if p.config.Source == "" {
    errs = packer.MultiErrorAppend(errs,
      errors.New("Filedownload: Source must be specified."))
  }

  log.Println(fmt.Sprintf("Filedownload: source: %s, destination: %s", p.config.Source, p.config.Destination))

  if errs != nil && len(errs.Errors) > 0 {
    return errs
  }
//  return fmt.Errorf("Stop here for now !!! !!!")
  return nil
}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error {
  log.Println(fmt.Sprintf("Downloading %s => %s", p.config.Source, p.config.Destination))
  ui.Say(fmt.Sprintf("Downloading %s => %s", p.config.Source, p.config.Destination))

  // We're downloading a file...
  log.Println(fmt.Sprintf("Opening the target file: %s", p.config.Destination))
  f, err := os.OpenFile(p.config.Destination, os.O_WRONLY|os.O_CREATE, 0666)
  if err != nil {
    log.Println(fmt.Sprintf("Opening the target file failed: %s", err))
    return err
  }
  defer f.Close()

  err = comm.Download(p.config.Source, f)
  if err != nil {
    ui.Error(fmt.Sprintf("Download failed: %s", err))
  }
  if err != nil { // TODO: Remove it later on
    return err
  }
  return fmt.Errorf("Stop here for now !!! !!!")
//  return nil
}

func (p *Provisioner) Cancel() {
  // Just hard quit. It isn't a big deal if what we're doing keeps
  // running on the other side.
  os.Exit(0)
}
