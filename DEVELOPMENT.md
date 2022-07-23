# Development Guide

## Changing the `README.md`

Our `README.md` is generated.
Manual modifications in the `README.md` file will be overwritten.

If you want to change the `README.md` file, please make your changes to `assets/README.template`.

## Extending the *yml* format in `/podcasts`

If you aim to modify (add, change or remove) a yaml proprty in the `/podcasts/*.yml`, 
please ensure making this change in all *.yml* files.

Our tooling in `/app` also need adjustment.
Mainly in

* `/app/cmd/types.go`: Adjusting the type structure
* `app/cmd/convertYamlToJson.go` -> `mergePodcastInformation`: Adjusting the Yaml to JSON merge logic