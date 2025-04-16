# HomeDrive a Photo Gallery made with Go + Gin

A simple server built in Go using the [Gin](https://github.com/gin-gonic/gin) web framework. It allows users to create directories, upload images, browse them, and delete selected ones — all through a clean and minimal web interface.

## Features

- Create new directories
- Upload multiple files to a selected directory
- Browse image content within a directory (with next/previous buttons)
- Delete selected images
- User interface built with HTML + CSS + vanilla JavaScript

## Technologies

- **Backend**: Go + Gin
- **Frontend**: HTML, CSS, JavaScript
- **Storage**: Local filesystem (`/upload/`)

## Installation & Running

1. **Clone the repository:**
2. **Run**
```bash
  go run server.go
```
3. **Open in browser [localhost](http://localhost:8080) on your local machine or `http://<ip of the machine that the code is running on>:8080` to access it from a device in the same local network
