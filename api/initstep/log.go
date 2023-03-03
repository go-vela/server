package initstep

import (
	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"
)

// SaveLog is a helper function for APIs to create or update
// an InitStep and Log.
func SaveLog(c *gin.Context, b *library.Build, i *library.InitStep, l *library.Log) (*library.InitStep, *library.Log, error) {
	if len(l.GetData()) == 0 {
		return i, l, nil
	}

	var err error

	if i.GetNumber() == 0 {
		count, err := database.FromContext(c).CountInitStepsForBuild(b)
		if err != nil {
			return i, l, err
		}

		// this might be a naive way to increment the step number
		i.SetNumber(int(count + 1))

		// send API call to create the InitStep
		err = database.FromContext(c).CreateInitStep(i)
		if err != nil {
			return i, l, err
		}
	}

	// send API call to capture the InitStep (now with an ID)
	i, err = database.FromContext(c).GetInitStepForBuild(b, i.GetNumber())
	if err != nil {
		return i, l, err
	}

	l.SetInitStepID(i.GetID())

	if l.GetID() == 0 {
		// send API call to create the InitStep's Log
		err = database.FromContext(c).CreateLog(l)
		if err != nil {
			return i, l, err
		}
	} else {
		// send API call to update the InitStep's Log
		err = database.FromContext(c).UpdateLog(l)
		if err != nil {
			return i, l, err
		}
	}

	// send API call to capture the Log (now with an ID)
	l, err = database.FromContext(c).GetLogForInitStep(i)
	if err != nil {
		return i, l, err
	}

	return i, l, nil
}
