# OpReiADRCLIGo

***Op***inionated ***Rei***nterpretation of ***ADR*** ***CLI*** in ***Go*** === OpReiADRCLIGo. I am just that bad in naming.

This project exists when I learn about ADR (Architecture Decisions Record), and I use it as an excuse for me to develop something in [Go](https://go.dev/) and [Cobra](https://cobra.dev/).


## So what's the opinion?
* This is a CLI that works similar with [existing ADR tools](https://github.com/npryce/adr-tools), but templates, configuration, and portability are my own interpretation.
  - Code is written in Go and using Cobra, so I can have binary that run in both Windows and Linux (_easily_).
  - The tool can still be run by using the command "adr" like existing ADR tools.
  - I want to define my own [Markdown Architectural Decision Record](https://adr.github.io/madr/) format -- not too long, not too short.
  - And I want the ADRs created to always be at "docs/adr", with an index file called README.md


## Installation

### Download binary

Pre-built binaries for Linux, macOS, and Windows are available on the
[Releases page](https://github.com/rearrange/opreiadrcligo/releases).
Download the archive for your platform, extract it, and place the `adr`
binary somewhere on your `PATH`.

**Linux / macOS:**

```sh
# Replace vX.Y.Z and linux-amd64 with your version and platform.
# Available platforms: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64
curl -L https://github.com/rearrange/opreiadrcligo/releases/download/vX.Y.Z/adr-linux-amd64.tar.gz \
  | tar xz
sudo mv adr /usr/local/bin/
```

**Windows:**

Download `adr-windows-amd64.zip` from the Releases page, extract `adr.exe`,
and add its location to your `PATH`.

### Build from source

Requires [Go](https://go.dev/) 1.26 or later.

```sh
git clone https://github.com/rearrange/opreiadrcligo.git
cd opreiadrcligo
go build -o adr .
```

Move the resulting `adr` binary to a directory on your `PATH`.


## Usage

### `adr init`

Initialises the ADR workspace by creating the `docs/adr/` directory and an
index file (`docs/adr/README.md`). Run this **once per project** before
creating any ADRs.

```
$ adr init
✓ Created directory  : docs/adr
✓ Created index file : docs/adr/README.md

You can now create your first ADR:
  adr new "Record your first decision"
```

The command is safe to re-run: it returns an error if the workspace already
exists, and recovers gracefully if the directory is present but the index is
missing.

---

### `adr new <title>`

Creates a new ADR file in `docs/adr/` with a sequential 4-digit number and
the title slugified as the filename. The index (`docs/adr/README.md`) is
updated automatically with the new entry.

```
$ adr new "Use Go for the CLI toolchain"
✓ Created ADR   : docs/adr/0001-use-go-for-the-cli-toolchain.md
✓ Updated index : docs/adr/README.md

Title  : Use Go for the CLI toolchain
Number : 0001
```

The generated file follows the [MADR](https://adr.github.io/madr/) format
with sections for context, considered options, decision outcome, and
consequences. The status defaults to **Draft**.

Each subsequent `adr new` increments the number automatically:

```
$ adr new "Adopt hexagonal architecture"
✓ Created ADR   : docs/adr/0002-adopt-hexagonal-architecture.md
✓ Updated index : docs/adr/README.md

Title  : Adopt hexagonal architecture
Number : 0002
```

---

### `adr list`

Lists all ADRs as a formatted table. The **Status** column reflects the
current value inside each file — so a decision manually updated from
`Draft` to `Accepted` is shown correctly.

```
$ adr list
| #    | Title                          | Date        | Status   |
|------|--------------------------------|-------------|----------|
| 0001 | Use Go for the CLI toolchain   | 2 Apr 2026  | Accepted |
| 0002 | Adopt hexagonal architecture   | 2 Apr 2026  | Draft    |
```

If no ADRs exist yet, a helpful hint is printed instead.

---

The decisions for this tool are recorded as [architecture decision records in this project repository](docs/adr/).
