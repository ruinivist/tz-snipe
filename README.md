# tz-snipe

Cli tool to infer someone's country with a probability based on timezone info sources.

> could've been a one off script but part of learning go

## Usage

Run it directly with go using one of the flags below.

**By tz:**
list countries in a specific timezone offset.

```bash
tz-snipe --tz +1000
```

**By github username:**
get a user's location and time from their commit patches. for higher rate limits, export `GH_TOKEN`.

```bash
tz-snipe --github torvalds
```

**By time:**
see where in the world it is currently a specific time.

```bash
tz-snipe --time 15:00
```

## Build

build the binary to `./build/tz-snipe`:

```bash
make build
```
