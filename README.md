# logstox

A small logging facade and hub: log anywhere, route everywhere, zero lock-in. Keep your stack; fits beside any logger.

> [!WARNING]
>
> This project is currently in its very early stages and is subject to volatile changes until this notice is removed.
>
> This project is also a learning experience since there are several concepts and technologies I want to explore that
> aren't my expertise... yet :smirk:
>
> TL;DR: this whole thing is a WIP :shipit:
>
> This README will be filled out with further details as the project progresses.

## Why logstox?

What it is:

- A tiny [logging facade](https://en.wikipedia.org/wiki/Facade_pattern): one small, stable interface across projects.
- Drop-in: sits beside whichever logger you use today (zap, slog, stdlib), or want to use tomorrow.
- An optional hub that can [fan-out](https://en.wikipedia.org/wiki/Fan-out_(software)) entries to multiple data sinks (files,
  s3, http).
- A handful of sane defaults and helpers such that you don't need to write a bespoke app-local module from scratch
  everytime.

What it isn't:

- Not a "unify every feature" abstraction.
- Not sticky -- zero lock-in and escape hatches to the underlying logger.
- Not an exhaustive list of backend providers -- you can write your own (perhaps contribute it?).

## Ideology

- Keep the surface small (easy to learn, hard to misuse, impossible to wrangle with).
- Choose backends freely, keep [call sites](https://en.wikipedia.org/wiki/Call_site) unchanged.
- Route logs anywhere via a separate hub/ fan-out layer (files, s3, http) without touching call sites.

## FAQ

- **Q:** Does this bloat my binary?

  **A:** No. Go links only what you import. If you import `logstox` and `provider/zapx`, you don't "pay" for `slogx` or
  data sinks you don't import.

- **Q:** Why not wrap every feature of `zap`/ `slog`?

  **A:** Scope creep and potential lock-in. This facade exposes only the common path; escape hatches let you drop down
  to the underlying logger when needed.

  Your app-local wrapper will leverage this facade and can tailor the interface to suit your needs rather than forcing
  it on everyone.

  That said, if there is enough demand, we can explore expanding the standard layer.

## Roadmap (preview)

This is not a comprehensive outline nor necessarily the order of priorities but it's currently the way I'll tackle this
project (subject to change):

- [x] Stellar README
- [x] Core interfaces and fields
- [ ] Add logging backend provider (eg zap).
- [ ] Add stdlib backend provider
- [ ] Data sink: jsonl encoder, file sink
- [ ] Data sink: s3 sink?
- [ ] Examples

## Contributing

Bug reports and PRs are welcome. Please open an issue first for discussion.

Keep the core surface minimal: avoid adding provider-specific concepts to the facade unless they belong in the standard
interface. This should be clear through a conversation in an issue.

Feel free to open an issue if you spot something iffy or have a hot tip :shrug:

## License

Apache-2.0. See [LICENSE](./LICENSE). If applicable, see [NOTICE](./NOTICE).
