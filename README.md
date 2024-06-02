# IOProdz

requirements: [Golang](https://go.dev/dl/), [Make](https://www.gnu.org/software/make/)

### 🦺 install dev tools

```
sh .tooling/devtools.sh
```

### Setup config

```
cp .env.example .env
# Open .env in editor and setup the values
```

### 👷🏽 run for dev

```
make dev
// or run tests
make test
```

### 🏗️ build and 🚜 run for prod

```
make build && ./ioprodz

```
