# WavToMp3Converter

## Description
WavToMp3Converter is a Go-based application that monitors a directory for `.wav` and `.json` files, 
converts the `.wav` files to `.mp3` using FFmpeg and the configuration from `.json` file, and stores the conversion records in an SQLite database.

## Prerequisites
- Go 1.23.4 or later

## Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/imaneELMAZZOUZY/WavToMp3Converter.git
    ```
2. Navigate to the project directory:
    ```sh
    cd WavToMp3Converter
    ```

## Usage

1. Run the application:
    ```sh
    go run cmd/main.go
    ```
    The application will be running on port 5000

2. To change the port , you can use the flag -port :
 ```sh 
     go run ./cmd -port ":new_port"
 ```

3. Move the samples in `samples` directory into `samples/Watched_folder` directory

4. To change the path of the watched directory, you can use the flag -d :
 ```sh 
     go run ./cmd -d "new_path"
 ```


## Developed APIs

Below are the available APIs for managing and querying job statuses:

### 1. `/api/jobs/current`
- **Method**: `GET`
- **Description**: Retrieves the jobs that are currently being processed by the job goroutine.

Example `curl` command to view current jobs:
```bash
curl -i http://localhost:5000/api/jobs/current
```

### 2. `/api/jobs/finished`
- **Method**: `GET`
- **Description**: Retrieves the jobs that have been completed and stored in the database. 
  - Optionally, you can pass a `status` parameter to filter finished jobs based on their success or failure (e.g., `status=successful` or `status=failed`).

Example `curl` command to view finished jobs:
```bash
curl -i http://localhost:5000/api/jobs/finished?status=successful
```

### 3. `/api/jobs/waiting`
- **Method**: `GET`
- **Description**: Retrieves the jobs that have been uploaded but are waiting to be processed by the job goroutine. These jobs are in a queued state and have not started processing yet.

Example `curl` command to view waiting jobs:
```bash
curl -i http://localhost:5000/api/jobs/waiting
```


## JSON Configuration

Each `.json` file should have the following structure:

```json
{
    "input_file": "example.wav",
    "output_file": "example.mp3",
    "codec": "libmp3lame",
    "bitrate": "192k",
    "sample_rate": "44100",
    "channels": "2"
}
```


