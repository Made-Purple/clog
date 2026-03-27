# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [0.0.2] - 2026-03-27
### Added
- `init` now creates a sample.yaml in changelog.d/ as a copy-paste template for manual use
- `release` warns about uncommitted changes (staged, unstaged, untracked) before proceeding
### Fixed
- `release` tagging now uses annotated tags, fixing failure when tag.gpgSign is enabled
- `release` auto-commit no longer removes sample.yaml from changelog.d/
- Git errors now include stderr output instead of only showing exit codes

## [0.0.1] - 2026-03-27
### Added
- Initial build
- Installation script added
- Github actions added for release

# Notes
[Deployment] Notes for deployment
[Added] for new features.
[Changed] for changes in existing functionality.
[Deprecated] for once-stable features removed in upcoming releases.
[Removed] for deprecated features removed in this release.
[Fixed] for any bug fixes.
[Security] to invite users to upgrade in case of vulnerabilities.
[YANKED] Note the emphasis, used for Hotfixes
