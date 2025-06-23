# Neoman

[![Static Badge](https://img.shields.io/badge/Read_Docs-%24_nman_inodaf%2Fnman-black)
](https://nman.local/inodaf/nman)

It is unintuitive nowadays searching and reading documentations for the dependencies of our software stack and general software manuals. Heavily inspired by `man`, this software wants to bring back the practicality of _man_ and add new features for enhanced productivity and user experience. Neoman aims to provide a streamlined experience for \*
\*maintainers** and **readers\*\* by offering key features.

**Zero-Deployment**

Maintainers can focus on creating high-quality documentation and manuals and spend less time on having to reason about deployment of the documentation site and have all the burden to manage it. Also decreasing (or zeroing) costs with infrastructure which is especially important of open-source softwares.

**Git on the Core**

By leveraging Git, maintainers have not to worry about notifying updates in other place than pushing to a repository on any Git Remote Provider (GitHub, GitLab). Readers will get all updates automatically after a `git push`.

**Convention**

Markdown files are the only thing Neoman cares about. To "publish" your documentation on Neoman, add a _/docs_ directory, an _index.md_ and a bunch more of _.md_ files, git-push to your remote and that is it.

**One Place**

Developers should find every documentations they want in one place and quickly add more. Whether you are part of an
organization, documentations and manuals of all softwares in your company should be available instantly. Use the full-text search to get results from all or specific documentation.

```sh
nman inodaf/nman
```

**Local & Secure**

Readers of documentations from their organization don't need to reason about data leaks. Neoman add all documentations locally (even not shared across users in the same device) and spawns a tiny background process. Neoman will never offer
any cloud-based solutions.

**Shareable**

Despite being a local-only software, Neoman pages are shareable with colleagues or other users that have it installed.

Share Neoman's local URLs [https://nman.local/:yourorg/:project](https://nman.local/inodaf/nman) so readers will be at the same page. Maintainers can share the same link or a badge in their _README.md_ file.

[![Static Badge](https://img.shields.io/badge/Read_Docs-%24_nman_inodaf%2Fnman-black)
](https://nman.local/inodaf/nman)

## Installation

```sh
curl https://raw.githubusercontent.com/inodaf/neoman/refs/heads/main/install.sh | bash
```

You will be prompted with your password. Although you can reject, your permission allow the installer to configure the
[https://nman.local](https://nman.local) domain in your `/etc/hosts` file.

## Documentation

Use Neoman itself (of course) or find everything under the [/docs](/docs) directory.

```sh
nman inodaf/nman
```
