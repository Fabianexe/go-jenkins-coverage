# go2jenkins
Create coverage files for golang that the Jenkins coverage plugin can read and display.

## The accuracy of `go test -coverprofile`
The `go test -coverprofile` command is a great tool to get coverage information about your project.
However, it measures the coverage on a bock level. This means that if you function conatins empty lines, only comments, 
or lines with only a closing bracket, they will be counted in line metrics.

This project tries to solve this problem by using the `go/ast` package to determine the actual lines of code from the source.

Another result from this is that branches on a line level can be determined. If a line contains an `if` statement,
with multiple conditions, it is still one block for the coverage profile. There are projects that try to solve this problem
for example [gobco](https://github.com/rillig/gobco). However, they for the moment not compatible with the Jenkins coverage plugin.
Thus, we add branch coverage on method and file level. Where such multi condition statements are counted as one branche.

## Others
So far we are aware about two other projects that do something similar:
* [gocov-xml](https://github.com/axw/gocov-xml)
* [gocover-cobertura](https://github.com/boumenot/gocover-cobertura) 

However, both of them focus on the coverage part and take over a big downsides of the `go test -coverprofile` command.
Only packages with any coverage are included in the report. 
This means that if you have a package with no tests at all, it will not be included in the report. 
This is a big problem if you want to have a complete report of your project or any meaningfully coverage metrics.

Further this project is more about the Jenkins integration. So it is more than just a coverage tool.
It adds complexity metrics, more options to determine coverage, and branch coverage.