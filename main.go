package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Entry struct {
	Plaats     int
	Naam       string
	Level      int
	Nodes      int
	ColorName  string
	Tier       int
	Commentaar string
	Foreground string
	Color      string
}

var leagues = []struct {
	Name       string
	Background string
	Foreground string
}{
	{"White", "#FFFFFF", "black"},
	{"Yellow", "#FFFF00", "black"},
	{"Salmon", "#FA8072", "black"},
	{"Orange", "#FFA500", "black"},
	{"Lime", "#00FF00", "black"},
	{"Green", "#008000", "white"},
	{"Cyan", "#00FFFF", "black"},
	{"Blue", "#0000FF", "white"},
	{"Dark Blue", "#00008B", "white"},
	{"Magenta", "#FF00FF", "white"},
	{"Purple", "#800080", "white"},
	{"Indigo", "#400040", "white"},
	{"Brown", "#8B4513", "white"},
	{"Red", "#FF0000", "white"},
	{"Dark Red", "#8B0000", "white"},
	{"Black", "#000000", "white"},
}

func getColorAndForeground(level int) (string, string) {
	tierIndex := (level - 1) % 16
	if tierIndex >= len(leagues) {
		tierIndex = len(leagues) - 1
	}
	return leagues[tierIndex].Name, leagues[tierIndex].Foreground
}

func getTier(level int) int {
	return ((level - 1) / 16) + 1
}

func getColorBackground(level int) string {
	tierIndex := (level - 1) % 16
	if tierIndex >= len(leagues) {
		tierIndex = len(leagues) - 1
	}
	return leagues[tierIndex].Background
}

func main() {
	// Open the input file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	var entries []Entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line by tabs
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			fmt.Println("Skipping invalid line:", line)
			continue
		}

		level, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Error parsing level:", err, "in line:", line)
			continue
		}
		nodes, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("Error parsing nodes:", err, "in line:", line)
			continue
		}
		comment := ""
		if len(parts) == 4 {
			comment = parts[3]
		}
		colorName, foreground := getColorAndForeground(level)
		colorBackground := getColorBackground(level)

		entries = append(entries, Entry{
			Naam:       parts[0],
			Level:      level,
			Nodes:      nodes,
			ColorName:  colorName,
			Tier:       getTier(level),
			Commentaar: comment,
			Foreground: foreground,
			Color:      colorBackground,
		})
	}

	// Sort entries by Level and then by Nodes
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Level == entries[j].Level {
			return entries[i].Nodes > entries[j].Nodes
		}
		return entries[i].Level > entries[j].Level
	})

	// Assign correct place values
	for i := range entries {
		entries[i].Plaats = i + 1
	}

	// Generate HTML
	tmpl := template.Must(template.New("report").Parse(htmlTemplate))
	outputFile, err := os.Create("output.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, entries)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	fmt.Println("HTML report generated successfully.")
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Level Report</title>
	<style>
		table { width: 100%; border-collapse: collapse; }
		th, td { padding: 8px; text-align: left; border: 1px solid #ddd; text-align: center; }
		th { background-color: #f2f2f2; }
	</style>
</head>
<body>
	<h1>Level Report</h1>
	<table>
		<tr>
			<th>Plaats</th>
			<th>Naam</th>
			<th>Level</th>
			<th>Color</th>
			<th>Tier</th>
			<th>Nodes</th>
			<th>Commentaar</th>
		</tr>
		{{range .}}
		<tr style="background-color: {{.Color}}; color: {{.Foreground}}">
			<td>{{.Plaats}}</td>
			<td>{{.Naam}}</td>
			<td>{{.Level}}</td>
			<td>{{.ColorName}}</td>
			<td>{{.Tier}}</td>
			<td>{{.Nodes}}</td>
			<td>{{.Commentaar}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>
`
