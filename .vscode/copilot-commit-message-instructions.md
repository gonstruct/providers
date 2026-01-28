## Generating commit messages instruction

Always start with a proper scope:
- Use "feat:" for new features.
- Use "fix:" for bug fixes.
- Use "docs:" for documentation changes.
- Use "style:" for code style changes (formatting, missing semicolons, etc).
- Use "refactor:" for code changes that neither fix a bug nor add a feature.
- Use "perf:" for performance improvements.
- Use "test:" for adding or updating tests.
- Use "chore:" for changes to the build process or auxiliary tools and libraries.
- Use "ci:" for changes to CI configuration files and scripts.
- Use "build:" for changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm).
- Use "revert:" for reverting a previous commit.
- Use "wip:" for work in progress commits that are not ready to be reviewed.
- Use "hotfix:" for critical bug fixes that need to be addressed immediately.

When writing the commit message:
- Only write a subject line, no body or footer, just 1 line.
- Use the imperative mood in the subject line (e.g., "Fix bug" instead of "Fixed bug" or "Fixes bug").
- Limit the subject line to a maximum of 56 characters.
- DO NOT Capitalize the first letter of the subject line.
- Do not end the subject line with a period.