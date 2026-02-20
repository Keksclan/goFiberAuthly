# Contributing

Thanks for your interest! This is an **example / test repository** for
[goAuthly](https://github.com/Keksclan/goAuthly) + Fiber v3, so contributions
should keep it **minimal and focused**.

## Guidelines

- **No big feature additions.** The goal is to demonstrate goAuthly integration,
  not to build a full application.
- **Correctness first.** PRs that fix bugs, improve compatibility with goAuthly,
  or clarify documentation are very welcome.
- Keep the code formatted and passing basic checks before opening a PR:

  ```bash
  gofmt -w .
  go vet ./...
  go build ./...
  go test ./...
  ```

## Commit Style

There is no strict convention enforced. A clear, concise commit message is
enough — for example:

```
fix: correct audience check when list is empty
docs: add troubleshooting section for JWKS errors
```

## Questions?

Open an issue if anything is unclear. Thank you!
