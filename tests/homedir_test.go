package tests

import (
   "testing"
   "github.com/temp25/hdl/helper"
)

func TestHomeDir(t *testing.T) {
   expectedHomeDirectoryPath := "/home/user"
   actualHomeDirectoryPath, _ := helper.HomeDir()
   if expectedHomeDirectoryPath != actualHomeDirectoryPath {
      t.Error("Expected",expectedHomeDirectoryPath , " but got", actualHomeDirectoryPath)
   }
}
