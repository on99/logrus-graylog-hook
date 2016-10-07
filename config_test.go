package graylog

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateConfig(t *testing.T) {
	Convey("Given config with invalid hostname", t, func() {
		c := Config{Hostname: ""}

		Convey("When validate", func() {
			err := ValidateConfig(c)

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given valid config", t, func() {
		c := NewConfig()

		Convey("When validate", func() {
			err := ValidateConfig(c)

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
