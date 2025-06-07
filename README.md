# KiraForm Service

![Go](https://img.shields.io/badge/Go-1.x-blue) ![Open Source](https://img.shields.io/badge/Open%20Source-%E2%9C%94-green)

## 🚀 Overview
Form generator is a powerful and flexible backend service built with **Go** for dynamic form generation. Designed with scalability, security, and extensibility in mind, this open-source project enables seamless creation and management of customizable forms.

We welcome contributions from the community to enhance features, optimize performance, and expand integrations. 🎉

---

## ✨ Features
- Dynamic form generation with structured JSON
- Secure and scalable API
- Role-based access control (RBAC)
- Multi-database support
- High-performance request handling
- Extensible plugin system
- Fully open-source and community-driven

---

## 📦 Installation & Setup
### Prerequisites
- Go 1.x installed
- PostgreSQL
- Git (for cloning the repository)

### Steps
```sh
# Clone the repository
git clone https://github.com/kiritoo9/kiraform-service.git
cd kiraform-service

# Install dependencies
go mod tidy

# Set local environment in your local machine
export MIGRATION=true # to run migration
export SEEDER=true # to run seeder

export ENV=development # to run app as development (it also decided to choose .env.development file as environment)
export ENV=production # to run app as production (same)

# Set up environment variables (see .env.example)
# you can create .env.development or .env.production file
# otherwise it won't be read at all

# Run the application
go run src/entry/main.go
```

---

## 🛠 Contribution Guide
We 💙 contributions! If you'd like to contribute:
1. Fork the repository
2. Create a feature branch (`git checkout -b feature-name`)
3. Commit your changes (`git commit -m "Added new feature"`)
4. Push to the branch (`git push origin feature-name`)
5. Open a Pull Request 🎉

For detailed contribution guidelines, check [CONTRIBUTING.md](CONTRIBUTING.md).

---

## 📄 License
This project is licensed under the **MIT License** – see the [LICENSE](LICENSE) file for details.

---

## 🌎 Community & Support
Have questions, suggestions, or ideas? Feel free to open an issue or join our discussions!

---

## ✒️ Author
**Kiritoo9**  
[GitHub](https://github.com/kiritoo9)

---

## 📄 Version
App Version 0.0.1<br />
Architecture Version 1.1