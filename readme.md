# go-gedcom

go-gedcom is a Go library for parsing GEDCOM files. It was created primarily for use in manipulating 
[GEDCOM v7 files](https://www.familysearch.org/developers/docs/guides/gedcom7_0). The intent is to provide a
simple way to validate, batch edit, and otherwise process genealogical data.

## Description

This package is currently compatible with GEDCOM v7.0.11. ABNF grammars are imported when the `gedcom7` 
package is initialized at runtime. So long as there are no breaking changes to the specifications, newer
versions of the GEDCOM standard should be compatible with this package. Potential breaking changes include

* Changes to the ABNF grammar core data types

There are a few parts to this repository:

1. CLI tool which uses the packages to provide functionality for manipulating GEDCOM files.
2. Core packages to parse and manipulate GEDCOM files.
   - `gedcom` - Interface definitions for gedcom documents and record structures.
   - `gedcom7` - Implementation of the gedcom interface for GEDCOM v7 files.
     - `gedcom7/gc70val` - Validation rules for GEDCOM v7 files derived from the published abnf grammar.
3. Supporting packages
    - `abnf` - A parser for ABNF grammars.
    - `stack` - A fork of a useful stack implementation used to parse GEDCOM files.
    - `uuid` - A client wrapper to generate UUIDs for GEDCOM records.

## Getting Started

First, get the package:

```
git clone ...
```

Modify the existing cli tool in `/cli` or create your own.

### Dependencies

* go 1.18+

## Help

## Authors

[Bill Shaw](https://github.com/funwithbots)

## Contributing

To contribute to this project, fork the repository, create a new branch, and submit a pull request.

## Roadmap

1. Support for GEDCOM v5.5.1
2. Converter to/from GEDCOM v5.5.1 and v7
3. Split tool to export a subset of the GEDCOM file.
4. Graphing tool to allow batch creation of typical genealogical charts.
5. Validation for documentation structure.
6. Reporting tool for create reports from GEDCOM files.

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under the [GNU General Public License v3.0] - see the LICENSE file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [FamilySearch GEDCOM](https://github.com/FamilySearch/GEDCOM)
* [Caleb Doxsey Stack package](https://github.com/golang-collections/collections)
* [https://github.com/elimity-com/abnf](https://github.com/elimity-com/abnf)