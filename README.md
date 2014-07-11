# bikeshed

* Computes placeholder colors for artistic images (movie posters, DVD boxes etc.)
* Requires imagemagick, works with the old imagemagick that you can get for Ubuntu 12.04 LTS (useful when running this in Amazon EC2)

## Installation
* Install imagemagick. The old one in Ubuntu 12.04 LTS will do just fine.
* Get the [binary](https://dl.dropboxusercontent.com/u/2098438/Permanent/bikeshed/bin/bikeshed) (64 bits) or [compile it yourself](http://golang.org/doc/code.html#remote).

## Usage:

bikeshed imagefilename

## Output:

RGB as hex without alpha. E.g. c0ffee
