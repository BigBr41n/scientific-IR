# Scientific Information Retrieval System

An implementation of popular Information Retrieval (IR) models, including:

- Classic Inverted Index Intersection
- Vector Space Model (VSM)
- Latent Semantic Indexing (LSI)
- BM25

## Usage

### Data

Dummy documents are available in the following directory:

```sh
./data/documents
```

### Running the Program

1. To run the program directly, use:

```sh
go run /cmd/main
```

2. Alternatively, to build and run the program based on the OS:

- Windows :
  ```sh
  GOOS=windows GOARCH=amd64 go build -o main-windows.exe ./cmd/main
  ```
- linux :
  ```sh
  GOOS=linux GOARCH=amd64 go build -o main-linux ./cmd/main
  ```

### Querying

1. Once the program is running, enter a query in the TUI (Text User Interface) input prompt.

2. View the results, which are generated using all the implemented Information Retrieval models.
