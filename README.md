# clog

A changelog fragment manager. Developers create small YAML files (one per branch) describing their changes, and at release time `clog` merges them into `CHANGELOG.md` — no more merge conflicts.

## Installation

### One-line install

```bash
curl -sSfL https://raw.githubusercontent.com/Made-Purple/clog/master/install.sh | sh
```

This detects your OS and architecture, downloads the latest release, installs the binary to `/usr/local/bin`, and sets up shell completions for bash, zsh, or fish.

To install to a different directory:

```bash
INSTALL_DIR=~/.local/bin curl -sSfL https://raw.githubusercontent.com/Made-Purple/clog/master/install.sh | sh
```

### From source

```bash
go install github.com/made-purple/clog/cmd/clog@latest
```

## Quick Start

Initialize a project:

```bash
clog init
```

Create a changelog fragment on your branch:

```bash
clog new --edit
```

Preview the next release:

```bash
clog preview
```

Cut a release:

```bash
clog release
```

---

## Changelog Workflow

This project uses changelog fragments to avoid merge conflicts in `CHANGELOG.md`. Instead of editing the changelog directly, each branch adds a small YAML file describing its changes. At release time, these fragments are merged into `CHANGELOG.md` automatically.

### Project Setup

Run the following once in your repository to create the `changelog.d/` directory, a sample template, and an initial `CHANGELOG.md`:

```bash
clog init
```

This creates:
- `changelog.d/` — where fragment files live
- `changelog.d/sample.yaml` — a reference template (ignored during releases)
- `CHANGELOG.md` — the project changelog (if it doesn't already exist)

If you don't have `clog` installed, you can create these manually. See the manual workflows below.

---

## For Developers

Every branch that introduces a user-facing change should include a changelog fragment. This is a small YAML file in `changelog.d/` that describes what changed.

### Using clog

From your feature branch, run:

```bash
clog new
```

This creates a fragment file named after your branch (e.g. `feature-login-page.yaml` for a branch called `feature/login-page`). Open the file and fill in your entries under the relevant categories.

To create the file and open it in your editor immediately:

```bash
clog new --edit
```

### Manual workflow

If you don't have `clog` installed, copy the sample template:

```bash
cp changelog.d/sample.yaml changelog.d/your-branch-name.yaml
```

Use a descriptive filename based on your branch — replace `/` with `-` and lowercase everything (e.g. `feature/login-page` becomes `feature-login-page.yaml`).

### Filling in your fragment

The fragment file has eight categories. Add your entries as bullet items under the relevant categories and leave the rest as empty strings:

```yaml
deployment:
  - ""
added:
  - "User login page with OAuth support"
  - "Password reset flow"
changed:
  - ""
deprecated:
  - ""
removed:
  - ""
fixed:
  - "Session timeout no longer logs out during active use"
security:
  - ""
yanked:
  - ""
```

**Categories:**

| Category     | Use for                                                  |
|-------------|----------------------------------------------------------|
| deployment  | Notes relevant to deploying this change                   |
| added       | New features                                              |
| changed     | Changes to existing functionality                         |
| deprecated  | Features that will be removed in a future release         |
| removed     | Features removed in this release                          |
| fixed       | Bug fixes                                                 |
| security    | Security-related changes (encourages users to upgrade)    |
| yanked      | Hotfixes — used for emergency releases                    |

Leave categories empty (`- ""`) if they don't apply. Only non-empty entries are included in the release.

### Validating your fragment

Before pushing, you can check that your fragment is valid:

```bash
clog validate
```

To preview what the next release entry will look like with all current fragments:

```bash
clog preview
```

### Commit your fragment

Commit the fragment file along with your code changes. It should be part of your pull request:

```bash
git add changelog.d/your-branch-name.yaml
git commit -m "Add changelog fragment"
```

---

## For Maintainers

At release time, all fragment files in `changelog.d/` are merged into a single release entry in `CHANGELOG.md`.

### Using clog

```bash
clog release
```

This will:

1. Read and validate all fragment files in `changelog.d/`
2. Prompt you for a version number
3. Prompt for optional metadata (e.g. `(98%)(Dev)(Prod)`)
4. Show a preview of the release entry
5. Insert the entry into `CHANGELOG.md`
6. Optionally commit the changes and delete the fragment files
7. Optionally create an annotated git tag

### Manual workflow

If you don't have `clog` installed:

1. **Review fragments** — read through all `.yaml` files in `changelog.d/` (ignoring `sample.yaml`)

2. **Build the release entry** — combine all non-empty entries grouped by category into a markdown block:

   ```markdown
   ## [1.2.0] - 2026-03-27

   ### Added
   - User login page with OAuth support
   - Password reset flow

   ### Fixed
   - Session timeout no longer logs out during active use
   ```

   Categories should appear in this order: Deployment, Added, Changed, Deprecated, Removed, Fixed, Security, YANKED. Omit any categories with no entries.

3. **Insert into CHANGELOG.md** — add the new entry below the header and above any existing version entries. Preserve the `# Notes` footer section at the bottom of the file.

4. **Delete fragments** — remove all `.yaml` files from `changelog.d/` except `sample.yaml`

5. **Commit and tag**:
   ```bash
   git add CHANGELOG.md
   git rm changelog.d/*.yaml
   git checkout changelog.d/sample.yaml
   git commit -m "Release v1.2.0"
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin main --tags
   ```
