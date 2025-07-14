<h1 align="center">
  <svg width="200" height="200" viewBox="0 0 150 150" xmlns="http://www.w3.org/2000/svg">
    <circle cx="75" cy="75" r="75" fill="none" stroke="black" />
    <image href="static/reconeer-icon.png" x="35" y="35" width="80" height="80" />
    <g id="hero-rotate">
      <animateTransform attributeName="transform" type="rotate" from="0 75 75" to="360 75 75" dur="30s" repeatCount="indefinite"/>
      <circle cx="150.0" cy="75.0" r="15.0" fill="rgba(255,255,255,0.8)" />
      <circle cx="98.17627457812105" cy="146.3292387221365" r="15.0" fill="rgba(255,255,255,0.8)" />
      <circle cx="14.323725421878947" cy="119.0838939219355" r="15.0" fill="rgba(255,255,255,0.8)" />
      <circle cx="14.323725421878933" cy="30.916106078064523" r="15.0" fill="rgba(255,255,255,0.8)" />
      <circle cx="98.17627457812104" cy="3.670761277863477" r="15.0" fill="rgba(255,255,255,0.8)" />
      <rect x="140.77653323377163" y="89.93862482821741" width="20" height="4" rx="2" fill="#d8bffd" transform="rotate(110 150.77653323377163 91.93862482821741)" />
      <rect x="112.16025403784438" y="134.68395609140177" width="20" height="4" rx="2" fill="#fffacd" transform="rotate(150 122.16025403784438 136.68395609140177)" />
      <rect x="48.06137517178259" y="148.77653323377163" width="20" height="4" rx="2" fill="#b5e7a0" transform="rotate(200 58.06137517178259 150.77653323377163)" />
      <rect x="3.3160439085982283" y="120.16025403784438" width="20" height="4" rx="2" fill="#f1948a" transform="rotate(240 13.316043908598228 122.16025403784438)" />
      <rect x="-10.776533233771634" y="56.06137517178259" width="20" height="4" rx="2" fill="#85c1e9" transform="rotate(290 -0.7765332337716329 58.06137517178259)" />
      <rect x="17.83974596215558" y="11.316043908598" width="20" height="4" rx="2" fill="#a2d9ce" transform="rotate(330 27.83974596215558 13.316043908598)" />
      <rect x="81.93862482821743" y="-2.776533233771632" width="20" height="4" rx="2" fill="#f7dc6f" transform="rotate(380 91.93862482821743 -0.776533233771632)" />
      <rect x="126.68395609140174" y="25.839745962155575" width="20" height="4" rx="2" fill="#d2b4de" transform="rotate(420 136.68395609140174 27.839745962155575)" />
    </g>
  </svg>
  <br>
</h1>

<h4 align="center">Fast subdomain enumeration client for reconeer.com API.</h4>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/reconeer/reconeer"><img src="https://goreportcard.com/badge/github.com/reconeer/reconeer"></a>
  <a href="https://github.com/reconeer/reconeer/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
  <a href="https://github.com/reconeer/reconeer/releases"><img src="https://img.shields.io/github/release/reconeer/reconeer"></a>
  <a href="https://twitter.com/reconeer"><img src="https://img.shields.io/twitter/follow/reconeer.svg?logo=twitter"></a>
  <a href="https://discord.gg/reconeer"><img src="https://img.shields.io/discord/123456789.svg?logo=discord"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Install</a> •
  <a href="#running-reconeer">Usage</a> •
  <a href="#api-setup">API Setup</a> •
  <a href="#reconeer-go-library">Library</a> •
  <a href="https://discord.gg/reconeer">Join Discord</a>
</p>

---

`Reconeer` is a subdomain enumeration tool designed to discover valid subdomains for websites using the reconeer.com API. It features a modular architecture optimized for speed and efficiency. `Reconeer` is built for one purpose—passive subdomain enumeration—and it excels at it.

The tool complies with the usage policies of the reconeer.com API. Its passive approach ensures rapid and discreet enumeration, making it ideal for penetration testers and bug bounty hunters.

# Features

<h1 align="left">
  <img src="static/reconeer-logo.png" alt="Reconeer" width="700px">
  <br>
</h1>

- Fast and efficient subdomain enumeration via reconeer.com API
- **Curated** integration with reconeer.com for reliable results
- Multiple input and output options (file, stdout)
- Optimized for speed and lightweight resource usage
- **STDIN/OUT** support for seamless workflow integration

# Usage

```sh
reconeer -h
