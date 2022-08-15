package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type keyVal map[string]string
type parser struct {
	data map[string]keyVal
}

func NewParser() parser {
	p := parser{}
	p.data = make(map[string]keyVal)
	return p
}

func (p *parser) String() string {
	ret := ""
	for section, pairs := range p.data {
		if section != "" {
			ret += "[" + section + "]\n"
		}
		for key, val := range pairs {
			ret += key + " = " + val + "\n"
		}
		ret += "\n"
	}
	return ret
}

func (p *parser) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := ""
	_ = p.AddSection(currentSection)

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \n\t")
		if line == "" {
			continue
		}
		if isSection(line) {
			currentSection = strings.Trim(line, "[]")
			_ = p.AddSection(currentSection)
		} else if isKeyVal(line) {
			before, after, _ := strings.Cut(line, "=")
			before = strings.Trim(before, " ")
			after = strings.Trim(after, " ")
			_ = p.AddKeyVal(currentSection, before, after)
		} else if line[0] != '#' {
			return fmt.Errorf("line \"%s\" is invalid", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func isSection(line string) bool {
	return line[0] == '[' && line[len(line)-1] == ']'
}

func isKeyVal(line string) bool {
	key, val, equalSign := strings.Cut(line, "=")
	if !equalSign {
		return false
	}
	key = strings.Trim(key, " ")
	if key == "" || strings.Contains(key, " ") {
		return false
	}

	val = strings.Trim(val, " ")
	if val == "" || strings.Contains(val, " ") {
		return false
	}

	return true
}

func (p *parser) AddSection(section string) error {
	if !isSection("[" + section + "]") {
		return fmt.Errorf("section %s is invalid", section)
	}
	if _, ok := p.data[section]; !ok {
		p.data[section] = make(map[string]string)
	}
	return nil
}

func (p *parser) AddKeyVal(section, key, val string) error {
	err := p.AddSection(section)
	if err != nil {
		return err
	}
	if !isKeyVal(key + " = " + val) {
		return fmt.Errorf("key value pair \"%s, %s\" is invalid", key, val)
	}
	p.data[section][key] = val
	return nil
}

func (p *parser) WriteToFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for section, pairs := range p.data {
		if section != "" {
			_, err := file.WriteString("[" + section + "]\n")
			if err != nil {
				return err
			}
		}
		for key, val := range pairs {
			_, err := file.WriteString(key + " = " + val + "\n")
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return err
		}

	}
	return nil
}
