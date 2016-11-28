# BoxCLI
Command line interface for Aidbox and FHIR.

##Build

BoxCLI provides a wrapper for [JQ](https://stedolan.github.io/jq/) to allow filtering of JSON results.
Before building BoxCLI, you will need to install [JQ](https://stedolan.github.io/jq/) first.

After installing [JQ](https://stedolan.github.io/jq/) perform the following to build BoxCLI:

```$xslt
$ cd $GOPATH/src/github.com/harperd/boxcli 
$ go get .
$ go build -o box
```

##Usage
box [get|put|post|delete] [resource] [options] \<JQ query\>

Options:

	-M	monochrome (don't colorize JSON)
	-u	unformatted output
	-c	get the count of the query results only
	-i:n	get the resourced at index n
	
## Examples

Retrieve all FHIR Patient resources as a FHIR Bundle
```$xslt
$ box get patients
```
Get results without syntax highlighting
```$xslt
$ box get patient -M
```
Get the number of Patient resources
```$xslt
$ box get patient -c
```
Get the 5th Patient resources
```$xslt
$ box get patient -i:5
```
Get a specific Patient resource
```$xslt
$ box get patient/c7f17c3f-414c-4404-bd77-15aaf948ce7c
```
Delete a specific patient resource
```$xslt
$ box delete patient/c7f17c3f-414c-4404-bd77-15aaf948ce7c
```
Do a FHIR based search
```$xslt
$ box get patient?subject=c7f17c3f-414c-4404-bd77-15aaf948ce7c
```
Do a JQ type search (Get patient where last name is Baker)
```$xslt
$ box get patient "select(.name[0].family[0]=='Baker')"
```
Get all names for all Patient resources
````$xslt
$ box get patient ".name[]"
```
	
# To Do

1. Complete PUT and POST
2. Add FHIR based search functionality