# Dubby: Dual Universe Builder
Dubby is a small CLI tool to support the development of scripts for Dual Universe.  

You can convert to/from the JSON format which the game uses to "Import/Export to Clipboard",  
and to/from the YAML format which the game uses for the auto configure modules.

I'm working hard on making it more mature and stable, any feedback or bug reports are *very* welcome!

*auto configure is still WIP ...*

## Install
The easiest way to install Dubby is to grab one of the binaries from the 
[releases page](https://github.com/rubensayshi/dubby/releases).

Alternatively if you have `go` installed you could `go install https://github.com/rubensayshi/dubby/src/cmd/dubby`.

## Usage
Use `dubby -h` for the most up-to-date description of the available commands;

 - `dubby parse-to-src import.json ./src`
 - `dubby export-to-json ./src export.json`

### Minifying
the `dubby export-to-json` command has a `--minify` flag, which expects a `luamin` binary to be present on your machines,
 which is a NPM package (https://www.npmjs.com/package/luamin) and you can easily install this using `npm install -g luamin`.

## Why convert to (separate) lua files?
There's a few reasons to want to convert to lua files, though some of them are subjective...
 - Easier to maintain.
 - Easier to develop on small pieces of code.
 - Easier for others to find and use snippets of code.
 - Possible to write unit tests.

Then some undisputable ones;
 - You can easily minify the lua code when compiling.
 - You can run static analysis tools over the code, such as a linter 
    to ensure there's no syntax errors and that it somewhat conforms to a code style.

## Source structure
The source structure is simple, and the main goal in its design are to make it easy to separate out code, at minimal effort.

There's 2 main folders `lib/` and `slots/`; 

#### `lib/`
The `lib/` folder can contain any number of lua files and they will all be concatenated together 
and placed in a `unit.start()` filter when compiled.  

It's recommended to place most globals here, such as functions, "classes" and global state.  

Files are concatenated in order, which shouldn't matter if they're not dependenant on each other,
but if it does matter then prefixing them with a number.  
I recommended you start with `00`, `10`, `20`, etc. as that makes it easy to add something in between without having to rename the others.

When not minifying on export, the compiled code with contain markers to make sure that if you'd parse it again,
 it will be placed in the same directories. 

#### `slots/`
The `slots/` folder contains the filters for each slots,
each slot is contained in a single file in the format of `%d.%s.lua` where `%d` is the number of the slot and `%s` the name.  
Currently the game limits you to 10 slots, and the 3 default slots: `-1.unit`, `-2.system` and `-3.library`.

Any code which hasn't been marked with a filter, will be placed in the `start()` filter of the slot.

Filters can be added by declaring them in `do end` block with a `-- !DU: filter([args])` comment, for example;
```
someGlobal = 1
function doSomethingInLiveTick() end

do -- !DU: tick([Live])
    doSomethingInLiveTick()
end

do -- !DU: actionStart([button])
    someGlobal += 1
end
```
