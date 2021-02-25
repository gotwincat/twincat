<h1 align="center">Twincat</h1>

A native Go implementation of the [Beckhoff Twincat V3 protocol](https://infosys.beckhoff.com/english.php?content=../content/1033/tc3_adsnetref/7312567947.html&id=).

[![CircleCI](https://circleci.com/gh/gotwincat/twincat.svg?style=shield)](https://circleci.com/gh/gotwincat/twincat)
[![GoDoc](https://godoc.org/github.com/gotwincat/twincat?status.svg)](https://godoc.org/github.com/gotwincat/twincat)
[![GolangCI](https://golangci.com/badges/github.com/gotwincat/twincat.svg)](https://golangci.com/r/github.com/gotwincat/twincat)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/gotwincat/twincat/blob/main/LICENSE)
[![Version](https://img.shields.io/github/tag/gotwincat/twincat.svg?color=blue&label=version)](https://github.com/gotwincat/twincat/releases)

You need go1.16 or higher. We test with the current and previous Go version.

<table>
   <tr>
      <td width="60%" align="center">
         <img width="25%" src="https://github.com/gotwincat/twincat/blob/main/gopher.png">
      </td>
      <td width="40%">
        Artwork by <a href="https://twitter.com/ashleymcnamara">Ashley McNamara</a><br/>
        Inspired by <a href="http://reneefrench.blogspot.co.uk/">Renee French</a><br/>
        Taken from <a href="https://gopherize.me">https://gopherize.me</a> by <a href="https://twitter.com/matryer">Mat Ryer</a>
      </td>
   </tr>
</table>

## Quickstart

```sh
go get -u github.com/gotwincat/twincat
```

## Sponsors

The `gotwincat` project is sponsored by the following organizations by supporting the active committers to the project:

<p align="center">
  <a href="https://northvolt.com/">
    <img alt="Northvolt" width="50%" src="https://github.com/gotwincat/twincat/blob/main/logo/northvolt.png">
  </a>
</p>

### Users

We would also like to list organizations which use `gotwincat` in production. Please open a PR to include your logo below.

## Authors

The [gotwincat Team](https://github.com/gotwincat/twincat/graphs/contributors).

If you need to get in touch with us directly you may find us on [Keybase.io](https://keybase.io)
but try to create an issue first.

## Supported Features

| Request            | Supported | Notes |
|--------------------|-----------|-------|
| Read               | Yes       |       |
| ReadWrite          | Yes       |       |
| Write              | Yes       |       |
| GetSymHandleByName | Yes       |       |

## License

[MIT](https://github.com/gotwincat/twincat/blob/main/LICENSE)
