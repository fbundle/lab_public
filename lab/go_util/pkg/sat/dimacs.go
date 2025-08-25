package sat

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func Parse(r io.Reader) (formula Formula, err error) {
	numVariables := -1
	numClauses := 0

	var current []Literal

	scanner := bufio.NewScanner(r)
	read := 0
	for scanner.Scan() {
		raw := scanner.Bytes()

		if len(raw) == 0 {
			continue
		}

		if numVariables == -1 {
			switch raw[0] {
			case 'c':
			case 'p':
				fields := bytes.Fields(raw)
				if len(fields) != 4 {
					return nil, fmt.Errorf(
						"problem line should have 4 fields whitespace separated: %q", raw)
				}

				if string(fields[1]) != "cnf" {
					return nil, fmt.Errorf(
						"problem type must be 'cnf', got: %q", fields[1])
				}

				vars, err := strconv.Atoi(string(fields[2]))
				if err != nil {
					return nil, fmt.Errorf(
						"error converting variable count %q: %s", fields[2], err)
				}

				clauses, err := strconv.Atoi(string(fields[3]))
				if err != nil {
					return nil, fmt.Errorf(
						"error converting clauses count %q: %s", fields[3], err)
				}

				numVariables = vars
				numClauses = clauses

			default:
				return nil, fmt.Errorf(
					"invalid start of line character: %q", raw[0])
			}

			continue
		}

		fields := bytes.Fields(raw)

		end := false
		for _, raw := range fields {
			val, err := strconv.Atoi(string(raw))
			if err != nil {
				return nil, fmt.Errorf(
					"invalid literal %q", raw)
			}

			if val == 0 {
				end = true
				break
			}

			current = append(current, val)
		}

		if end {
			formula = append(formula, current)
			current = nil

			read++
			if read >= numClauses {
				break
			}
		}
	}

	return formula, nil
}
