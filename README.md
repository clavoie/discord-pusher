# discord-pusher
A small Go web app that pushes build notifications to [Discord](https://discordapp.com/) for [BitBucket](https://bitbucket.org/) and [Unity](https://unity3d.com/).

![discord-pusher-ui](https://user-images.githubusercontent.com/8398/28479135-1db39f1e-6e29-11e7-8ccc-08eebc422fcb.png)

![discord-build-image](https://user-images.githubusercontent.com/8398/28479297-cbb1f4a8-6e29-11e7-9bbe-cdc26508720c.png)

## Usage
- deploy the application to wherever you would like it to run. currently only Google App Engine is supported but it can be extended to other hosting environments via the [discord-pusher-deps](https://github.com/clavoie/discord-pusher-deps) package
- open `http://yourhost.com/` in your browser where the application is deployed
- generate a new webhook link in Discord
- copy the link and paste it into the Discord Web Url form field of the application
- select the type of notifications you'd like to recieve and click Add
- a webhook link will be generated for you. copy that link and paste it into either Bitbucket or Unity

## Google Cloud App Engine Deployment
`gcloud app deploy app.yaml`

## Extending
To support other hosting environments, clone [discord-pusher-deps](https://github.com/clavoie/discord-pusher-deps), implement `github.com/clavoie/discord-pusher-deps/types/HookContext`, and add it back to discord-pusher-deps. 

Go to `github.com/clavoie/discord-pusher/init.go` and update `pusherTypes` with your implementation. You can then start up the server with your implementation by either:

- setting the environment variable `DISCORD_PUSHER_TYPE=your_type`
- providing your pusher type to the application as a command line argument

## Why is HookContext in a separate repository
Go App Engine sadness
