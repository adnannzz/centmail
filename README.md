CentMail is a standalone, self-hosted, newsletter and mailing list manager. It is fast, feature-rich, and packed into a single binary. It uses a PostgreSQL database as its data store.

CentMail is a fork of the open source [listmonk](https://listmonk.app) project, rebranded and maintained independently at [github.com/adnannzz/centmail](https://github.com/adnannzz/centmail).

[![listmonk-dashboard](https://github.com/user-attachments/assets/689b5fbb-dd25-4956-a36f-e3226a65f9c4)](https://github.com/adnannzz/centmail)

*(Screenshot above is from upstream listmonk. CentMail's UI is based on it, plus its own additions — see below.)*

Beyond the rebrand, CentMail adds a few things not in upstream listmonk:
- **Multiple saved subscription forms** (Lists → Forms) — each with its own list selection and redirect URL, with ready-to-use HTML and `<iframe>` embed code, instead of a single ad-hoc list picker.
- **Light / Dark / Auto theme** — a switcher in the profile menu, following the OS theme by default.

## Installation

### Docker

No public Docker image is published for CentMail yet. Build your own image from source and use the included [docker-compose.yml](https://github.com/adnannzz/centmail/blob/master/docker-compose.yml):

```shell
# Clone the repo.
git clone https://github.com/adnannzz/centmail.git
cd centmail

# Build the binary and Docker image.
make dist
docker build -t centmail/centmail:latest .

# Run the services in the background.
docker compose up -d
```
Visit `http://localhost:9000`

__________________

### Binary

- `git clone https://github.com/adnannzz/centmail.git && cd centmail`
- `make dist` to build the self-contained `./centmail` binary (requires Go and Yarn; see [Developers](#developers) below).
- `./centmail --new-config` to generate config.toml. Edit it.
- `./centmail --install` to setup the Postgres DB (or `--upgrade` to upgrade an existing DB. Upgrades are idempotent and running them multiple times have no side effects).
- Run `./centmail` and visit `http://localhost:9000`

__________________


## Developers
CentMail's codebase is derived from listmonk, which is free and open source software licensed under AGPLv3. The [listmonk developer setup guide](https://listmonk.app/docs/developer-setup) mostly applies here too, since the tooling and project layout are unchanged. The backend is written in Go and the frontend is Vue with Buefy for UI.


## License
CentMail, like the listmonk project it's derived from, is licensed under the AGPL v3 license.
