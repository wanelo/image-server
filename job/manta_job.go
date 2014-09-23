package job

import (
	"fmt"

	"github.com/wanelo/image-server/uploader/manta/client"
)

var ImageExecutable = "/wanelo/public/images/bin/images-solaris-1.1.2"
var ImageInitScript = "/wanelo/stor/images/init.sh"

type JobCreator interface {
	CreateJob(opts client.CreateJobOpts) (string, error)
}

type MantaJob struct {
	BasePath  string
	Outputs   string
	Namespace string
}

func (j *MantaJob) ToImageCommand() string {
	return fmt.Sprintf("/assets%s --remote_base_path %s --namespace %s --outputs %s process $MANTA_INPUT_FILE", ImageExecutable, j.BasePath, j.Namespace, j.Outputs)
}

func (j *MantaJob) ToJobOpts() client.CreateJobOpts {
	phases := []client.Phase{
		{Type: "map",
			Exec: j.ToImageCommand(),
			Init: fmt.Sprintf("/assets%s", ImageInitScript),
			Assets: []string{
				ImageExecutable,
				ImageInitScript,
			},
		},
		{Type: "reduce", Exec: "cat"},
	}

	return client.CreateJobOpts{Phases: phases}
}

func (j *MantaJob) CreateMantaJob(mantaClient JobCreator) (string, error) {
	return mantaClient.CreateJob(j.ToJobOpts())
}
