package config

import (
	"bytes"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var fname = "config.yaml"
var lastCfg = []byte{}

func CheckChange() []byte {
	bs, err := os.ReadFile(fname)
	if err != nil {
		return nil
	}
	if bytes.Equal(lastCfg, bs) {
		return nil
	}
	return bs
}

func LoadOnce(i interface{}) error {
	bs, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	return LoadFromBytes(i, bs)
}

func LoadFromBytes(i interface{}, bs []byte) error {
	err := yaml.Unmarshal(bs, i)
	if err == nil {
		lastCfg = bs
	}
	return err
}

func Monitor(i interface{}) chan struct{} {
	ret := make(chan struct{})
	LoadOnce(i)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			bs := CheckChange()
			if bs != nil {
				err := LoadFromBytes(i, bs)
				if err != nil {
					ret <- struct{}{}
				}
			}
		}
	}()
	return ret
}
