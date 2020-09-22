# Dual Universe LUA / JSON / YAML converter

This little tool can convert to / from the JSON format that DU uses when you choose "Export/Import LUA Script to Clipboard" in-game.

It's super-duper experimental atm ... !

## Converting to / from JSON

It will convert to / from a filestructure which looks like;
```
slots
 \- -1.unit
  |- 0.stop().lua
  |- 1.start().lua
  |- 2.tick(timerId).lua
 |- -2.system
 |- -3.library
 |- 0.buttonuno
 |- 1.screen
 |- 2.mycontainer
 |- 3.slot4
```

For handlers with arguments (such as the `tick(timerId)` above), the LUA files will contain a little header like;
```
-- !DU: tick("Live")
```

This header will be parsed for the arguments use in the JSON, and will be filled with the arguments when parsing from JSON.
