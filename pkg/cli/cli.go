package cli

import "time"

const (
	MonthLayout = "2006-01"
	DateLayout  = "2006-01-02"
)

type Date struct {
	Layout string
	time.Time
}

// Set implements pflag.Value
func (d *Date) Set(v string) (err error) {
	d.Time, err = time.Parse(d.Layout, v)
	return
}

// String implements pflag.Value
func (d *Date) String() string {
	return d.Format(d.Layout)
}

// Type implements pflag.Value
func (d *Date) Type() string {
	return "Date"
}
