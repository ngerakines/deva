# deva

A command line interface for Elgato Keylight devices.

# Usage

First, create a configuration directory.

    $ mkdir -p ~/.config/deva/

Next, use the `keylight:discovery` subcommand to find your keylights.

    $ deva keylight:discovery > ~/.config/deva/10-keylight.json

Use the `mode:meeting` subcommand to turn on your lights.

    $ deva mode:meeting

Use the `mode:normal` subcommand to turn off your lights.

    $ deva mode:normal

Use the `config:validate` subcommand to print merged configuration

## Streamdeck UI

Want to add it to streamdeck?

```json
{
  "streamdeck_ui_version": 1,
  "state": {
    "your_device_id": {
      "buttons": {
        "0": {
          "keys": "",
          "write": "",
          "text": "On Video",
          "command": "/home/ngerakines/bin/deva mode:meeting --config /home/ngerakines/.config/deva/"
        },
        "1": {
          "keys": "",
          "write": "",
          "text": "Off Video",
          "command": "/home/ngerakines/bin/deva mode:normal --config /home/ngerakines/.config/deva/"
        }
      }
    }
  }
}
```

# Configuration

This tool loads configuration from JSON files. The lookup order is:

1. The value of the `DEVA_CONFIG` environment variable
2. The value of the `--config` program argument.
3. The file ./deva.json
4. The file $HOME/.deva.json
5. The directory $HOME/.deva/
6. The directory $HOME/.config/deva/
7. The file /etc/deva.json
8. The directory /etc/deva/

# License

MIT License

Copyright (c) 2020 Nick Gerakines
