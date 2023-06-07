# ha-discord-mute-sensor
This is a simple go project to light up a lamp on my desk when i'm muted in discord. i mute myself and then forget while playing games with my friends and instead of fixing myself, i'm going to solve the problem with technology :)

### usage
- copy the `compose-example.yml` and name it something more to your liking (i recommend docker-compose.yml or compose.yml)
- add your Discord token (can be gotten from the discord dev portal. no special intents or anything are needed, just the "view channels" permission)
- add your base home assistant url (http://home_assistant.stinky:8123 or whatever),
- add your discord user ID. google how to find this because i'm too lazy to explain it in a github readme
- add your home assistant auth token (https://developers.home-assistant.io/docs/auth_api/#long-lived-access-token)