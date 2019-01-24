package tests

import (
   "testing"
   "os/user"
   "github.com/temp25/hdl/helper"
)

func TestHomeDir(t *testing.T) {
  user, err := user.Current()
   if err != nil {
      panic(err)
   }
   expectedHomeDirectoryPath := "/home/" + user.Username
   actualHomeDirectoryPath, _ := helper.HomeDir()
   if expectedHomeDirectoryPath != actualHomeDirectoryPath {
      t.Error("Expected",expectedHomeDirectoryPath , " but got", actualHomeDirectoryPath)
   }
}

