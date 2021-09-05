package project


import (
	. "github.com/smartystreets/goconvey/convey"
	typhoon "github.com/vortex14/gotyphoon"
	"testing"
)

func TestCreateProject(t *testing.T) {

	Convey("Create a new project", t, func() {

		project := typhoon.Project{
			Name: "test-project",
		}

		Convey("checking name", func() {
			So(project.Name, ShouldEqual, "test-project")
		})

	})

}