# TMDB CLI

**TMDB CLI** is a command-line tool that allows you to get the latest movies from TMDB (The Movie Database)

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/Diwanshumidha/go.git
   ```

2. Build the project:

```sh
cd tmdb
go install
```

---

## Usage

### General Usage

```sh
tmdb [flags]
tmdb [command]
```

### Example

Get the latest movies of type `popular` (default):

```sh
tmdb
```

Get `upcoming` movies:

```sh
tmdb --type upcoming
```

---

## Available Commands

### **help**

Display help information about any command.

```sh
tmdb help
```

---

### **key**

Manage your TMDB API key stored in the system keyring.

```sh
tmdb key [command]
```

#### Subcommands:

- **set** – Set your TMDB API key.

  ```sh
  tmdb key set <key>
  ```

- **get** – List your stored TMDB API key.

  ```sh
  tmdb key get
  ```

- **delete** – Delete your TMDB API key.

  ```sh
  tmdb key delete
  ```

---

## Flags

| Flag            | Description                                                      | Default   |
| --------------- | ---------------------------------------------------------------- | --------- |
| `-h, --help`    | Show help for the `tmdb` command.                                | -         |
| `--type string` | Type of movies to get (`popular`, `top`, `upcoming`, `playing`). | `popular` |

---

## Example Usage

1. **Set the API key**:

```sh
tmdb key set <key>
```

2. **Get popular movies**:

```sh
tmdb --type popular
```

3. **Get upcoming movies**:

```sh
tmdb --type upcoming
```

4. **Delete the API key**:

```sh
tmdb key delete
```

---

## License

[MIT](./LICENSE)

---

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

---

## Author

- [Diwanshu Midha](https://github.com/Diwanshumidha)
