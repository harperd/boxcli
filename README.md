# BoxCLI
Command line interface for Aidbox and FHIR.

##Build
BoxCLI provides a wrapper for [JQ](https://stedolan.github.io/jq/) to allow filtering of JSON results.
Before building BoxCLI you will need to install [JQ](https://stedolan.github.io/jq/) first.

After installing [JQ](https://stedolan.github.io/jq/) perform the following to build BoxCLI:

```$xslt
$ cd $GOPATH/src/github.com/harperd/boxcli 
$ go get . && go build -o box
```
##Setup
To use BoxCLI you will need to set an environment variable for each box you want to be able to access. It is recommended to add it to your .bashrc, .profile, .bash_profile, etc.
BoxCLI environment variables should be all upper case and start with BOX_ followed by the name or alias of the box.

For example, to configure a box alias called mybox for (in this case the same as the actual box name) you would set the following environment variable:

```$xslt
export BOX_MYBOX=http://mybox.aidbox.io
```

##Usage
box [box name] [get|put|post|delete] [doc|fhir] [resource] [options] \<JQ query\>

Options:

	-M      monochrome (don't colorize JSON)
	-u      unformatted output
	-c      get the count of the query results only
	-i:n    get the resource at index n. Other value for n is 'last'.
	
## Examples
Retrieve all FHIR Patient resources as a FHIR Bundle
```$xslt
$ box mybox get fhir Patients
```
Get the number of Patient resources
```$xslt
$ box mybox get fhir Patient -c
```
Get the 5th Patient resource
```$xslt
$ box mybox get fhir Patient -i:5
```
Get a specific Patient resource
```$xslt
$ box mybox get fhir Patient/c7f17c3f-414c-4404-bd77-15aaf948ce7c
```
Delete a specific patient resource
```$xslt
$ box mybox delete fhir Patient/c7f17c3f-414c-4404-bd77-15aaf948ce7c
```
Do a FHIR based search then a JQ filter to just get the last modified date/time
```$xslt
$ box mybox get fhir Patient/5c170415-f585-4129-b9e5-ee7cef6af850 ".meta.lastUpdated"
```
Do a JQ type search (Get all patient resources where last name is Baker)
```$xslt
$ box mybox get fhir Patient "select(.name[0].family[0]=='Baker')"
```
Get all names for all Patient resources
````$xslt
$ box mybox get fhir Patient ".name[]"
```
# To Do
1. Complete PUT and POST
