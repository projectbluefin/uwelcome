
# How to configure uWelcome

uWelcome has a default built-in look, but it's actually made to be highly customizable.
If you're managing a custom system, it might interest you.

## Where to put your config

You can create a custom config file at `/etc/uwelcome/config.json` (system-wide) or `~/.config/uwelcome/config.json` (user-specific).

## Translations ?

uWelcome supports translations for its interface, as well as for some command descriptions and link names.
They have specific names / codes that are used to get translated strings.

Any other option not listed won't be translated.

## Breaking down the configuration file

Here's a breakdown of the config file options - there's also the example folder if you want to see concrete use cases.

## Greeting

These options allow to customize the welcome message.
You can just add a prefix and a suffix to the existing localized welcome message, but you can also add a fully custom message.

This :

```json
{
  "greeting": {
    "prefix": "> ",
    "suffix": " !"
  }
}
```

would result in :

```txt
> Welcome to Bluefin !
```

And this :

```json
{
  "greeting": {
    "suffix": " :D",
    "message": "Booper dooper"
  }
}
```

would result in :

```txt
Booper dooper :D
```

## Commands

This option allows you to define a list of commands to display in the banner.

Here are the unique codes you can use to get translated strings for command descriptions :

| Strings | English Display (translated) |
| --- | --- |
| `cmd_list` | "List of available commands" |
| `cli_pkg` | "Manage command line packages" |
| `term_bling` | "Enable terminal bling" |
| `banner_toggle` | "Toggle this banner on/off" (there are built-in commands for this) |
| `sys_info` | "View system info" |
| `man_upd` | "Manually update the system" |

```json
{
  "commands": [
    {
      "cmd": "uwelcome toggle",
      "desc": "banner_toggle"
    },
    {
      "cmd": "fastfetch",
      "desc": "sys_info"
    },
    {
      "cmd": "brew help",
      "desc": "cli_pkg"
    },
    {
      "cmd": "cowsay",
      "desc": "Display a cow saying something"
    }
  ]
}
```

## Motd

The `motd` option allows to add messages to the banner !
They will get chosen at random and can both be straight up messages or the result of commands.

```json
{
  "motd": {
    "messages" : [
      "This is a custom tip by yours truly ! :D",
      "This is another custom tip (they won't get translated)"
    ],
    "commands" : [
      "umotd"
    ]
  }
}
```

## Links

This option allows to add custom links to the banner.

There are unique names you can use to get a translated name for the link :

| Strings | English Display (translated) |
| --- | --- |
| `bluesky` | "Bluesky" |
| `discord` | "Discord" |
| `discuss` | "Discuss" |
| `docs` | "Documentation" |
| `donate` | "Donate" |
| `issues` | "Report an issue" |
| `mastodon` | "Mastodon" |
| `matrix` | "Matrix" |
| `source` | "Source Code" |
| `website` | "Website" |

```json
{
  "links": [
    {
      "name": "issues",
      "url": "https://issues.bazzite.gg/"
    },
    {
      "name": "docs",
      "url": "https://docs.bazzite.gg/"
    },
    {
      "name": "discord",
      "url": "https://discord.gg/bazzite"
    },
    {
      "name": "bluesky",
      "url": "https://bluesky.bazzite.gg/"
    },
    {
      "name": "discuss",
      "url": "https://github.com/ublue-os/aurora/discussions"
    },
    {
      "name": "Custom Link (it won't get translated)",
      "url": "https://www.innersloth.com/games/among-us/"
    }
  ]
}
```

## Color

This option allows the banner to use an accent color.

Not setting this option or setting it to something other than the options below will result in the banner using the default colors.

Here are all the different options:

- `auto`

> "auto" sets the color automatically to the system's accent color - the latter being predefined or fully custom.

- `blue`
- `green`
- `orange`
- `pink`
- `purple`
- `red`
- `slate`
- `teal`
- `yellow`

> Those set the color to a predefined color. More could be added in the future.

```json
{
  "color": "orange"
}
```
