## List All Docs For a Project

$ nman list n26/valium # List all docs available in the n26/valium project

This is a list of all the documentation available in the **n26/valium** project. Each entry includes the title of the document and its path within the project. Use this list to quickly find and access specific documents.

---
project: n26/valium
docs: # formatted as "Title: Path"
  - Publishing: 'Authoring Docs/Publishing.md'
  - Private Repositories: 'Teams/Private Repositories.md'
---

View full content of a specific doc (e.g., Publishing):
`$ nman view n26/valium "Authoring Docs/Publishing.md"`

+++ +++

## List All Projects with Docs in an Org/User

$ nman list n26/* # List all available projects with documentation under the n26 organization

This is a list of all projects under the **n26** organization that have available documentation. Use this list to explore the different projects and their associated documentation.

---
author: n26
docs:
  - valium
  - another-project
---

List docs of a specific project (e.g., valium):
`$ nman list n26/valium`

+++ +++

## List All Projects with Docs Across All Orgs/Users

$ nman list # List all available projects with documentation across all organizations

This is a list of all documentation installed. Use this list to discover projects and their documentation from different organizations. Follows the format "org/repo".

---
docs:
  - n26/valium
  - another-org/some-project
---

+++ +++

## Edge Case: Project Not Found

$ nman list unknown/repo

No documentation found for project "unknown/repo". Please check the project name and try again.