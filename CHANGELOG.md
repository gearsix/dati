# CHANGELOG

## v1.3.0

- added (*Template).ExecuteToFile
- made sure the errors returned by `LoadData` and `WriteData` always indicates the DataFormat being parsed. 
- Added 'force' param to `WriteFileData`
- Bugfix in `WriteFileData`, now it creates non-existing files

## v1.2.2

- Added standardised package errors:
  `ErrUnsupportedData`, `ErrUnsupportedTemplate`, `ErrUnsupportedTemplateType`, `ErrRootPathIsDir`, `ErrNilTemplate`.
- Minor improvements to comment docs in a few places
- Formatting changes to README.md
- Added some TODO items

## v1.2.1

- updated LICENSE dates
- added `DataWrite` and `DataWriteFile` (with passing tests)
- added TODO file
- Renamed some parameters (more uniform & sensical)
- Renamed a few variables (`err` not `e`)

## v1.1.0

- Added `IsDataFormat` and `IsTemplateLanguage` (with passing tests)

## v1.0.0

- Depreciated functions were removed
- Fix to ReadDataFormat and ReadTemplateLanguage
- Updated cmd/dati.go to use new functions
- Renamed examples/suti.cfg to examples/dati.cfg (fixed dati_test)

## v0.8.0

- rebranded to dati - data and template interface
- Added `LoadDataFile`, depreciated `LoadDataFilepath`
- Added `DataFormat`, `ReadDataFormat` depreciated `SupportedDataLangs`, `IsSupportedDataLang`
- Added `TemplateLanguage`, `ReadTemplateLanguage` depreciated `SupportedTemplateLangs`, `IsSupportedTemplate`

## v0.7.0

- added LoadTemplate, loads templates from io.Reader params.
  - Template language and name must be specified.
  - This function is the template.go counterpart to data.go#LoadData
- added LoadTemplateStrings loads templates from string params (calls LoadTemplate)
- renamed LoadTemplateFile -> LoadTemplateFilepath
- added []SupportedTemplateLangs, IsSupportedTemplateLang(), []SupportedDataLangs, IsSupportedDataLang()
- renamed suti.Template.Source -> suti.Template.Name
  - The .Name value assigned to templates when LoadTemplateFilpath() is called will
  be the **base name minus the file extension**.
- Lots of refactors to tidyup, improve consistency & error handling/catches
- Testing is all upto date :)

## v0.2.2

- added .Source to suti.Template, contains the filepath of the source file for the template
- More tidyup

## v0.2.1

- Tidyup work

## v0.2.0

- massive refactored to data API:
  - removed type Data, user will want to use their own data types
  - LoadData() and LoadDataFile() now takes an interface{} (pointer) to write the result to
  - removed LoadDataFiles(), bulk loading is now user responsibility
  - removed GenerateSuperData() & MergeData(), these were superfluous fucntions
- added file.go, which currently only has SortFileList() to stand-in for the missing data sort
functionality found in LoadDataFiles()
- refactored cmd/suti.go to work with the new API, no functional changes
- ".mustache" is now accepted as a mustache file extension (as intended)
- bugfixes

## v0.1.1

- updated everything to use the new (t *Template).Execute()

## v0.1.0

- type Template now a struct, the Template interface is accessed via (t *Template).Template
- ExecuteTemplate -> (t *Template).Execute()
- (t *Template).Execute() takes an interface{} param to execute against (not Data)
- improved error messaging

## v0.0.0

first release!

- fully doc'd in README & src comments
- examples/
- first iteration of suti API
- cmd/suti.go

