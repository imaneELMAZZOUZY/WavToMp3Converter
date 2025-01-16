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
     go run ./cmd -port ":port"
 ```

3. Move the samples in `samples` directory into `samples/Watched_folder` directory

4. To change the path of the watched directory, update DIRECTORY_TO_WATCH in .env file

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


