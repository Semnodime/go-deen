package hashs

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/takeshixx/deen/pkg/types"
	"golang.org/x/crypto/sha3"
)

// NewPluginSHA3224 creates a plugin
func NewPluginSHA3224() (p types.DeenPlugin) {
	p.Name = "sha3-224"
	p.Aliases = []string{}
	p.Type = "hash"
	p.Unprocess = false
	p.ProcessStreamFunc = func(reader io.Reader) ([]byte, error) {
		var err error
		hasher := sha3.New224()
		if _, err := io.Copy(hasher, reader); err != nil {
			return *new([]byte), err
		}
		hashSum := hasher.Sum(nil)
		outBuf := make([]byte, hex.EncodedLen(len(hashSum[:])))
		_ = hex.Encode(outBuf, hashSum[:])
		return outBuf, err
	}
	p.ProcessStreamWithCliFlagsFunc = func(flags *flag.FlagSet, reader io.Reader) ([]byte, error) {
		return p.ProcessStreamFunc(reader)
	}
	p.AddCliOptionsFunc = func(self *types.DeenPlugin, args []string) *flag.FlagSet {
		sha3Cmd := flag.NewFlagSet(p.Name, flag.ExitOnError)
		sha3Cmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s: \n\n", p.Name)
			fmt.Fprintf(os.Stderr, "SHA3 is the latest member of the Secure Hash Algorithm family\nof standards, released by NIST.\n\n")
			sha3Cmd.PrintDefaults()
		}
		sha3Cmd.Parse(args)
		return sha3Cmd
	}
	return
}

// NewPluginSHA3256 creates a plugin
func NewPluginSHA3256() (p types.DeenPlugin) {
	p.Name = "sha3-256"
	p.Aliases = []string{}
	p.Type = "hash"
	p.Unprocess = false
	p.ProcessStreamFunc = func(reader io.Reader) ([]byte, error) {
		var err error
		hasher := sha3.New256()
		if _, err := io.Copy(hasher, reader); err != nil {
			return *new([]byte), err
		}
		hashSum := hasher.Sum(nil)
		outBuf := make([]byte, hex.EncodedLen(len(hashSum[:])))
		_ = hex.Encode(outBuf, hashSum[:])
		return outBuf, err
	}
	p.ProcessStreamWithCliFlagsFunc = func(flags *flag.FlagSet, reader io.Reader) ([]byte, error) {
		return p.ProcessStreamFunc(reader)
	}
	p.AddCliOptionsFunc = func(self *types.DeenPlugin, args []string) *flag.FlagSet {
		sha3Cmd := flag.NewFlagSet(p.Name, flag.ExitOnError)
		sha3Cmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s: \n\n", p.Name)
			fmt.Fprintf(os.Stderr, "SHA3 is the latest member of the Secure Hash Algorithm family\nof standards, released by NIST.\n\n")
			sha3Cmd.PrintDefaults()
		}
		sha3Cmd.Parse(args)
		return sha3Cmd
	}
	return
}

// NewPluginSHA3384 creates a plugin
func NewPluginSHA3384() (p types.DeenPlugin) {
	p.Name = "sha3-384"
	p.Aliases = []string{}
	p.Type = "hash"
	p.Unprocess = false
	p.ProcessStreamFunc = func(reader io.Reader) ([]byte, error) {
		var err error
		hasher := sha3.New384()
		if _, err := io.Copy(hasher, reader); err != nil {
			return *new([]byte), err
		}
		hashSum := hasher.Sum(nil)
		outBuf := make([]byte, hex.EncodedLen(len(hashSum[:])))
		_ = hex.Encode(outBuf, hashSum[:])
		return outBuf, err
	}
	p.ProcessStreamWithCliFlagsFunc = func(flags *flag.FlagSet, reader io.Reader) ([]byte, error) {
		return p.ProcessStreamFunc(reader)
	}
	p.AddCliOptionsFunc = func(self *types.DeenPlugin, args []string) *flag.FlagSet {
		sha3Cmd := flag.NewFlagSet(p.Name, flag.ExitOnError)
		sha3Cmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s: \n\n", p.Name)
			fmt.Fprintf(os.Stderr, "SHA3 is the latest member of the Secure Hash Algorithm family\nof standards, released by NIST.\n\n")
			sha3Cmd.PrintDefaults()
		}
		sha3Cmd.Parse(args)
		return sha3Cmd
	}
	return
}

// NewPluginSHA3512 creates a plugin
func NewPluginSHA3512() (p types.DeenPlugin) {
	p.Name = "sha3-512"
	p.Aliases = []string{}
	p.Type = "hash"
	p.Unprocess = false
	p.ProcessStreamFunc = func(reader io.Reader) ([]byte, error) {
		var err error
		hasher := sha3.New512()
		if _, err := io.Copy(hasher, reader); err != nil {
			return *new([]byte), err
		}
		hashSum := hasher.Sum(nil)
		outBuf := make([]byte, hex.EncodedLen(len(hashSum[:])))
		_ = hex.Encode(outBuf, hashSum[:])
		return outBuf, err
	}
	p.ProcessStreamWithCliFlagsFunc = func(flags *flag.FlagSet, reader io.Reader) ([]byte, error) {
		return p.ProcessStreamFunc(reader)
	}
	p.AddCliOptionsFunc = func(self *types.DeenPlugin, args []string) *flag.FlagSet {
		sha3Cmd := flag.NewFlagSet(p.Name, flag.ExitOnError)
		sha3Cmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s: \n\n", p.Name)
			fmt.Fprintf(os.Stderr, "SHA3 is the latest member of the Secure Hash Algorithm family\nof standards, released by NIST.\n\n")
			sha3Cmd.PrintDefaults()
		}
		sha3Cmd.Parse(args)
		return sha3Cmd
	}
	return
}