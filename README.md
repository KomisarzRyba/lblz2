# lblz2

Updated labeling system for UNCSA Percussion.

## Description

The `lblz2` project is an updated labeling system designed specifically for the UNCSA Percussion program. It aims to streamline and improve the process of labeling percussion instruments and equipment.

## Features

- **Efficient Labeling**: Simplified process for creating and managing labels.
- **Database Management**: Stays in sync with the Airtable database.
- **QR Code Generation**: Automatically generates QR codes for easy identification.
- **QR Code Printing**: Automatically prints the generated codes, meant to be used with an external label printer.
- **Go Language**: Entirely developed in Go for high performance and scalability.

## Installation

To install and run the project locally, follow these steps:

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/KomisarzRyba/lblz2.git
   cd lblz2
   ```

2. **Install Dependencies:**
   Ensure you have Go installed. Then, run:
   ```bash
   go mod tidy
   ```

3. **Run the Application:**
   ```bash
   go run main.go
   ```

## Usage

To use the application:

1. **Add Labels:**
   Utilize the provided interface to add new labels for instruments.
2. **Manage Labels:**
   Edit or remove QR labels as needed.
3. **Scan QR Codes:**
   Use a QR scanner to quickly access the items information. The generated code is displayed in the terminal window in ASCII.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries or feedback, please contact:

- **Repository Owner**: KomisarzRyba
- **Email**: [antek.olesik@gmail.com](mailto:antek.olesik@gmail.com)
