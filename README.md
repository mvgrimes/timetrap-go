# Timetrap-go

An implmentation of the simple command line time tracking tool
[timetrap](https://github.com/samg/timetrap) in go. 

## Usage

See `tt help` or [timetrap](https://github.com/samg/timetrap) for usage.

## Why

I've been successfully using and enjoying
[timetrap](https://github.com/samg/timetrap) for years. Thank you @samg and the
contributors for a simple and elegant tool.

It was reimplemented in go primarily for speed and simplicity of installation
(no gems, no Ruby version issues, just a single binary). I've also started
adding a few new features (ie, `tt in --sheet "work"` to both switch sheet and
checkin in one command).

## WIP

Some of the features of timetrap have not yet been implemented including
auto_sheets and external formatters. I don't use either of these, so they are
not a high priority.

The database should be compatible with timetrap, so you can switch
back-and-forth as needed.
