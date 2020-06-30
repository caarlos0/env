# env

[![Build Status](https://img.shields.io/github/workflow/status/caarlos0/env/build?style=for-the-badge)](https://github.com/caarlos0/env/actions?workflow=build)
[![Coverage Status](https://img.shields.io/codecov/c/gh/caarlos0/env.svg?logo=codecov&style=for-the-badge)](https://codecov.io/gh/caarlos0/env)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](http://godoc.org/github.com/caarlos0/env)

Simple lib to parse envs to structs in Go.

## Example

A very basic example:

```go
package main

import (
	"fmt"
	"time"

	// if using go modules
	"github.com/caarlos0/env/v6"

	// if using dep/others
	"github.com/caarlos0/env"
)

type config struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT" envDefault:"3000"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS" envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
}
```

You can run it like this:

```sh
$ PRODUCTION=true HOSTS="host1:host2:host3" DURATION=1s go run main.go
{Home:/your/home Port:3000 IsProduction:true Hosts:[host1 host2 host3] Duration:1s}
```

## Supported types and defaults

Out of the box all built-in types are supported, plus a few others that
are commonly used.

Complete list:

- `string`
- `bool`
- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`
- `float32`
- `float64`
- `string`
- `time.Duration`
- `encoding.TextUnmarshaler`
- `url.URL`

Pointers, slices and slices of pointers of those types are also supported.

You can also use/define a [custom parser func](#custom-parser-funcs) for any
other type you want.

If you set the `envDefault` tag for something, this value will be used in the
case of absence of it in the environment.

By default, slice types will split the environment value on `,`; you can change
this behavior by setting the `envSeparator` tag.

If you set the `envExpand` tag, environment variables (either in `${var}` or
`$var` format) in the string will be replaced according with the actual value
of the variable.

Unexported fields are ignored.

## Custom Parser Funcs

If you have a type that is not supported out of the box by the lib, you are able
to use (or define) and pass custom parsers (and their associated `reflect.Type`)
to the `env.ParseWithFuncs()` function.

In addition to accepting a struct pointer (same as `Parse()`), this function
also accepts a `map[reflect.Type]env.ParserFunc`.

`env` also ships with some pre-built custom parser funcs for common types. You
can check them out [here](parsers/).

If you add a custom parser for, say `Foo`, it will also be used to parse
`*Foo` and `[]Foo` types.

This directory contains pre-built, custom parsers that can be used with `env.ParseWithFuncs`
to facilitate the parsing of envs that are not basic types.

Check the example in the [go doc](http://godoc.org/github.com/caarlos0/env)
for more info.

## Required fields

The `env` tag option `required` (e.g., `env:"tagKey,required"`) can be added
to ensure that some environment variable is set.  In the example above,
an error is returned if the `config` struct is changed to:


```go
type config struct {
    Home         string   `env:"HOME"`
    Port         int      `env:"PORT" envDefault:"3000"`
    IsProduction bool     `env:"PRODUCTION"`
    Hosts        []string `env:"HOSTS" envSeparator:":"`
    SecretKey    string   `env:"SECRET_KEY,required"`
}
```


## From file

The `env` tag option `file` (e.g., `env:"tagKey,file"`) can be added
to in order to indicate that the value of the variable shall be loaded from a file. The path of that file is given
by the environment variable associated with it
Example below

```go
package main

import (
	"fmt"
	"time"
	"github.com/caarlos0/env"
)

type config struct {
	Secret       string   `env:"SECRET,file"`
	Password     string   `env:"PASSWORD,file" envDefault:"/tmp/password"`
	Certificate  string   `env:"CERTIFICATE,file" envDefault:"${CERTIFICATE_FILE}" envExpand:"true"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
}
```

```sh
$ echo qwerty > /tmp/secret
$ echo dvorak > /tmp/password
$ echo coleman > /tmp/certificate

$ SECRET=/tmp/secret  \
	CERTIFICATE_FILE=/tmp/certificate \
	go run main.go
{Secret:qwerty Password:dvorak Certificate:coleman}
```


## Options

### Decryptor

The `env` tag option `decrypt` (e.g., `env:"tagKey,decrypt"`) can be added
to in order to indicate that the value of the variable is encrypted and should be decrypted before set and parsed further. If tag `file` is set the value will be loaded from the file and then decrypted and set.

Having the filename decrypted and then loaded is not supported. It will assume the filename is unencrypted but the content of the file should be decrypted.

Other actions such as `expand` will happen **after** decryption is done if the `file` tag **is not set**. This is due to loading from file is the last action. So no other actions are possible after the file has been loaded.

If the struct contains an `decrypt` tag it must call `Parse` with `Options.Decryptor` set.
Otherwise an error will be returned.

The `Options.Decryptor` is an struct that should implement the `Decrypt(string) (string, error)` function from `Decryptor` interface.

Example below using AWS and KMS to decrypt

```go
package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/caarlos0/env"
)

type Config struct {
	kms     *kms.Client
	envData struct {
		Password string `env:"PASSWORD,decrypt"`
	}
}

// Decrypt will decrypt val using KMS.
func (c *Config) Decrypt(val string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return "", fmt.Errorf("couldn't base64 decode string")
	}

	resp, err := c.kms.DecryptRequest(&kms.DecryptInput{
		CiphertextBlob: raw,
	}).Send(context.Background())
	if err != nil {
		return "", fmt.Errorf("couldn't decrypt string. %s", err.Error())
	}

	return string(resp.DecryptOutput.Plaintext), nil
}

func main() {
	// Load the default AWS config and create KMS service.
	awsCfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg := &Config{kms: kms.New(awsCfg)}
	opts := &env.Options{Decryptor: cfg}

	// Load env vars.
	if err := env.Parse(&cfg.envData, opts); err != nil {
		log.Fatal(err)
	}

	// Print the loaded data.
	fmt.Printf("%+v\n", cfg.envData)
}
```

### Environment

By setting the `Options.Environment` map you can till `Parse` to set those `keys` and `values`
as env vars before parsing is done.

This can make your testing scenarios a bit more clean and easy to handle.

```go
package main

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	Password string `env:"PASSWORD"`
}


func main() {
	cfg := &Config{}
	opts := &env.Options{Environment: map[string]string{
		"PASSWORD": "MY_PASSWORD",
	}}

	// Load env vars.
	if err := env.Parse(&cfg.envData, opts); err != nil {
		log.Fatal(err)
	}

	// Print the loaded data.
	fmt.Printf("%+v\n", cfg.envData)
}
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/env.svg)](https://starchart.cc/caarlos0/env)

