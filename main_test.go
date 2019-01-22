package main

import "testing"

func TestMain(t *testing.T) {
  var v float64
  v = 1.5 //Average([]float64{1,2})
  if v != 1.5 {
    t.Error("Expected 1.5, got ", v)
  }
}