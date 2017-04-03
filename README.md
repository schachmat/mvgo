**mvgo** is a command line tool to list public transport departure times.

## Features and limitations

* Requires Go 1.5
* Automatic config management with [ingo](https://github.com/schachmat/ingo)
* Supported backends:
  * [mvg live](https://www.mvg-live.de/)
  * [mvv efa](http://efa.mvv-muenchen.de/)
  * merger: Merges input from two different backends
* Supported frontends:
  * ascii table: Lists the departure times in rows
  * json: Can be useful to be consumed by other tools
  * template: Builds a html page from a template usable for live ticker screens

## Installation

```shell
go get -u github.com/schachmat/mvgo
```

## Usage example

```shell
$ mvgo "Münchner Freiheit"
The next departures from Münchner Freiheit are:
  0  59    Ackermannbogen
  0  59    Giesing Bf.
  1  142   Scheidplatz
  1  U6    Garching-Forschungszentrum
  1  54    Lorettoplatz
  3  U3    Moosach
  3  U6    Harras
  6  53    Aidenbachstraße
  7  U6    Klinikum Großhadern
  9  23    Schwabing Nord
```

A default station can also be saved in the `~/.mvgorc` configuration file. The
config path is customizable with the `MVGORC` environment variable.

## Contributing

If you want your local public transport provider to be supported, write your own
backend and make a Pull Request. Ideas for new frontends are also welcome.

## License - ISC

Copyright (c) 2016-2017,  <teichm@in.tum.de>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
