package ethapi

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/rollup"
	"github.com/ethereum/go-ethereum/rpc"
)

type Utils struct{}

func (f *Utils) SaveString(data string, fileName string) error {
	if _, err := os.Stat(fileName); !errors.Is(err, os.ErrNotExist) { //over write file is not allowed, avoid rewrite
		// important file by accident
		return errors.New("duplicated file exist")
	}
	err := ioutil.WriteFile(fileName, []byte(data), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func UtilsApis(backend *rollup.RollupBackend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "utils",
			Version:   "1.0",
			Service:   &Utils{},
			Public:    false,
		},
	}
}
