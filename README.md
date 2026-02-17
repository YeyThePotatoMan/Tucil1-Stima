# Tucil1-13524040
# Queens Puzzle Solver (Tucil Stima)

A GUI-based software to solve the Queens-Linkedin prblem with brute force.
This program allows sers to input puzzles via text  or image up to 26 colours and save the solution with .png or.txt format.
Live update is available for those who is curious with the process! :D
Built with **Go** and **Fyne** for GUI.

## Features
- GUI Interface for easy interaction.
- Input puzzle via text file (`.txt`).
- Input puzzle via image (`.txt`).
- Save solution as Image (`.png`) with grid visualization.
- Save solution as Text (`.txt`).
- Life update brute-force simulation.

## Requirements
1.  **Go** (go1.25.5 or higher)

## Installation
1.  Clone this repository:
    ```bash
    gh repo clone YeyThePotatoMan/Tucil1-13524040
    cd Tucil1-13524040
    ```

2.  Download dependencies:
    ```bash
    go mod tidy
    ```

3.  Run the application:
    ```bash
    cd src
    go run .
    ```

alternatively, you can use the executable:
1.  Clone this repository:
    ```bash
    gh repo clone YeyThePotatoMan/Tucil1-13524040
    cd Tucil1-13524040
    ```

2.  Download dependencies:
    ```bash
    go mod tidy
    ```

3. ```bash
    cd bin
    ./solver
    ```

## How to Use
1.  **Load Puzzle**:
    * Click **"Load File"** to upload a `.txt` file representing the grid.
    * (Bonus) Click **"Upload Image"** to upload a screenshot of the puzzle.
2.  **Solve**: Click the **"Solve"** button.
3.  **Save**:
    * Enter a filename (e.g., `solution1`).
    * Click **"Save Image"** to generate a visual solution in the `test/` folder.
    * Click **"Save Txt"** to generate a text solution in the `test/` folder.

## Project Structure
```bash
├── assets            # Assets used in the source code
├── bin/              # Location of the executable file
├── src/              # Source codes
├── test/             # Output folder for solutions
├── testcase/         # Optional folder for test cases
├── go.mod
├── go.sum
└── README.md
```

## Author
* **Name**: Kloce Paul William Saragih
* **Student ID**: 13524040
* **Class**: IF2211 Strategi Algoritma - K2
```bash
  __
<(o )___
 ( ._> /
  `---'
```