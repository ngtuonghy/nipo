<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
<div align="center">

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

</div>

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/ngtuonghy/nipo">
    <img src="assets/logo.png" alt="Logo" width="160" height="160" style="border-radius: 20px;">
  </a>

<h3 align="center">NIPO</h3>

  <p align="center">
    A production-ready secure HTTP tunneling tool from localhost to the internet. Built in Go.
    <br />
    <a href="docs/README-vi.md"><strong>🇻🇳 Đọc bằng Tiếng Việt</strong></a>
    <br />
    <br />
    <a href="https://github.com/ngtuonghy/nipo">View Demo</a>
    ·
    <a href="https://github.com/ngtuonghy/nipo/issues">Report Bug</a>
    ·
    <a href="https://github.com/ngtuonghy/nipo/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Nipo is a lightweight CLI tool that exposes local servers to the public internet instantly, featuring a beautiful terminal UI (TUI).

### Key Features
* **Interactive TUI**: Displays real-time bandwidth, connection status, and logs.
* **Multi-Tunneling**: Expose multiple ports simultaneously via a local `nipo.yml` config.
* **Automatic Setup**: Instantly downloads the native `cloudflared` daemon on the first run.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Go][Go-badge]][Go-url]
* [![Node.js][Node-badge]][Node-url]
* [![Cloudflare Workers][Cloudflare-badge]][Cloudflare-url]
* [![Charm][Charm-badge]][Charm-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Follow these instructions to install and use Nipo.

### Prerequisites

* **Node.js** (v14 or higher) - used to install the global binary runner.

### Installation

Install Nipo globally using npm:

```bash
npm install -g nipo
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

### 1. HTTP Tunneling
To forward public traffic to a local port (e.g., port 3000):
```bash
nipo http 3000
```

### 2. Custom Subdomains
Start a tunnel with a custom subdomain:
```bash
nipo http 3000 --subdomain myapp
# Or shorthand:
nipo http 3000 -s myapp
```

### 3. Multi-tunnel Execution (via nipo.yml)
You can define multiple tunnels in a local configuration file named `nipo.yml`:
```yaml
default_subdomain: "myproject"
tunnels:
  - name: web
    port: 3000                     # Uses default subdomain: myproject-web
  - name: api
    port: 8080
    subdomain: "my-custom-api"  # Overrides and uses: my-custom-api
```
* **Start all tunnels simultaneously**:
  ```bash
  nipo start
  ```
* **Start specific tunnels from the config file**:
  ```bash
  nipo start web
  ```

### 4. Configurations & Language Preferences
Check the active configuration settings:
```bash
nipo config
```
Nipo detects your system locale, but you can change the display language dynamically:
```bash
nipo config --language vi  # Switch to Vietnamese
nipo config --language en  # Switch to English
```
These preferences are saved globally in `~/.nipo/config.yml`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- TROUBLESHOOTING & FAQ -->
## Troubleshooting & FAQ

### Interactive Subdomain Conflict Resolution
If the custom subdomain you request is already in use by another active tunnel, Nipo will pause and prompt you with an interactive selection:
* **Use a random subdomain**: Generates a random fallback subdomain immediately so you can run the tunnel.
* **Exit**: Cancels the execution so you can try another subdomain.

### Can I run Nipo without Node.js?
Yes! The npm package is simply a convenient wrapper. If you do not have Node.js installed, you can go to the [Releases](https://github.com/ngtuonghy/nipo/releases) section, download the native compiled Go binary for your OS, rename it to `nipo` (or `nipo.exe` on Windows), and add it to your system PATH.

### Executable Permission Issues (Linux / macOS)
Nipo automatically sets execution permissions (`chmod +x`) on all downloaded helper binaries. However, if you encounter a `permission denied` error running the cloudflared daemon, you can run:
```bash
chmod +x ~/.nipo/bin/cloudflared
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [ ] TLS / HTTPS local backend support.
- [ ] Custom Domain mapping (CNAME verification).
- [ ] Visual web-based dashboard for request inspection.
- [ ] Add TCP raw port forwarding.

See the [open issues](https://github.com/ngtuonghy/nipo/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

ngtuonghy - [@ngtuonghy](https://github.com/ngtuonghy)

Project Link: [https://github.com/ngtuonghy/nipo](https://github.com/ngtuonghy/nipo)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Cloudflare Tunnel / cloudflared](https://github.com/cloudflare/cloudflared)
* [Charm.sh (Bubble Tea, Lip Gloss)](https://github.com/charmbracelet)
* [Cobra CLI](https://github.com/spf13/cobra)
* [nport](https://github.com/tuanngocptn/nport)
* [Best-README-Template](https://github.com/othneildrew/Best-README-Template)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/ngtuonghy/nipo.svg?style=for-the-badge
[contributors-url]: https://github.com/ngtuonghy/nipo/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/ngtuonghy/nipo.svg?style=for-the-badge
[forks-url]: https://github.com/ngtuonghy/nipo/network/members
[stars-shield]: https://img.shields.io/github/stars/ngtuonghy/nipo.svg?style=for-the-badge
[stars-url]: https://github.com/ngtuonghy/nipo/stargazers
[issues-shield]: https://img.shields.io/github/issues/ngtuonghy/nipo.svg?style=for-the-badge
[issues-url]: https://github.com/ngtuonghy/nipo/issues
[license-shield]: https://img.shields.io/github/license/ngtuonghy/nipo.svg?style=for-the-badge
[license-url]: https://github.com/ngtuonghy/nipo/blob/master/LICENSE
[Go-badge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Node-badge]: https://img.shields.io/badge/Node.js-339933?style=for-the-badge&logo=node-dot-js&logoColor=white
[Node-url]: https://nodejs.org/
[Cloudflare-badge]: https://img.shields.io/badge/Cloudflare_Workers-F38020?style=for-the-badge&logo=cloudflare&logoColor=white
[Cloudflare-url]: https://workers.cloudflare.com/
[Charm-badge]: https://img.shields.io/badge/Charm-8A2BE2?style=for-the-badge&logo=charm&logoColor=white
[Charm-url]: https://charm.land/
