siteshot
========
![MIT license](https://img.shields.io/github/license/writeas/siteshot.svg)

Website screenshot capturer written in Go and designed to run without an X session (thanks to [Xvfb](http://www.x.org/archive/X11R7.6/doc/man/man1/Xvfb.1.xhtml)).

## Dependencies

* Python
* [ImageMagick](http://www.imagemagick.org/)
* [webkit2png](https://github.com/adamn/python-webkit2png)
* [Xvfb](http://www.x.org/archive/X11R7.6/doc/man/man1/Xvfb.1.xhtml)

## Usage

Send a POST request with a `url` field to the server. Any other request returns `400 Bad Request`.

`curl --data "url=https://html.house/xpeoccu2.html" http://localhost:3333`

## Installation

Tested on Ubuntu 14.04 with Python 2.7 / pip 1.5.4.

**TL;DR** `sudo apt-get install xvfb imagemagick python-qt4 libqt4-webkit && pip install git+https://github.com/adamn/python-webkit2png.git && go get github.com/writeas/siteshot`

1. Get ImageMagick: `sudo apt-get install imagemagick`
2. Get Xvfb: `sudo apt-get install xvfb`
3. Get webkit2png dependencies: `sudo apt-get install python-qt4 libqt4-webkit`
4. Get webkit2png: `pip install git+https://github.com/adamn/python-webkit2png.git` (from a [GitHub comment](https://github.com/adamn/python-webkit2png/issues/55#issuecomment-150974776))
5. Get siteshot: `go get github.com/writeas/siteshot`
