package arg

import (
	"fmt"
	"os"
	"strings"

	"mit-6.824-test"
)

type FlagLab struct {
	value string
	lab   *test.Lab
}

func NewFlagLab(value string) *FlagLab {
	return &FlagLab{value: value}
}

func (fl *FlagLab) String() string {
	return fl.value
}

func (fl *FlagLab) Set(val string) error {
	for _, choice := range test.Labs {
		if strings.ToLower(val) == strings.ToLower(choice.Name) {
			fl.value, fl.lab = choice.Name, &choice
			return nil
		}
	}
	var names []string
	for _, choice := range test.Labs {
		names = append(names, choice.Name)
	}
	return fmt.Errorf("invalid value '%s'. Choose from %s", val, strings.Join(names, ", "))
}

func (fl *FlagLab) GetLab() *test.Lab {
	return fl.lab
}

type FlagRun struct {
	lab   *test.Lab
	value string
}

func NewFlagRun(lab string) FlagRun {
	return FlagRun{lab: test.LabMap[lab]}
}

func (fr *FlagRun) GetLab() *test.Lab {
	return fr.lab
}

func (fr *FlagRun) String() string {
	return fr.value
}

func (fr *FlagRun) Set(val string) error {
	if fr.lab == nil || len(fr.lab.Tests) == 0 {
		for i, lab := range test.Labs {
			for _, choice := range lab.Tests {
				if strings.ToLower(val) == strings.ToLower(choice) {
					fr.lab = &test.Labs[i]
					fr.value = choice
					return nil
				}
			}
		}
		// 没找到test
		var tests []string
		for _, lab := range test.Labs {
			var names []string
			for _, choice := range lab.Tests {
				names = append(names, choice)
			}
			tests = append(tests, strings.Join(names, ", "))
		}
		return fmt.Errorf("invalid value '%s'. Choose from \n%s", val, strings.Join(tests, "\n"))
	} else {
		for _, choice := range fr.lab.Tests {
			if strings.ToLower(val) == strings.ToLower(choice) {
				fr.value = choice
				return nil
			}
		}
	}
	return fmt.Errorf("invalid value '%s'. Choose from %s", val, strings.Join(fr.lab.Tests, ", "))
}

type FlagEnvs []string

func (fe *FlagEnvs) String() string {
	return strings.Join(*fe, " ")
}

func (fe *FlagEnvs) Set(val string) error {
	*fe = append(*fe, val)
	env := strings.Split(val, "=")
	_ = os.Setenv(env[0], env[1])
	return nil
}
