# Development

This repository uses the [workchat-plugin-starter-template](https://gitlab.com/w1572/workchat-plugin-starter-template). Therefore, developing this plugin is roughly the same as it is with every plugin using the template. All the necessary steps to develop are in the template's repository.

If you are a Nagios admin/user and think there is something this plugin lacks or something that it does could be done another way, let us know! We are trying to develop this plugin based on users' needs. If there is a certain feature you or your team needs, open up an issue, and explain your needs. We will be happy to help.

This plugin only contains a server portion. Read our documentation about the [Developer Workflow](https://developers.workchat.com/extend/plugins/developer-workflow/) and Developer Setup for more information about developing and extending plugins.

## Developing the watcher

To build the watcher, you can use the following command:

```shell script
env GOOS=linux GOARCH=amd64 go build -o dist/watcherX.Y.Z.linux-amd64 -a -v cmd/watcher/main.go
```

Of course, you can build the watcher for other operating systems and architectures too.

## I saw a bug, I have a feature request or a suggestion

Please file a GitHub issue, it will be very useful!

Pull Requests are welcome! You can contact us on the Workchat Community ~Plugin: Nagios channel.

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Workchat server is running locally, you can enable local mode to streamline deploying your plugin. After configuring it, just run:

```shell script
make deploy
```