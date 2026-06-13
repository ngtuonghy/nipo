# Nipo

Production-ready tunneling platform. Create secure HTTP tunnels from localhost to the internet.

## Installation

Install globally via npm:

```bash
npm install -g nipo
```

## Quick Start

Start a secure HTTP tunnel on port 3000:

```bash
nipo http 3000
```

Start a tunnel with a custom subdomain:

```bash
nipo http 3000 -s myapp
```

Start all tunnels defined in `nipo.yml`:

```bash
nipo start
```

## Configuration

To view and configure settings like language:

```bash
nipo config
nipo config --language vi
```

## License

MIT
