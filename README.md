# behind-the-scenes

- A lot of us love pets, and post the cute pictures in #pets channel. But why not share the cuteness with the world? 
- Working from mountains today? Why not share that awesome view of your office with others?
- Setting up booth for our next conference? or at a offsite event with other co-workers? Let's share the glimpse of hard-work that goes into making this happen?

With that in mind, I am proposing `behind-the-scenes` slack app. 

We can chose which way we want to go:

- We'll start a new channel called #behind-the-scenes, where slack members will share pictures of their pets, their awesome mechanical keyboards, their pictures of pit-stop on way to offsite event or setting up a booth for the conference. Just about anything that gives a glimpse of "behind the scenes @your-company" and post these pictures to BlueSky/Twitter (can be triggered by adding a specific emoji reaction). 

Following is the manifest of the app:

```
display_information:
  name: behind-the-scenes
features:
  bot_user:
    display_name: behind-the-scenes
    always_online: false
oauth_config:
  scopes:
    bot:
      - files:read
      - channels:history
      - reactions:read
      - app_mentions:read
settings:
  event_subscriptions:
    request_url: https://bts-test.usingspin.com/api/slack
    bot_events:
      - reaction_added
  org_deploy_enabled: false
  socket_mode_enabled: false
  token_rotation_enabled: false

```


FAQ

### can the app accidently share pictures from other channels?

No, it will only receive events from channels where the app bot account has been added. We have also built a safegaurd in app to ignore events from all channels except the one configured in `spin.toml`
