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
  <a href="https://github.com/ngtuonghy/nipo-tunnel">
    <img src="assets/logo.png?v=2" alt="Logo" width="240">
  </a>
<h3 align="center">Nipo Tunnel</h3>

  <p align="center">
    Fast, secure, and 100% free HTTP tunnels from local to the internet. Powered by Cloudflare - no account or config required.
    <br />
    <a href="docs/README-vi.md"><strong>Tiếng Việt</strong></a>
    <br />
    <br />
    <a href="https://github.com/ngtuonghy/nipo-tunnel">View Demo</a>
    ·
    <a href="https://github.com/ngtuonghy/nipo-tunnel/issues">Report Bug</a>
    ·
    <a href="https://github.com/ngtuonghy/nipo-tunnel/issues">Request Feature</a>
  </p>
</div>

<div align="center">
  <video src="https://github.com/user-attachments/assets/6382c1f1-bde6-4d60-b1e4-e370e45cc122" controls="controls" muted="muted" autoplay="autoplay" width="100%"></video>
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

Nipo is a lightweight CLI tool that exposes local servers to the public internet instantly.

### Key Features
* **Interactive TUI**: Displays real-time bandwidth, connection status, and logs.
* **Multi-Tunneling**: Expose multiple ports simultaneously via a local `nipo.yml` config.
* **Automatic Setup**: Instantly downloads the native `cloudflared` daemon on the first run.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

[![Go][Go-badge]][Go-url]
[![Node.js][Node-badge]][Node-url]
[![Cloudflare Workers][Cloudflare-badge]][Cloudflare-url]
[![Charm][Charm-badge]][Charm-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Follow these instructions to install and use Nipo.

### Prerequisites

* **Node.js** (v14 or higher) - used to install the global binary runner.

### Installation

Install Nipo globally using npm:

```bash
npm install -g nipo-tunnel
```

### Uninstallation

To completely remove Nipo Tunnel and its configuration from your system:

```bash
# 1. Uninstall the npm package
npm uninstall -g nipo-tunnel

# 2. Remove configuration and core binary
# Windows (PowerShell)
Remove-Item -Recurse -Force ~/.nipo
# Mac/Linux
rm -rf ~/.nipo
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
Yes! The npm package is simply a convenient wrapper. If you do not have Node.js installed, you can go to the [Releases](https://github.com/ngtuonghy/nipo-tunnel/releases) section, download the native compiled Go binary for your OS, rename it to `nipo` (or `nipo.exe` on Windows), and add it to your system PATH.

### Executable Permission Issues (Linux / macOS)
Nipo automatically sets execution permissions (`chmod +x`) on all downloaded helper binaries. However, if you encounter a `permission denied` error running the cloudflared daemon, you can run:
```bash
chmod +x ~/.nipo/bin/cloudflared
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [ ] Custom Domain mapping (CNAME verification).

See the [open issues](https://github.com/ngtuonghy/nipo-tunnel/issues) for a full list of proposed features (and known issues).

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

Project Link: [https://github.com/ngtuonghy/nipo-tunnel](https://github.com/ngtuonghy/nipo-tunnel)

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
[contributors-shield]: https://img.shields.io/github/contributors/ngtuonghy/nipo-tunnel.svg?style=for-the-badge
[contributors-url]: https://github.com/ngtuonghy/nipo-tunnel/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/ngtuonghy/nipo-tunnel.svg?style=for-the-badge
[forks-url]: https://github.com/ngtuonghy/nipo-tunnel/network/members
[stars-shield]: https://img.shields.io/github/stars/ngtuonghy/nipo-tunnel.svg?style=for-the-badge
[stars-url]: https://github.com/ngtuonghy/nipo-tunnel/stargazers
[issues-shield]: https://img.shields.io/github/issues/ngtuonghy/nipo-tunnel.svg?style=for-the-badge
[issues-url]: https://github.com/ngtuonghy/nipo-tunnel/issues
[license-shield]: https://img.shields.io/github/license/ngtuonghy/nipo-tunnel.svg?style=for-the-badge
[license-url]: https://github.com/ngtuonghy/nipo-tunnel/blob/main/LICENSE
[Go-badge]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Node-badge]: https://img.shields.io/badge/Node.js-339933?style=for-the-badge&logo=node-dot-js&logoColor=white
[Node-url]: https://nodejs.org/
[Cloudflare-badge]: https://img.shields.io/badge/Cloudflare_Workers-F38020?style=for-the-badge&logo=cloudflare&logoColor=white
[Cloudflare-url]: https://workers.cloudflare.com/
[Charm-badge]: https://img.shields.io/badge/Charm-8A2BE2?style=for-the-badge&logo=charm&logoColor=white
[Charm-url]: https://charm.land/
