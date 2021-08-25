# github2orgmode

github2orgmode is a small tool to convert issues from defined repositories into
an org file.

## License

MIT see [LICENSE](https://github.com/dragonchaser/github2orgmode/blob/master/LICENSE) file in this repository.

## Build

Just run `make` in the project folder, binary can be found in `bin/`

## Run

```
$> GH_TOKEN=<your-personal-github-access-token> GH_REPOS="org1/repo1 org1/repo2"
bin/github2orgmode > output.org
```

The personal github access token can be created in your profile page [https://github.com/settings/tokens](https://github.com/settings/tokens).
