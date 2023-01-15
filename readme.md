# go-gedcom

go-gedcom is a Go library for parsing GEDCOM files. It was created primarily for use in manipulating 
[GEDCOM v7 files](https://www.familysearch.org/developers/docs/guides/gedcom7_0). The intent is to provide a
simple way to validate, batch edit, and otherwise process genealogical data.

## Description

This package is currently compatible with GEDCOM v7.0.11. ABNF grammars are loaded when the `gedcom7` 
package is initialized at runtime. So long as there are no breaking changes to the specifications, newer
versions of the GEDCOM standard should be compatible with this package. Potential breaking changes include

* Changes to the ABNF grammar core data types

There are a few parts to this repository:

1. CLI tool which uses the packages to provide functionality for manipulating GEDCOM files.
2. Core packages to parse and manipulate GEDCOM files.
   - `gedcom` 
     - Interface definition for gedcom documents.
     - node.go provides the tree structure for records.
   - `gedcom7` - Implementation of the gedcom interface for GEDCOM v7 files.
     - `gedcom7/gc70val` - Validation rules for GEDCOM v7 files derived from the published abnf grammar.
3. Supporting packages
    - `abnf` - A parser for ABNF grammars.
    - `stack` - A fork of a useful stack implementation used to parse GEDCOM files.
    - `uuid` - A client wrapper to generate UUIDs for GEDCOM records.

## Getting Started

First, get the package:

```bash
git clone https://github.com/funwithbots/go-gedcom.git
```

Modify the existing cli tool in `/cli` or create your own.

### Dependencies

* go 1.18+

## Help

## Authors

[Bill Shaw](https://github.com/funwithbots)

## Contributing

To contribute to this project, fork the repository, create a new branch, and submit a pull request.

At a minimum, pull requests must pass the following checks:
- all tests must pass
- linting must pass using golangci-lint

Workflows run these checks automatically when branches are pushed.

## Roadmap

1. Support for GEDCOM v5.5.1
2. Converter to/from GEDCOM v5.5.1 and v7
3. Split tool to export a subset of the GEDCOM file.
4. Graphing tool to allow batch creation of typical genealogical charts.
5. Validation for GEDCOM file structures.
6. Reporting tool for creating reports from GEDCOM files.
7. Full Record comparison tool

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under the [GNU General Public License v3.0] - see the LICENSE file for details

## RELEASE OF LIABILITY

This software is provided AS-IS without any warranty of any kind, either expressed or implied, including, but not limited to, the implied warranties of merchantability and fitness for a particular purpose. The entire risk as to the quality and performance of the software is with you. Should the software prove defective, you assume the cost of all necessary servicing, repair, or correction.

## Acknowledgments

Inspiration, code snippets, etc.
* [FamilySearch GEDCOM](https://github.com/FamilySearch/GEDCOM)
* [Caleb Doxsey Stack package](https://github.com/golang-collections/collections)
* [https://github.com/elimity-com/abnf](https://github.com/elimity-com/abnf)