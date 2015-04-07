## Purpose
XMLGen is a tool for generating native Golang types from XML. This automates what is otherwise a very tedious and error prone task when working with XML.

[![Build Status](http://img.shields.io/travis/dutchcoders/XMLGen.svg?style=flat)](https://travis-ci.org/dutchcoders/XMLGen)
[![GPLv3 License](http://img.shields.io/badge/license-GPLv3-blue.svg?style=flat)](http://choosealicense.com/licenses/gpl-3.0/)

XMLGen is based on and using code from [JSONGen](https://github.com/bemasher/JSONGen).

## Install (using brew)

```
brew tap dutchcoders/homebrew-xmlgen
brew install xmlgen
```

## Usage

```
$ xmlgen -h
Usage of xmlgen:
  -dump="NUL": Dump tree structure to file.
  -normalize=true: Squash arrays of struct and determine primitive array type.
  -title=true: Convert identifiers to title case, treating '_' and '-' as word boundaries.
```

Reading from stdin can be done as follows:
```
$ cat test.xml | xmlgen
```

Or a filename can be passed:
```
$ xmlgen test.xml
```

Using [test.xml](test.xml) as input the example will produce:
```go
➜ xmlgen git:(master) ✗ curl http://www.ibiblio.org/xml/examples/shakespeare/all_well.xml | xmlgen
type _ struct {
        PLAY struct {
                TITLE string `xml:"TITLE"`
                FM    struct {
                        P []string `xml:"P"`
                } `xml:"FM"`
                PERSONAE struct {
                        TITLE   string   `xml:"TITLE"`
                        PERSONA []string `xml:"PERSONA"`
                        PGROUP  struct {
                                PERSONA  []string `xml:"PERSONA"`
                                GRPDESCR string   `xml:"GRPDESCR"`
                        } `xml:"PGROUP"`
                        PERSONA []string `xml:"PERSONA"`
                        PGROUP  struct {
                                PERSONA  []string `xml:"PERSONA"`
                                GRPDESCR string   `xml:"GRPDESCR"`
                        } `xml:"PGROUP"`
                        PERSONA string `xml:"PERSONA"`
                } `xml:"PERSONAE"`
                SCNDESCR string `xml:"SCNDESCR"`
                PLAYSUBT string `xml:"PLAYSUBT"`
                ACT      []struct {
                        TITLE string `xml:"TITLE"`
                        SCENE []struct {
                                TITLE    string `xml:"TITLE"`
                                STAGEDIR string `xml:"STAGEDIR"`
                                SPEECH   []struct {
                                        SPEAKER  string `xml:"SPEAKER"`
                                        LINE     string `xml:"LINE"`
                                        STAGEDIR string `xml:"STAGEDIR"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                                SPEECH   []struct {
                                        SPEAKER  string   `xml:"SPEAKER"`
                                        LINE     []string `xml:"LINE"`
                                        STAGEDIR string   `xml:"STAGEDIR"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                                SPEECH   []struct {
                                        SPEAKER  string   `xml:"SPEAKER"`
                                        LINE     []string `xml:"LINE"`
                                        STAGEDIR []string `xml:"STAGEDIR"`
                                        LINE     []string `xml:"LINE"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                                SPEECH   struct {
                                        SPEAKER  string `xml:"SPEAKER"`
                                        LINE     string `xml:"LINE"`
                                        STAGEDIR string `xml:"STAGEDIR"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                                SPEECH   []struct {
                                        SPEAKER  string   `xml:"SPEAKER"`
                                        LINE     []string `xml:"LINE"`
                                        STAGEDIR string   `xml:"STAGEDIR"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                                SPEECH   struct {
                                        SPEAKER  string   `xml:"SPEAKER"`
                                        LINE     []string `xml:"LINE"`
                                        STAGEDIR string   `xml:"STAGEDIR"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                        } `xml:"SCENE"`
                        EPILOGUE struct {
                                TITLE  string `xml:"TITLE"`
                                SPEECH struct {
                                        SPEAKER string   `xml:"SPEAKER"`
                                        LINE    []string `xml:"LINE"`
                                } `xml:"SPEECH"`
                                STAGEDIR string `xml:"STAGEDIR"`
                        } `xml:"EPILOGUE"`
                } `xml:"ACT"`
        } `xml:"PLAY"`
}
```

## Parsing
### Field Names
  * Field names are sanitized and written as exported fields of the generated type.
  * If sanitizing produces an empty string the identifier is changed to `_`, this will need to be set by hand in order to properly decode the type.
  * Spaces and `-` are converted to `_`.
  * Field names are converted to title case treating `_` and `-` as word boundaries along with spaces. This can be disabled using `-title=false`.

## Types
### Primitive
  * Primitive types are parsed and stored as-is.
  * Valid types are currently only strings

### Object
  * Object types are treated as structs.
  * Fields of structures are sorted lexicographically by sanitized field name.

### Lists
  * A homogeneous list of recurring elements with the same name the primitive type e.g.: `[]string`
  * Lists with object elements are treated as a list of structs.
    * Fields of each element are "squashed" into a single struct. The result is an array of a struct containing all encountered fields.   

Examples of all of the above can be found in [test.xml](test.xml).

## Caveats
  * Currently sibling field names are not guaranteed to be unique.

## License
The source of this project is licensed under GNU GPL v3.0, according to [http://choosealicense.com/licenses/gpl-3.0/](http://choosealicense.com/licenses/gpl-3.0/):

#### Required:

 * **Disclose Source:** Source code must be made available when distributing the software. In the case of LGPL, the source for the library (and not the entire program) must be made available.
 * **License and copyright notice:** Include a copy of the license and copyright notice with the code.
 * **State Changes:** Indicate significant changes made to the code.

#### Permitted:

 * **Commercial Use:** This software and derivatives may be used for commercial purposes.
 * **Distribution:** You may distribute this software.
 * **Modification:** This software may be modified.
 * **Patent Grant:** This license provides an express grant of patent rights from the contributor to the recipient.
 * **Private Use:** You may use and modify the software without distributing it.

#### Forbidden:

 * **Hold Liable:** Software is provided without warranty and the software author/license owner cannot be held liable for damages.
 * **Sublicensing:** You may not grant a sublicense to modify and distribute this software to third parties not included in the license.

## Feedback
If you find a case that produces incorrect results or you have a feature suggestion, let me know: submit an issue.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

Douglas Hall is the creator of JSONGen. XMLGen is based and using code from [JSONGen](https://github.com/bemasher/JSONGen).

