package integer_range //nolint:revive

import "fmt"

// Config is the configuration for the integer_range resource.
type Config struct {
	Min  int64  `yaml:"min"`
	Max  int64  `yaml:"max"`
	Step int64  `yaml:"step"`
	Pad  uint32 `yaml:"pad"`
}

// validate ensures that the configuration is valid.
func (i *Config) validate() error {
	// Fix me now that we allow negative values
	if i.Max <= i.Min {
		return fmt.Errorf("max (%d) in integer_range must be greater than zero and the minimum value", i.Max)
	}

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler for defaulting the integer_range
func (i *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type ConfigDefaulted Config
	var defaults = ConfigDefaulted{
		Min:  DefaultIntegerRangeMin,
		Max:  DefaultIntegerRangeMax,
		Step: DefaultIntegerRangeStep,
		Pad:  DefaultIntegerRangePad,
	}

	out := defaults
	if err := unmarshal(&out); err != nil {
		return err
	}

	tmpl := Config(out)
	if err := tmpl.validate(); err != nil {
		return err
	}

	*i = tmpl
	return nil
}
