// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

type (
	// VolumeSlice is the yaml representation of
	// the volumes block for a step in a pipeline.
	VolumeSlice []*Volume

	// Volume is the yaml representation of a volume
	// from a volumes block for a step in a pipeline.
	Volume struct {
		Source      string `yaml:"source,omitempty"      json:"source,omitempty"      jsonschema:"required,minLength=1,description=Set the source directory to be mounted.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-volume-key"`
		Destination string `yaml:"destination,omitempty" json:"destination,omitempty" jsonschema:"required,minLength=1,description=Set the destination directory for the mount in the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-volume-key"`
		AccessMode  string `yaml:"access_mode,omitempty" json:"access_mode,omitempty" jsonschema:"default=ro,description=Set the access mode for the mounted volume.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-volume-key"`
	}
)

// ToPipeline converts the VolumeSlice type
// to a pipeline VolumeSlice type.
func (v *VolumeSlice) ToPipeline() *pipeline.VolumeSlice {
	// volume slice we want to return
	volumes := new(pipeline.VolumeSlice)

	// iterate through each element in the volume slice
	for _, volume := range *v {
		// append the element to the pipeline volume slice
		*volumes = append(*volumes, &pipeline.Volume{
			Source:      volume.Source,
			Destination: volume.Destination,
			AccessMode:  volume.AccessMode,
		})
	}

	return volumes
}

// UnmarshalYAML implements the Unmarshaler interface for the VolumeSlice type.
func (v *VolumeSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// string slice we try unmarshalling to
	stringSlice := new(raw.StringSlice)

	// attempt to unmarshal as a string slice type
	err := unmarshal(stringSlice)
	if err == nil {
		// iterate through each element in the string slice
		for _, volume := range *stringSlice {
			// split each slice element into source, destination and access mode
			parts := strings.Split(volume, ":")

			switch {
			case len(parts) == 1:
				// append the element to the volume slice
				*v = append(*v, &Volume{
					Source:      parts[0],
					Destination: parts[0],
					AccessMode:  "ro",
				})

				continue
			case len(parts) == 2:
				// append the element to the volume slice
				*v = append(*v, &Volume{
					Source:      parts[0],
					Destination: parts[1],
					AccessMode:  "ro",
				})

				continue
			case len(parts) == 3:
				// append the element to the volume slice
				*v = append(*v, &Volume{
					Source:      parts[0],
					Destination: parts[1],
					AccessMode:  parts[2],
				})

				continue
			default:
				return fmt.Errorf("volume %s must contain at least 1 but no more than 2 `:`(colons)", volume)
			}
		}

		return nil
	}

	// volume slice we try unmarshalling to
	volumes := new([]*Volume)

	// attempt to unmarshal as a volume slice type
	err = unmarshal(volumes)
	if err != nil {
		return err
	}

	// iterate through each element in the volume slice
	for _, volume := range *volumes {
		// implicitly set `destination` field if empty
		if len(volume.Destination) == 0 {
			volume.Destination = volume.Source
		}

		// implicitly set `access_mode` field if empty
		if len(volume.AccessMode) == 0 {
			volume.AccessMode = "ro"
		}
	}

	// overwrite existing VolumeSlice
	*v = VolumeSlice(*volumes)

	return nil
}

func (v *Volume) ToYAML() *yaml.Volume {
	if v == nil {
		return nil
	}

	return &yaml.Volume{
		Source:      v.Source,
		Destination: v.Destination,
		AccessMode:  v.AccessMode,
	}
}

func (v *VolumeSlice) ToYAML() *yaml.VolumeSlice {
	// volume slice we want to return
	volumeSlice := new(yaml.VolumeSlice)

	// iterate through each element in the volume slice
	for _, volume := range *v {
		// append the element to the yaml volume slice
		*volumeSlice = append(*volumeSlice, volume.ToYAML())
	}

	return volumeSlice
}
