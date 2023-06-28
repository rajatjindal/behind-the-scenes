# pets-of-fermyon

A lot of us love pets, and post the cute pictures in #pets channel. But why not share the cuteness with the world?

With that in mind, I am proposing `pets-of-fermyon` slack app. 

Idea is when someone react using a specific emoji to a picture posted in #pets channel, we will automatically share the cuteness with world by posting the picture on BlueSky (and/or Twitter).

Following is the manifest of the app:

```
display_information:
  name: pets-of-fermyon
features:
  bot_user:
    display_name: pets-of-fermyon
    always_online: false
oauth_config:
  scopes:
    bot:
      - files:read
      - channels:history
      - reactions:read
settings:
  event_subscriptions:
    request_url: https://pets.fermyon.app
    bot_events:
      - reaction_added
  org_deploy_enabled: false
  socket_mode_enabled: false
  token_rotation_enabled: false

```


FAQ

### can the app accidently share pictures from other channels?

No, it will only receive events from channels where the app bot account has been added. We have also built a safegaurd in app to ignore events from all channels except the one configured in `spin.toml`
