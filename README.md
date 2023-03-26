# Career Page Monitor

I found myself lately checking the career page of a few companies I'm interested in working for. The number of companies I'm interested in is growing and I don't want to check them manually anymore. So I decided to automate the process and write a CLI tool that will notify me when a new job is posted.

If you are following multiple companies, or you want to scale your job search during the current Tech recession, this tool might be useful for you.

## Installation
If you have Go installed, you can install the script by running the following command:
```
go install github.com/Mo-Fatah/cpm@latest
```

## Usage
First you need to run 
```
cpm init
```
This will create a config directory in your home directory.

To add a company, run:
```
cpm add <career-page-url>
```
Ideally, the page url should include the query parameters of your desired job type if it is supported.
for example, this url from hashicorp will only show jobs that are related to software engineering:
```
https://www.hashicorp.com/careers/open-positions?department=Research+%26+Development
```

to check the status of the companies you added, run:
```
cpm status
```
this will inform you if the career pages has changed, and whether the changes regards a new jobs, job removal or something else.



