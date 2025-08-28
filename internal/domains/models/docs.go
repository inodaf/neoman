package models

import "errors"

type Docs struct {
  Owner             Owner
  Name              string
  GitRemoteProvider string
}

type Owner struct {
  Value string
}

func NewOwner(name string) (*Owner, error) {
  if len(name) <= 0 {
    return nil, errors.New("asd")
  }
  return &Owner{ Value: name }, nil
}

func run() {
  owner, err := NewOwner("isac")
  if err != nil {
    panic("")
  }

  docs := Docs{
    Owner: *owner,
  }
}