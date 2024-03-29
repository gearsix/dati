
dati
====

data and template interface
 
USAGE
-----

  dati [OPTIONS]

DESCRIPTION
-----------

  dati aims to provide a universal interface for executing data files,
  written in any data-serialization language, against template files,
  written in any templating languages.
  Ideally dati will support any language you want to use.
 
  dati works by using various libraries that do all the hard work to
  parse data and template files passed to it. It generates a data
  structure of all the passed data files combined (a super-data
  structure) and executes that structure against a set of root template
  files.
  The used libraries are listed below for credit/reference.
  
  dati can also be imported as a golang package to be used as a library.
 
OPTIONS
-------

  - **-r**, **-root** *PATH*<br/>
  Path of the root template file to execute against.

  - **-p**, **-partial** *PATH ...*<br/>
  Path of (multiple) template files that are called upon by at least
  one root template
    - If a directory is passed then all files within that directory
	will (recursively) be loaded.

  - **-gd**, **-global-data** *PATH ...*<br/>
  Path of (multiple) data files to load as "global data".
  If a directory is passed then all files within that directory will
  (recursively) be loaded.

  - **-dk**, **-data-key** *NAME*<br/>
  Set the name of the key used for the generated array of data. The
  default *data key* is "data".

  - **-sd**, **-sort-data** *ATTRIBUTE*<br/>
  The file attribute to order data files by. If no value is provided,
  the data will be provided in the order it's loaded.
    - *Accepted values*: "filename", "modified".
    - A suffix can be appended to each value to set the sort order:
	"-asc" (for ascending), "-desc" (for descending).
	If not specified, this defaults to "-asc".

  - **-cfg** **-config** *FILE*<br/>
  A data file to provide default values for the above options (CONFIG).

CONFIG
------

  It's possible you'll want to set the same options if you run dati multiple
  times for the same project. This can be done by creating a file (written as
  a data file) and passing the filepath to the -cfg argument.
  
  The key names for the options set in the config file must match the name of
  the argument option to set (long or short). For example (a config file in
  toml):

	root="~/templates/blog.mst"
	partial="~/templates/blog/"
	gd="./blog.json"
	data="./posts/"
	dk="posts"

DATA
----

  dati generates a single super-structure of all the data files passed to it.
  This super-structure is executed against each "root" template.

  The super-structure generated by dati will only have 1 definite key: "data"
  (or the value of the "data-key" option). This key will overwrite any "global
  data" keys in the root of the super-structure. Its value will be an array,
  where each element is the resulting data structure of each parsed "data"
  file.
 
  Parsed "global data" will be written to the root of the super-structure and
  into the root of each "data" array object. If a key within one of these
  objects conflicts with one of the "global data" keys, then that
  "global data" key will not be written to the object.

TEMPLATES
---------

  All "root" template files passed to dati that have a file extension matching
  one of the supported templating languages will be parsed and executed
  against the super-structure generated by dati.

  All "parital" templates will be parsed into any "root" templates that have a
  file extension that match the same templating language.

SUPPORTED FORMATS / LANGUAGES
-----------------------------

  Below is a list of the supported data-serialisation languages, used for
  "data" and "global data" files.

  - JSON (.json), see https://json.org/
  - YAML (.yaml), see https://yamllint.com/
  - TOML (.toml), see https://toml.io/

  These are the currently supported templating languages, used for files
  passed in the "root" and "partial" arguments.

  - mustache (.mu, .mustache), see https://mustache.github.io/
  - golang text/template (.tmpl, .gotmpl), see https://golang.org/pkg/text/template/
  - golang html/template (.hmpl, .gohmpl), see https://golang.org/pkg/html/template/
    - note that this and text/template are almost interchangable, with the
    exception that html/template will produce "HTML output safe against code
    injection".
<!--  - statix (.stx .statix), see https://gist.github.com/plugnburn/c2f7cc3807e8934b179e -->

EXAMPLES
--------

	dati -cfg ./dati.cfg -r templates/textfile.mst

	dati -r homepage.hmpl -p head.hmpl -p body.hmpl -gd meta.json -d posts/*

  see the examples/ directory in the dati repository for a cool example.

LIBRARIES
---------

  As stated above, all of these libraries do the hard work, dati just combines
  it all together - so thanks to the authors. Also here for reference.

  - The Go standard library is used for parsing JSON, .tmpl/.gotmpl, .hmpl/.gohmpl
  - github.com/pelletier/go-toml
  - gopkg.in/yaml.v3
  - github.com/cbroglie/mustache

AUTHORS
-------

  - gearsix <gearsix@tuta.io>

