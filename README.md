# API Fetcher & Data Writer

This is a simple Go program that concurrently fetches data from multiple API endpoints, handles retries, and processes the JSON response into structured data. The program demonstrates basic usage of Go's concurrency features such as goroutines, `sync.WaitGroup`, and channels.

## Features

- Fetches data from multiple APIs concurrently.
- Handles API retries on failure.
- Parses JSON responses into a `DataWriter` struct.
- Displays the parsed data in a user-friendly format.

## How It Works

1. **API Fetching**: 
   - The `fetchURL` function sends an HTTP GET request to each API endpoint and retries up to 3 times in case of an error.
   - The response status is sent to a channel for further processing.

2. **Data Writing**: 
   - The `csvWritter` function reads the status from the API fetching stage, fetches the data again if the API call was successful, and writes the structured data to a channel.
   - The parsed data (ID, Title, and Completed status) is printed to the console.

3. **Concurrency**: 
   - Goroutines and `sync.WaitGroup` are used to handle multiple API requests concurrently, improving performance.

## How to Run

1. Ensure that you have Go installed on your system.
2. Clone the repository:
   ```bash
   git clone https://github.com/your-username/api-fetcher-data-writer.git
