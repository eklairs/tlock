<img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white"> <img src="https://img.shields.io/github/stars/eklairs/tlock?style=for-the-badge" />

# TLock

TLock is an open-source tool to store and manage your authentication tokens securely. It gives users a centralized solution to generate and maintain time-based one-time passwords (TOTP) and other token types for secure authentication processes. By consolidating multiple authentication tokens into a single application, this 2FA manager streamlines the process of accessing accounts while ensuring high security.

## Features

- 📺  In Terminal
- ⌨️  Easily traverse through UI with keyboard
- 🔒  Proper encryption to secure your tokens at rest
- ⚡️  Supports Industry Standard TOTP and HOTP-based tokens
- ➕  Easily add tokens manually or from the screen (QR Code)
- 📁  Organize your tokens inside of folders
- 🎨  Themes to personalize tlock
- 🖼  Icon of the provider

>[!NOTE]
>For showing the provider's icon, you must have Nerd Fonts installed

## Installation

- **Arch Linux** (with AUR helper, like yay)

  ```fish
  yay -S tlock
  ```

- **MacOS** (with MacPorts)

  ```fish
  sudo port install tlock
  ```

- **Windows** (with scoop)

  ```fish
  scoop install tlock
  ```
- **Manually**

  You can also download the binary based on your operating system to use TLock from [releases](https://github.com/eklairs/tlock/releases)
- **go**

  ```fish
  go install github.com/eklairs/tlock@latest
  ```

## Screenshots

<img src="/assets/login.png" />
<img src="/assets/dashboard.png" />
<img src="/assets/help.png" />
<img src="/assets/add_token.png" />

## Running

Open your terminal and type `tlock` to start using tlock!

## Contributing

Did you come across a bug or want to introduce a new feature? Read the [CONTRIBUTING.md](https://github.com/eklairs/tlock/blob/main/CONTRIBUTING.md) and get started!

## License

[MIT](https://github.com/eklairs/tlock/raw/main/LICENSE)
