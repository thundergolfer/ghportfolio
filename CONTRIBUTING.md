# Contributing Guide / Dev Notes

## Releases

Update the semantic version and comment of the below commands and run them.

```bash
$ git tag -a v0.1.0 -m "First release"
$ git push origin v0.1.0
```

Then simply run `goreleaser` **from the root directory of the project**. 

```
goreleaser
```
