# Brainfuck Web Server
A Web Server written in Brainfuck

## Run
```bash
go run . <brainfuck file>
```
The server will start on port 8080 by default. You can change the port in line 181 of `main.go`

## Creating a web server
### Requests
The HTTP request will be sent to the Brainfuck input in the following format:

```
<METHOD> <URL> <QUERY>
<BODY>
```

Example:

```
GET /collection ?id=1
```

The request information will send an \x03 when its finished
You can read all the input with the following code
```brainfuck
---                 Remove 3 to validate if input is EOT
[                   If EOT (x03) then cell should be 0
  +++
  Do something with the input here
  ,---
]
```

### Response
All output will be send back to the request as `text/plain` or  `text/html` depending on its contents


### Example
The following example creates a basic web application

```brainfuck
,---                 Remove 3 to validate if input is EOT
[                    If EOT (x03) then cell should be 0
  Print HTML
  +++
  >>
  h1 = Brainfuck Web Server
  ++++++++++[>++++++>++++++++++>+++++>++++++>+++++++>+++++++++++>++++++++++>++++++++++>+++++++++++>++++++++++>++++++++++++>++++++++++>+++++++++++>+++>+++++++++>++++++++++>++++++++++>+++>++++++++>++++++++++>+++++++++++>++++++++++++>++++++++++>+++++++++++>++++++>+++++>++++++++++>+++++>++++++<<<<<<<<<<<<<<<<<<<<<<<<<<<<<-]>>++++>->++>---->++++>--->+++++>>++>--->->--->++>--->+>-->++>+++>+>++++>-->+>++++>>--->++++>->++<<<<<<<<<<<<<<<<<<<<<<<<<<<<[.>]<[[-]<]
  p = Open p tag
  ++++++++++[>++++++>+++++++++++>++++++<<<-]>>++>++<<[.>]<[[-]<]
]

Print request information
<<
---                 Remove 3 to validate if input is EOT
[                   If EOT (x03) then cell should be 0
  +++
  .
  ,---
]

Close p tag
++++++++++[>++++++>+++++>+++++++++++>++++++<<<<-]>>--->++>++<<<[.>]<[[-]<]
```