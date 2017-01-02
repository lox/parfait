Parfait [![Build status](https://badge.buildkite.com/626907f71e5a0fde836085a4aa28d1b22ee36be71d5b88f476.svg)](https://buildkite.com/lox/parfait) [![Latest Release](https://img.shields.io/github/release/lox/parfait.svg)](https://github.com/lox/parfait/releases)
=======

A command-line tool for creating and managing cloudformation stacks. 

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
parfait create-stack --template https://s3-ap-southeast-2.amazonaws.com/cloudformation-templates-ap-southeast-2/WordPress_Single_Instance.template wordpress Param1=blah Param2=blah

2016/12/18 17:35:13 CREATE_IN_PROGRESS -> wordpress [AWS::CloudFormation::Stack] => "User Initiated"
2016/12/18 17:35:19 CREATE_IN_PROGRESS -> WebServerSecurityGroup [AWS::EC2::SecurityGroup] 
2016/12/18 17:35:34 CREATE_IN_PROGRESS -> WebServerSecurityGroup [AWS::EC2::SecurityGroup] => "Resource creation Initiated"
2016/12/18 17:35:34 CREATE_COMPLETE -> WebServerSecurityGroup [AWS::EC2::SecurityGroup] 
2016/12/18 17:35:38 CREATE_IN_PROGRESS -> WebServer [AWS::EC2::Instance] 
...
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

