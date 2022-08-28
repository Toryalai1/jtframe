/*  This file is part of JTFRAME.
    JTFRAME program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    JTFRAME program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with JTFRAME.  If not, see <http://www.gnu.org/licenses/>.

    Author: Jose Tejada Gomez. Twitter: @topapate
    Version: 1.0
    Date: 7-7-2022 */

package jtcfgstr

import (
	//"text/template"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/jotego/jtframe/jtdef"
)

// appends arguments to a slice
func append_args(flag_name string, k *int, slice []string) []string {
	if *k < len(os.Args) && os.Args[*k] == flag_name {
		*k++
		added := 0
		for *k < len(os.Args) {
			if os.Args[*k][0] == '-' {
				*k--
				break
			}
			slice = append(slice, os.Args[*k])
			added++
			*k++
		}
		if added == 0 {
			log.Fatal("Expecting a macro name after ", flag_name)
		}
	}
	return slice
}

func parse_args(cfg jtdef.Config, args []string) {
	switch cfg.Target {
	case "mist", "mister", "sidi", "neptuno", "mc2", "mcp", "pocket":
		break
	default:
		{
			fmt.Printf("Unsupported target '%s'\n", cfg.Target)
			os.Exit(1)
		}
	}
	if len(cfg.Core) > 0 {
		cfg.Deffile = filepath.Join(os.Getenv("CORES"), cfg.Core, "/hdl/jt"+cfg.Core+".def")
	}
	if cfg.Verbose {
		fmt.Println("target=", cfg.Target)
		fmt.Println("def=", cfg.Deffile)
	}
	for k := 1; k < len(args); k++ {
		cfg.Discard = append_args("--undef", &k, cfg.Discard)
		cfg.Add = append_args("--def", &k, cfg.Add)
	}
	return
}

func make_cfgstr(cfg jtdef.Config, def map[string]string) (cfgstr string) {
	jtframe_path := os.Getenv("JTFRAME")
	if jtframe_path == "" {
		log.Fatal("Environment variable JTFRAME must be set")
	}
	var tpath string
	if cfg.Template == "" {
		tfolder := cfg.Target
		if cfg.Target == "sidi" { // SiDi shares the config string with MiST
			tfolder = "mist"
		}
		if cfg.Target == "pocket" { // Pocket doesn't have a config string
			return ""
		}
		tpath = jtframe_path + "/target/" + tfolder + "/cfgstr"
	} else {
		tpath = cfg.Template
	}
	t := template.Must(template.ParseFiles(tpath))
	var buffer bytes.Buffer
	t.Execute(&buffer, def)
	cfgstr = buffer.String()
	// Trim spaces
	chunks := strings.Split(cfgstr, ";")
	cfgstr = ""
	for _, s := range chunks {
		cfgstr = cfgstr + strings.TrimSpace(s) + ";"
	}
	// Removes any ; at the end
	for len(cfgstr) > 0 && cfgstr[len(cfgstr)-1] == ';' {
		cfgstr = cfgstr[0 : len(cfgstr)-1]
	}
	return
}

func dump_cfgstr(cfgstr string) {
	f, err := os.Create("cfgstr.hex")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for k, c := range cfgstr {
		fmt.Fprintf(f, "%02X ", c)
		if k&0xf == 0xf {
			fmt.Fprintln(f, "")
		}
	}
	for k := len(cfgstr); k < 1024; k++ {
		fmt.Fprintf(f, "00 ")
		if k&0xf == 0xf {
			fmt.Fprintln(f, "")
		}
	}
}

func dump_bash(def map[string]string) {
	for k, v := range def {
		fmt.Printf("export %s=\"%s\"\n", k, v)
	}
}

func dump_cpp(def map[string]string) {
	for k, v := range def {
		if k == "JTFRAME_PLL" {
			v = strings.TrimPrefix(v, "jtframe_pll")
		}
		// Get only the numerical part
		fmt.Printf("#define %s %s\n", k, v)
	}
}

func dump_verilog(def map[string]string, fmtstr string, esc_quotes bool) {
	pllsim := "10.416"
	for k, v := range def {
		if len(v) > 2 && v[0:2] == "0x" {
			val, _ := strconv.ParseInt(v, 0, 0)
			v = fmt.Sprintf("'h%X", val)
		}
		// Optionally escape quote characters
		if esc_quotes {
			v = strings.ReplaceAll(v, "\"", "\\\"")
		}
		// Output the key=value pair in the format
		// given by fmtstr, but skip it if the value
		// contains spaces, as simulators will get
		// confused
		if strings.Index(v, " ") == -1 {
			fmt.Printf(fmtstr+"\n", k, v)
		}
		if k == "JTFRAME_PLL" {
			// Converts to ns for simulation
			khz, err := strconv.Atoi(strings.TrimPrefix(v, "jtframe_pll"))
			if err != nil {
				log.Fatal("JTCFGSTR: while parsing JTFRAME_PLL ", nil)
			}
			ns := 1e6 / float32(khz*16)
			pllsim = fmt.Sprintf("%.3f", ns)
		}
	}
	// Output an extra macro used by fast_pll.v
	fmt.Printf(fmtstr+"\n", "JTFRAME_PLLSIM", pllsim)
}

func dump_parameter(def map[string]string, fmtstr string) {
	for k, v := range def {
		if !strings.HasPrefix(k, "JTFRAME_") {
			continue
		}
		if len(v) == 0 {
			v = "1"
		}
		if len(v) > 2 && v[0:2] == "0x" {
			val, _ := strconv.ParseInt(v, 0, 0)
			v = fmt.Sprintf("'h%X", val)
		}
		// Output the key=value pair in the format
		// given by fmtstr, but skip it if the value
		// contains spaces, as simulators will get
		// confused
		if strings.Index(v, " ") == -1 {
			fmt.Printf(fmtstr+"\n", k, v)
		}
	}
}

func Run(cfg jtdef.Config, args []string) {
	parse_args(cfg, args)
	def := jtdef.Make_macros(cfg)
	if !jtdef.Check_macros(def) {
		os.Exit(1)
	}
	switch cfg.Output {
	case "cfgstr":
		{
			// Make the config string
			cfgstr := make_cfgstr(cfg, def)
			dump_cfgstr(cfgstr)
			// show the config string
			if cfg.Verbose {
				fmt.Printf("jtdef.Config for target %s (%d bits)\n\n", cfg.Target, len(cfgstr)*8)
				fmt.Println(cfgstr, "\n\nBreak up:")
				aux := strings.Split(cfgstr, ";")
				for _, s := range aux {
					fmt.Println("\t", s)
				}
			}
		}
	case "cpp":
		dump_cpp(def)
	case "bash":
		dump_bash(def)
	case "quartus":
		dump_verilog(def, "set_global_assignment -name VERILOG_MACRO \"%s=%s\"", false)
		// dump_parameter(def, "set_parameter -name %s %s")
	case "iverilog", "verilator":
		dump_verilog(def, "+define+%s=%s", false) // do not escape quotes
	case "ncverilog", "synapticad", "modelsim", "questasim":
		dump_verilog(def, "+define+%s=%s", true) // escape quotes
	default:
		{
			fmt.Printf("JTCFGSTR: requested invalid output '%s'\n", cfg.Output)
			os.Exit(1)
		}
	}
}
