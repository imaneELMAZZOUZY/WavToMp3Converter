package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/app"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/models"
	"github.com/imaneELMAZZOUZY/WavToMp3Converter/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type cmdRunnerMock struct {
    mock.Mock
}

func (c *cmdRunnerMock) Run(name string, args []string) error {
    return c.Called(name, args).Error(0)
}

func createCmdRunnerMock(name string, args []string, err error) *cmdRunnerMock {
    c := &cmdRunnerMock{}
    c.On("Run", name, args).Return(err)
    return c
}

func TestJobConverterRun(t *testing.T) {

    mockTime := "2023-01-01T12:00:00Z"
	mockTimeProvider := func() time.Time {
		parsedTime, _ := time.Parse(time.RFC3339, mockTime)
		return parsedTime
	}
    
    type fields struct {
        FileChecker services.FileChecker
        CmdRunner   services.CmdRunner
    }
    testCases := []struct {
        title      string
        fields     fields
        jsonConfig models.ConversionConfig
        expected   models.ConversionRecord
        assertion  assert.ErrorAssertionFunc
    }{
        {
            title: "happy path",
            fields: fields{
                FileChecker: func(path string) bool {
                    return true
                },
                CmdRunner: createCmdRunnerMock(
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
                StartTime:        "_start_time_",
                EndTime:          mockTime,
                ConversionStatus: "successful",
            },
            assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
                return assert.Nil(tt, err)
            },
        },
        {
            title: "Conversion failed",
            fields: fields{
                FileChecker: func(path string) bool {
                    return true
                },
                CmdRunner: createCmdRunnerMock(
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
                StartTime:        "_start_time_",
                EndTime:          mockTime,
                ConversionStatus: "failed",
            },
            assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
                return assert.Nil(tt, err)
            },
        },
        {
            title: "FFmpeg not found error",
            fields: fields{
                FileChecker: func(path string) bool { return false },
                CmdRunner:   &cmdRunnerMock{},
            },
            jsonConfig: models.ConversionConfig{
                InputFile:  "__input_file__",
                OutputFile: "__output_file__",
                Codec:      "__codec__",
                Bitrate:    "__bitrate__",
                SampleRate: "__samplerate__",
                Channels:   "__channels__",
            },
            expected: models.ConversionRecord{},
            assertion: func(tt assert.TestingT, err error, i ...interface{}) bool {
                return assert.ErrorIs(tt, err, app.ErrFFMPEGNotFound)
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.title, func(t *testing.T) {
            sut := services.JobConverter{
                IsFileExist: tc.fields.FileChecker,
                CmdRunner:   tc.fields.CmdRunner,
                TimeProvider: mockTimeProvider,
            }

            res, err := sut.Run(tc.jsonConfig, tc.expected.StartTime)
            assert.Equal(t, tc.expected, res)
            tc.assertion(t, err)
        })
    }
}