package timecode

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRate(t *testing.T) {
	t.Run("NaN", func(t *testing.T) {
		_, err := newRate(1, 0)
		assert.Error(t, err)
	})
	t.Run("0fps", func(t *testing.T) {
		_, err := newRate(0, 1001)
		assert.Error(t, err)
	})
	t.Run("1fps", func(t *testing.T) {
		_, err := newRate(1, 1)
		assert.Error(t, err)
	})
	t.Run("23.976fps", func(t *testing.T) {
		r, err := newRate(24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, 24, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 24*60, r.framesPer1Min)
		assert.Equal(t, 24*600, r.framesPer10Min)
	})
	t.Run("24fps", func(t *testing.T) {
		r, err := newRate(24, 1)
		assert.NoError(t, err)
		assert.Equal(t, 24, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 24*60, r.framesPer1Min)
		assert.Equal(t, 24*600, r.framesPer10Min)
	})
	t.Run("25fps", func(t *testing.T) {
		r, err := newRate(25, 1)
		assert.NoError(t, err)
		assert.Equal(t, 25, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 25*60, r.framesPer1Min)
		assert.Equal(t, 25*600, r.framesPer10Min)
	})
	t.Run("29.97fps", func(t *testing.T) {
		r, err := newRate(30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, 30, r.fps)
		assert.Equal(t, 2, r.dropFrames)
		assert.Equal(t, 30*60-2, r.framesPer1Min)
		assert.Equal(t, 30*600-9*2, r.framesPer10Min)
	})
	t.Run("30fps", func(t *testing.T) {
		r, err := newRate(30, 1)
		assert.NoError(t, err)
		assert.Equal(t, 30, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 30*60, r.framesPer1Min)
		assert.Equal(t, 30*600, r.framesPer10Min)
	})
	t.Run("48fps", func(t *testing.T) {
		r, err := newRate(48, 1)
		assert.NoError(t, err)
		assert.Equal(t, 48, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 48*60, r.framesPer1Min)
		assert.Equal(t, 48*600, r.framesPer10Min)
	})
	t.Run("50fps", func(t *testing.T) {
		r, err := newRate(50, 1)
		assert.NoError(t, err)
		assert.Equal(t, 50, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 50*60, r.framesPer1Min)
		assert.Equal(t, 50*600, r.framesPer10Min)
	})
	t.Run("59.94fps", func(t *testing.T) {
		r, err := newRate(60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, 60, r.fps)
		assert.Equal(t, 4, r.dropFrames)
		assert.Equal(t, 60*60-4, r.framesPer1Min)
		assert.Equal(t, 60*600-9*4, r.framesPer10Min)
	})
	t.Run("60fps", func(t *testing.T) {
		r, err := newRate(60, 1)
		assert.NoError(t, err)
		assert.Equal(t, 60, r.fps)
		assert.Equal(t, 0, r.dropFrames)
		assert.Equal(t, 60*60, r.framesPer1Min)
		assert.Equal(t, 60*600, r.framesPer10Min)
	})
	t.Run("error/23.995fps", func(t *testing.T) {
		r, err := newRate(29995, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
	t.Run("error/23.997fps", func(t *testing.T) {
		r, err := newRate(29997, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
	t.Run("error/29.96fps", func(t *testing.T) {
		r, err := newRate(29960, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
	t.Run("error/29.98fps", func(t *testing.T) {
		r, err := newRate(29980, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
	t.Run("error/59.93fps", func(t *testing.T) {
		r, err := newRate(59930, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
	t.Run("error/59.95fps", func(t *testing.T) {
		r, err := newRate(59950, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
	t.Run("error/60.001fps", func(t *testing.T) {
		r, err := newRate(60001, 1000)
		assert.Equal(t, ErrUnsupportedFrameRate, err)
		assert.Nil(t, r)
	})
}

func TestNewTestcodeNonDF(t *testing.T) {
	t.Run("NaN", func(t *testing.T) {
		_, err := NewTimecode(1, 1, 0)
		assert.Error(t, err)
	})
	t.Run("0fps", func(t *testing.T) {
		_, err := NewTimecode(1, 0, 1001)
		assert.Error(t, err)
	})
	t.Run("1fps", func(t *testing.T) {
		_, err := NewTimecode(1, 1, 1)
		assert.Error(t, err)
	})
	t.Run("23.976fps", func(t *testing.T) {
		tc, err := NewTimecode(1439, 24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:23", tc.String())
		assert.Equal(t, uint64(1439), tc.Frames())
		assert.Equal(t, 60.018, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1440, 24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(1440), tc.Frames())

		tc, err = NewTimecode(1441, 24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(1441), tc.Frames())

		tc, err = NewTimecode(1440*10, 24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(1440*10), tc.Frames())
		assert.Equal(t, 600.6, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1440*10+1, 24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(1440*10+1), tc.Frames())

		maxFrames := uint64(24*6*(1440*10)) - 1
		tc, err = NewTimecode(maxFrames, 24000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:23", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())

		tc, err = NewTimecode(maxFrames+1, 24000, 1001)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("24fps", func(t *testing.T) {
		tc, err := NewTimecode(1439, 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:23", tc.String())
		assert.Equal(t, uint64(1439), tc.Frames())
		assert.Equal(t, 59.958, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1440, 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(1440), tc.Frames())

		tc, err = NewTimecode(1441, 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(1441), tc.Frames())

		tc, err = NewTimecode(1440*10, 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(1440*10), tc.Frames())
		assert.Equal(t, 600.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1440*10+1, 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(1440*10+1), tc.Frames())

		maxFrames := uint64(24*6*(1440*10)) - 1
		tc, err = NewTimecode(maxFrames, 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:23", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())
		assert.Equal(t, 86399.958, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(maxFrames+1, 24, 1)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("25", func(t *testing.T) {
		tc, err := NewTimecode(1499, 25, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:24", tc.String())
		assert.Equal(t, uint64(1499), tc.Frames())
		assert.Equal(t, 59.96, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1500, 25, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(1500), tc.Frames())

		tc, err = NewTimecode(1501, 25, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(1501), tc.Frames())

		tc, err = NewTimecode(1500*10, 25, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(1500*10), tc.Frames())
		assert.Equal(t, 600.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1500*10+1, 25, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(1500*10+1), tc.Frames())

		maxFrames := uint64(24*6*(1500*10)) - 1
		tc, err = NewTimecode(maxFrames, 25, 1)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:24", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())
		assert.Equal(t, 86399.96, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(maxFrames+1, 25, 1)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("30fps", func(t *testing.T) {
		tc, err := NewTimecode(1799, 30, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:29", tc.String())
		assert.Equal(t, uint64(1799), tc.Frames())
		assert.Equal(t, 59.967, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1800, 30, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(1800), tc.Frames())

		tc, err = NewTimecode(1801, 30, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(1801), tc.Frames())

		tc, err = NewTimecode(1800*10, 30, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(1800*10), tc.Frames())
		assert.Equal(t, 600.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(1800*10+1, 30, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(1800*10+1), tc.Frames())

		maxFrames := uint64(24*6*(1800*10)) - 1
		tc, err = NewTimecode(maxFrames, 30, 1)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:29", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())
		assert.Equal(t, 86399.967, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(maxFrames+1, 30, 1)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("48", func(t *testing.T) {
		tc, err := NewTimecode(2879, 48, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:47", tc.String())
		assert.Equal(t, uint64(2879), tc.Frames())
		assert.Equal(t, 59.979, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(2880, 48, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(2880), tc.Frames())

		tc, err = NewTimecode(2881, 48, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(2881), tc.Frames())

		tc, err = NewTimecode(2880*10, 48, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(2880*10), tc.Frames())
		assert.Equal(t, 600.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(2880*10+1, 48, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(2880*10+1), tc.Frames())

		maxFrames := uint64(24*6*(2880*10)) - 1
		tc, err = NewTimecode(maxFrames, 48, 1)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:47", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())
		assert.Equal(t, 86399.979, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(maxFrames+1, 48, 1)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("50fps", func(t *testing.T) {
		tc, err := NewTimecode(2999, 50, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:49", tc.String())
		assert.Equal(t, uint64(2999), tc.Frames())
		assert.Equal(t, 59.98, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(3000, 50, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(3000), tc.Frames())

		tc, err = NewTimecode(3001, 50, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(3001), tc.Frames())

		tc, err = NewTimecode(3000*10, 50, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(3000*10), tc.Frames())
		assert.Equal(t, 600.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(3000*10+1, 50, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(3000*10+1), tc.Frames())

		maxFrames := uint64(24*6*(3000*10)) - 1
		tc, err = NewTimecode(maxFrames, 50, 1)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:49", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())
		assert.Equal(t, 86399.98, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(maxFrames+1, 50, 1)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("60fps", func(t *testing.T) {
		tc, err := NewTimecode(0, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:00:00", tc.String())
		assert.Equal(t, uint64(0), tc.Frames())
		assert.Equal(t, 0.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(3599, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:59", tc.String())
		assert.Equal(t, uint64(3599), tc.Frames())

		tc, err = NewTimecode(3600, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:00", tc.String())
		assert.Equal(t, uint64(3600), tc.Frames())

		tc, err = NewTimecode(3601, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:01", tc.String())
		assert.Equal(t, uint64(3601), tc.Frames())

		tc, err = NewTimecode(3600*10, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(3600*10), tc.Frames())
		assert.Equal(t, 600.0, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(3600*10+1, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:01", tc.String())
		assert.Equal(t, uint64(3600*10+1), tc.Frames())

		maxFrames := uint64(24*6*(3600*10)) - 1
		tc, err = NewTimecode(maxFrames, 60, 1)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:59", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())
		assert.Equal(t, 86399.983, math.Round(tc.Duration().Seconds()*1000)/1000)

		tc, err = NewTimecode(maxFrames+1, 60, 1)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
}

func TestNewTestcodeDF(t *testing.T) {
	t.Run("30DF", func(t *testing.T) {
		tc, err := NewTimecode(1798, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:28", tc.String())
		assert.Equal(t, uint64(1798), tc.Frames())

		tc, err = NewTimecode(1799, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:29", tc.String())
		assert.Equal(t, uint64(1799), tc.Frames())

		tc, err = NewTimecode(1800, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:02", tc.String())
		assert.Equal(t, uint64(1800), tc.Frames())

		tc, err = NewTimecode(1800+1798*8, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:09:00:02", tc.String())
		assert.Equal(t, uint64(1800+1798*8), tc.Frames())

		tc, err = NewTimecode(1800+1798*9, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(1800+1798*9), tc.Frames())

		tc, err = NewTimecode(1800+1798*9+1799, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:59:29", tc.String())
		assert.Equal(t, uint64(1800+1798*9+1799), tc.Frames())

		tc, err = NewTimecode(1800+1798*9+1800, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:11:00:02", tc.String())
		assert.Equal(t, uint64(1800+1798*9+1800), tc.Frames())

		maxFrames := uint64(24*6*(1800+1798*9)) - 1
		tc, err = NewTimecode(maxFrames, 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:29", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())

		tc, err = NewTimecode(maxFrames+1, 30000, 1001)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
	t.Run("60DF", func(t *testing.T) {
		tc, err := NewTimecode(3596, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:56", tc.String())
		assert.Equal(t, uint64(3596), tc.Frames())

		tc, err = NewTimecode(3599, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:00:59:59", tc.String())
		assert.Equal(t, uint64(3599), tc.Frames())

		tc, err = NewTimecode(3600, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00:04", tc.String())
		assert.Equal(t, uint64(3600), tc.Frames())

		tc, err = NewTimecode(3600+3596*8, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:09:00:04", tc.String())
		assert.Equal(t, uint64(3600+3596*8), tc.Frames())

		tc, err = NewTimecode(3600+3596*9, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:00:00", tc.String())
		assert.Equal(t, uint64(3600+3596*9), tc.Frames())

		tc, err = NewTimecode(3600+3596*9+3599, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:10:59:59", tc.String())
		assert.Equal(t, uint64(3600+3596*9+3599), tc.Frames())

		tc, err = NewTimecode(3600+3596*9+3600, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:11:00:04", tc.String())
		assert.Equal(t, uint64(3600+3596*9+3600), tc.Frames())

		maxFrames := uint64(24*6*(3600+3596*9)) - 1
		tc, err = NewTimecode(maxFrames, 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "23:59:59:59", tc.String())
		assert.Equal(t, maxFrames, tc.Frames())

		tc, err = NewTimecode(maxFrames+1, 60000, 1001)
		assert.Equal(t, ErrTooManyFrames, err)
		assert.Nil(t, tc)
	})
}

func TestParseTimecode(t *testing.T) {
	t.Run("ParseTimecode", func(t *testing.T) {
		tc, err := ParseTimecode("00:01:00;00", 30000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:01:00;02", tc.String())
		assert.Equal(t, uint64(1800), tc.Frames())
		assert.Equal(t, ":", tc.optp.Sep)
		assert.Equal(t, ";", tc.optp.SepDF)
	})
	t.Run("ParseTimecode/19h", func(t *testing.T) {
		tc, err := ParseTimecode("19:00:00;00", 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1641600), tc.Frames())
		assert.Equal(t, ":", tc.optp.Sep)
		assert.Equal(t, ";", tc.optp.SepDF)
	})
	t.Run("ParseTimecode/23h", func(t *testing.T) {
		tc, err := ParseTimecode("23:00:00;00", 24, 1)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1987200), tc.Frames())
		assert.Equal(t, ":", tc.optp.Sep)
		assert.Equal(t, ";", tc.optp.SepDF)
	})
	t.Run("ParseTimecode/24h", func(t *testing.T) {
		tc, err := ParseTimecode("24:01:00;00", 24, 1)
		assert.Nil(t, tc)
		assert.Equal(t, ErrInvalidTimecode, err)
	})
	t.Run("ParseTimecode/skip timecode", func(t *testing.T) {
		tc, err := ParseTimecode("00:09:00:03", 60000, 1001)
		assert.NoError(t, err)
		assert.Equal(t, "00:09:00:04", tc.String())
	})
	t.Run("ParseTimecode/overflow", func(t *testing.T) {
		tc, err := ParseTimecode("00:09:00:99", 60000, 1001)
		assert.Error(t, err)
		assert.Nil(t, tc)
	})
	t.Run("ParseTimecode/lacks character", func(t *testing.T) {
		tc, err := ParseTimecode("0:01:00:00", 24, 1)
		assert.Nil(t, tc)
		assert.Equal(t, ErrInvalidTimecode, err)
	})
	t.Run("ParseTimecode/superfluous character", func(t *testing.T) {
		tc, err := ParseTimecode("0:01:000:00", 24, 1)
		assert.Nil(t, tc)
		assert.Equal(t, ErrInvalidTimecode, err)
	})
	t.Run("ParseTimecode/invalid character", func(t *testing.T) {
		tc, err := ParseTimecode("00:01:00:0x", 24, 1)
		assert.Nil(t, tc)
		assert.Equal(t, ErrInvalidTimecode, err)
	})
	t.Run("ParseTimecode/invalid separator", func(t *testing.T) {
		tc, err := ParseTimecode("00:01?00:0x", 24, 1)
		assert.Nil(t, tc)
		assert.Equal(t, ErrInvalidTimecode, err)
	})
	t.Run("ParseTimecode/separators are inconsistent", func(t *testing.T) {
		tc, err := ParseTimecode("00;01:00;00", 24, 1)
		assert.Nil(t, tc)
		assert.Equal(t, ErrInvalidTimecode, err)
	})
}

func TestReset(t *testing.T) {
	t.Run("Reset", func(t *testing.T) {
		const (
			maxFrames     = uint64(24*6*(3600+3596*9)) - 1
			tooManyFrames = maxFrames + 1
		)

		tc1, _ := NewTimecode(3596, 60000, 1001)
		tc2, _ := Reset(tc1, maxFrames)
		assert.Equal(t, "00:00:59:56", tc1.String())
		assert.Equal(t, "23:59:59:59", tc2.String())

		tcErr, err := Reset(tc2, tooManyFrames)
		assert.Equal(t, "23:59:59:59", tc2.String())
		assert.Nil(t, tcErr)
		assert.Equal(t, ErrTooManyFrames, err)
	})
}

func TestAdd(t *testing.T) {
	t.Run("Add timecode", func(t *testing.T) {
		tc1, _ := NewTimecode(1798, 30000, 1001)
		tc2, _ := NewTimecode(2, 30000, 1001)
		tc3, _ := tc1.Add(tc2)
		assert.Equal(t, "00:01:00:02", tc3.String())
	})
	t.Run("Add frames", func(t *testing.T) {
		tc1, _ := NewTimecode(17000, 30000, 1001)
		tc2, _ := tc1.AddFrames(982)
		assert.Equal(t, "00:10:00:00", tc2.String())
	})
	t.Run("Add/mismatch frame rate", func(t *testing.T) {
		tc1, _ := NewTimecode(1, 30, 1)
		tc2, _ := NewTimecode(1, 30000, 1001)
		tc3, err := tc1.Add(tc2)
		assert.Nil(t, tc3)
		assert.Equal(t, ErrMismatchFrameRate, err)
	})
	t.Run("Add/overflow", func(t *testing.T) {
		tc1, _ := NewTimecode(2589407, 30000, 1001)
		tc2, _ := NewTimecode(1, 30000, 1001)
		tc3, err := tc1.Add(tc2)
		assert.Nil(t, tc3)
		assert.Equal(t, ErrTooManyFrames, err)
	})
	t.Run("Add frames/overflow", func(t *testing.T) {
		tc1, _ := NewTimecode(2589407, 30000, 1001)
		tc2, err := tc1.AddFrames(1)
		assert.Nil(t, tc2)
		assert.Equal(t, ErrTooManyFrames, err)
	})
}

func TestSub(t *testing.T) {
	t.Run("Sub timecode", func(t *testing.T) {
		tc1, _ := NewTimecode(1800, 30000, 1001)
		tc2, _ := NewTimecode(1, 30000, 1001)
		tc3, _ := tc1.Sub(tc2)
		assert.Equal(t, "00:00:59:29", tc3.String())
	})
	t.Run("Sub frames", func(t *testing.T) {
		tc1, _ := NewTimecode(17982, 30000, 1001)
		tc2, _ := tc1.SubFrames(1798)
		assert.Equal(t, "00:09:00:02", tc2.String())
	})
	t.Run("Sub/mismatch frame rate", func(t *testing.T) {
		tc1, _ := NewTimecode(1, 24, 1)
		tc2, _ := NewTimecode(1, 24000, 1001)
		tc3, err := tc1.Sub(tc2)
		assert.Nil(t, tc3)
		assert.Equal(t, ErrMismatchFrameRate, err)
	})
	t.Run("Sub/underflow", func(t *testing.T) {
		tc1, _ := NewTimecode(1, 30000, 1001)
		tc2, _ := NewTimecode(10, 30000, 1001)
		tc3, err := tc1.Sub(tc2)
		assert.Nil(t, tc3)
		assert.Equal(t, ErrUnderflowFrames, err)
	})
	t.Run("Sub frames/underflow", func(t *testing.T) {
		tc1, _ := NewTimecode(10, 30000, 1001)
		tc2, err := tc1.SubFrames(11)
		assert.Nil(t, tc2)
		assert.Equal(t, ErrUnderflowFrames, err)
	})
}

func TestTimecodeOption(t *testing.T) {
	t.Run("single option/DF", func(t *testing.T) {
		opt := func(p *TimecodeOptionParam) {
			p.Sep = "."
			p.SepDF = ","
		}
		tc, err := NewTimecode(3596, 60000, 1001, opt)
		assert.NoError(t, err)
		assert.Equal(t, "00.00.59,56", tc.String())
	})
	t.Run("multiple options/DF", func(t *testing.T) {
		opt1 := func(p *TimecodeOptionParam) {
			p.Sep = ","
		}
		opt2 := func(p *TimecodeOptionParam) {
			p.SepDF = ";"
		}
		tc, err := NewTimecode(3596, 60000, 1001, opt1, opt2)
		assert.NoError(t, err)
		assert.Equal(t, "00,00,59;56", tc.String())
	})
	t.Run("multiple options/NDF", func(t *testing.T) {
		opt1 := func(p *TimecodeOptionParam) {
			p.Sep = "."
		}
		opt2 := func(p *TimecodeOptionParam) {
			p.SepDF = ";"
		}
		tc, err := NewTimecode(3596, 60, 1, opt1, opt2)
		assert.NoError(t, err)
		assert.Equal(t, "00.00.59.56", tc.String())
	})
}
