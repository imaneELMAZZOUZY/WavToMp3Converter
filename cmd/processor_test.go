package main

import (
	"errors"
	"testing"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cmdRunnerMock struct {
	mock.Mock
}

func (c *cmdRunnerMock) Run(name string, args []string) error {
	return c.Called(name, args).Error(0)
}

func createCmdRunnerMock(name string, args []string, err error) *cmdRunnerMock{
	c := &cmdRunnerMock{}

	c.On("Run", name, args).Return(err)

	return c
}

func TestJobConverterRun(t *testing.T) {
	type fields struct {
		fileChecker FileChecker
		cmdRunner   CmdRunner
	}
	testCases := []struct {
		title      string
		fields     fields
		jsonConfig models.ConversionConfig 
		expected   models.ConversionRecord
		assertion  assert.ErrorAssertionFunc
	}{
		{
			title : "happy path",
			fields: fields{
				fileChecker: func(path string) bool {
					return true
				},
				cmdRunner: createCmdRunnerMock(
					"bin\\ffmpeg.exe",
					[]string{
						"-i", "samples/Watched_folder/__input_file__",
						"-codec:a", "__codec__",
						"-b:a", "__bitrate__",
						"-ar", "__samplerate__",
						"-ac", "__channels__",
						"samples/Output_folder/__output_file__",
					},
					nil,
				),
			},
			jsonConfig: models.ConversionConfig{
				InputFile:  "__input_file__",
				OutputFile: "__output_file__",
				Codec:      "__codec__",
				Bitrate:    "__bitrate__",
				SampleRate: "__samplerate__",
				Channels:   "__channels__",
			},
			expected: models.ConversionRecord{
				InputFile:        "__input_file__",
				OutputFile:       "__output_file__",
				Codec:            "__codec__",
				Bitrate:          "__bitrate__",
				SampleRate:       "__samplerate__",
				Channels:         "__channels__",
				ConversionStatus: "successful",
			},
			assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(tt, err)
			},
	},
	{
		title : "Conversion failed",
		fields: fields{
			fileChecker : func(path string) bool {
				return true
			},
			cmdRunner : createCmdRunnerMock(
				"bin\\ffmpeg.exe",
				[]string{
					"-i", "samples/Watched_folder/__input_file__",
					"-codec:a", "__codec__",
					"-b:a", "__bitrate__",
					"-ar", "__samplerate__",
					"-ac", "__channels__",
					"samples/Output_folder/__output_file__",
				},
				errors.New("Conversion failed"),
			),
		},
		jsonConfig: models.ConversionConfig{
			InputFile:  "__input_file__",
				OutputFile: "__output_file__",
				Codec:      "__codec__",
				Bitrate:    "__bitrate__",
				SampleRate: "__samplerate__",
				Channels:   "__channels__",
		},
		expected: models.ConversionRecord{
			InputFile:        "__input_file__",
			OutputFile:       "__output_file__",
			Codec:            "__codec__",
			Bitrate:          "__bitrate__",
			SampleRate:       "__samplerate__",
			Channels:         "__channels__",
			ConversionStatus: "failed",
		},
		assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
			return assert.Nil(tt,err)
		},
	},
	{
		title : "FFmpeg not found error",
		fields: fields{
			fileChecker: func(path string) bool{ return false} ,
			cmdRunner: &cmdRunnerMock{},
		},
		jsonConfig: models.ConversionConfig{
			InputFile:  "__input_file__",
				OutputFile: "__output_file__",
				Codec:      "__codec__",
				Bitrate:    "__bitrate__",
				SampleRate: "__samplerate__",
				Channels:   "__channels__",
		},
		expected: models.ConversionRecord{} ,
		assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
			return assert.ErrorIs(tt, err, ErrFFMPEGNotFound)
		},
	},
		
	}
	
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			sut := JobConverter{
				IsFileExist: tc.fields.fileChecker,
				CmdRunner: tc.fields.cmdRunner,
			}	

			res, err := sut.Run(tc.jsonConfig, "")
			assert.Equal(t, tc.expected.InputFile, res.InputFile)
			assert.Equal(t, tc.expected.OutputFile, res.OutputFile)
			assert.Equal(t, tc.expected.Codec, res.Codec)
			assert.Equal(t, tc.expected.Bitrate, res.Bitrate)
			assert.Equal(t, tc.expected.SampleRate, res.SampleRate)
			assert.Equal(t, tc.expected.Channels, res.Channels)
			assert.Equal(t, tc.expected.ConversionStatus, res.ConversionStatus)
			tc.assertion(t, err)
		})
	}
}