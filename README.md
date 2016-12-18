Parfait 
=======

A command-line tool for creating and monitoring cloudformation stacks. 

## Features 

 * Easy command-line syntax for launching and updating stacks with parameters
 * Polling creates and stack updates
 * Polling Cloudwatch Logs

## Installing

Either download the binary from https://dl.equinox.io/lox/parfait/stable, or install with golang:

```
go get github.com/lox/parfait
```

## Usage

### Watch a Stack

This polls the events from a stack until a terminal event occurs.

```bash
parfait watch-stack my-stack
```

### Creating a Stack

```bash
parfait create-stack --file ./templates/vpc.yml my-stack Param1=blah Param2=blah
```

### Updating a Stack

```bash
parfait update-stack my-stack Param1=blah Param2=blah

Changeset:
- Llamas (replace)
- Blah (update)

Continue (Y|n)?
```

### Follow Cloudwatch Logs

This polls the events from a stack until a terminal event occurs.

```bash
parfait follow-logs --group my-log-group
```

