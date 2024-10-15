# Gofluence
[![GoMod](https://img.shields.io/github/go-mod/go-version/beto20/gofluence)](https://github.com/beto20/gofluence)
[![Size](https://img.shields.io/github/languages/code-size/beto20/gofluence)](https://github.com/beto20/gofluence)
[![License](https://img.shields.io/github/license/beto20/gofluence)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/beto20/gofluence)](https://goreportcard.com/report/github.com/beto20/gofluence)

## Description
Gofluence is a CLI application and DevOps tool that helps you build documentation on Confluence pages. It also works well with CI/CD pipelines because it is easy to integrate and keeps your project documentation updated in a simple way.

- Currently support:
    - Java projects

## Use case
- CLI application
    - You can use it as a CLI application in Java projects, and it will generate documentation for your project.

- DevOps tool
    - You can integrate it into a CI/CD pipeline as a new step or an existing step, and it will generate documentation for your project.

## Getting Started
- Requires Go version 1.22 or 1.23

1. Clone the repository
    ```
    git clone https://github.com/beto20/Gofluence.git
    ```
2. Execute unit test
    ```
    go test
    ```
3. Build the project
    ```
    go build cmd/gofluence/main.go
    ```
4. Move the binary to the root of the new project
5. Execute the binary with the following flag
    ```
    ./main -b branch -c commit-hash -p prefix-project -t confluence-token -u confluence-url -rn repository-name -ct storage-container -cs storage-string-connection
    ```

## Contribute
Gofluence is an Open-Source Software (OSS), so if you would like to contribute with fixes, new features, integrations, or other improvements, I will guide you step by step.

1. You must check the feature backlog or ticket issues.
2. Choose an activity you would like to develop.
3. Fork the repository.
4. Create a new branch.
5. Add and commit your changes. Use a proper commit message.
6. Push to the branch.
7. Open a pull request and wait for the review and feedback.
8. If there are any observations or improvements suggested for your code, take them positively.
9. If there are no observations, your PR will be merged.
10. Congratulations! You have successfully contributed to one of the greatest Go OSS projects.