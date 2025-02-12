#!/usr/bin/env python3
import sys

def format_text_file(input_file, output_file, extra_space=2):
    """
    Reads a whitespace-delimited text file, computes the maximum width for each column,
    and writes out a new file with fixed-width columns.
    
    :param input_file: Path to the input file.
    :param output_file: Path to the output file.
    :param extra_space: Extra spaces to add after each column (default is 2).
    """
    # Read non-blank lines from the file.
    with open(input_file, 'r') as f:
        lines = [line.rstrip("\n") for line in f if line.strip()]
    
    # Split each line into columns (assumes whitespace separation).
    data = [line.split() for line in lines]
    
    # Determine the number of columns (assumes all rows have roughly the same number)
    num_cols = max(len(row) for row in data)
    
    # Compute the maximum width for each column.
    col_widths = [0] * num_cols
    for row in data:
        for i, item in enumerate(row):
            col_widths[i] = max(col_widths[i], len(item))
    
    # Open the output file for writing.
    with open(output_file, 'w') as out:
        for row in data:
            # If a row has fewer items than num_cols, pad it with empty strings.
            if len(row) < num_cols:
                row += [''] * (num_cols - len(row))
            # Build the formatted row using ljust to pad each column.
            formatted_line = ""
            for i, item in enumerate(row):
                formatted_line += item.ljust(col_widths[i] + extra_space)
            out.write(formatted_line.rstrip() + "\n")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python format_txt.py input.txt output.txt")
    else:
        input_filename = sys.argv[1]
        output_filename = sys.argv[2]
        format_text_file(input_filename, output_filename)
        print(f"Formatted file written to {output_filename}")
