package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

const version = "0.0.2"

func inputFromStdin() string {
	var text string
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		text += scan.Text() + "\n"
	}
	return text
}

func extractTextFromChildren(node ast.Node) [][]string {
	var texts [][]string
	for _, raw := range node.GetChildren() {
		rawText := extractTextFromTableDocument(raw)
		for _, text := range rawText {
			texts = append(texts, text)
		}
	}
	return texts
}

func extractTextFromTableDocument(node ast.Node) [][]string {
	switch node := node.(type) {
	case *ast.Document:
		// TODO: Deal with multi tables.
		for _, node := range node.GetChildren() {
			if val := extractTextFromChildren(node); val != nil {
				return val
			}
		}
		return nil
	case *ast.Table:
		texts := extractTextFromChildren(node)
		return texts
	case *ast.TableHeader:
		texts := extractTextFromChildren(node)
		return texts
	case *ast.TableBody:
		texts := extractTextFromChildren(node)
		return texts
	case *ast.TableRow:
		var row [][]string
		var ss []string
		for _, c := range node.GetChildren() {
			if len(c.GetChildren()) == 0 {
				ss = append(ss, "")
				continue
			}
			leaf := c.GetChildren()[0].AsLeaf()
			ss = append(ss, string(leaf.Literal))
		}
		row = append(row, ss)
		return row
	default:
		return nil
	}
}

func dumpCSV(records [][]string, buf *bytes.Buffer) {

	w := csv.NewWriter(buf)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			panic("Write Error")
		}
		w.Flush()
	}
}

func main() {
	input := inputFromStdin()
	inputByte := []byte(input)
	parser := parser.New()
	output := parser.Parse(inputByte)
	records := extractTextFromTableDocument(output)
	buf := new(bytes.Buffer)
	dumpCSV(records, buf)
	fmt.Printf(buf.String())
}
