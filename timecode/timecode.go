package timecode

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// rate represents frame rate.
type rate struct {
	fps            int
	numerator      int32
	denominator    int32
	dropFrames     int
	framesPer1Min  int
	framesPer10Min int
}

var (
	// supportedRates represents supported frame rates 23.976, 24, 25, 29.97DF, 30, 48, 50, 59.94DF, 60.
	supportedRates = []*rate{
		{fps: 10, numerator: 10, denominator: 1, dropFrames: 0, framesPer1Min: 10 * 60, framesPer10Min: 10 * 600},             // 10
		{fps: 15, numerator: 15, denominator: 1, dropFrames: 0, framesPer1Min: 15 * 60, framesPer10Min: 15 * 600},             // 15
		{fps: 24, numerator: 24000, denominator: 1001, dropFrames: 0, framesPer1Min: 24 * 60, framesPer10Min: 24 * 600},       // 23.976
		{fps: 24, numerator: 24, denominator: 1, dropFrames: 0, framesPer1Min: 24 * 60, framesPer10Min: 24 * 600},             // 24
		{fps: 25, numerator: 25, denominator: 1, dropFrames: 0, framesPer1Min: 25 * 60, framesPer10Min: 25 * 600},             // 25
		{fps: 30, numerator: 30000, denominator: 1001, dropFrames: 2, framesPer1Min: 30*60 - 2, framesPer10Min: 30*600 - 9*2}, // 29.97DF
		{fps: 30, numerator: 30, denominator: 1, dropFrames: 0, framesPer1Min: 30 * 60, framesPer10Min: 30 * 600},             // 30
		{fps: 48, numerator: 48, denominator: 1, dropFrames: 0, framesPer1Min: 48 * 60, framesPer10Min: 48 * 600},             // 48
		{fps: 50, numerator: 50, denominator: 1, dropFrames: 0, framesPer1Min: 50 * 60, framesPer10Min: 50 * 600},             // 50
		{fps: 60, numerator: 60000, denominator: 1001, dropFrames: 4, framesPer1Min: 60*60 - 4, framesPer10Min: 60*600 - 9*4}, // 59.94DF
		{fps: 60, numerator: 60, denominator: 1, dropFrames: 0, framesPer1Min: 60 * 60, framesPer10Min: 60 * 600},             // 60
	}

	// timecodePattern represents timecode pattern.
	timecodePattern = regexp.MustCompile(`^([01][0-9]|2[0-3])([p:;.,])([0-5][0-9])([p:;.,])([0-5][0-9])([:;.,])([0-5][0-9])$`)
)

var (
	ErrNilTimecode          = errors.New("nil timecode")           // error for nil timecode
	ErrUnsupportedFrameRate = errors.New("unsupported frame rate") // error for unsupported frame rate
	ErrMismatchFrameRate    = errors.New("mismatch frame rate")    // error for mismatch frame rate
	ErrUnderflowFrames      = errors.New("underflow frames")       // error for underflow frames
	ErrInvalidTimecode      = errors.New("invalid timecode")       // error for invalid timecode
	ErrTooManyFrames        = errors.New("too many frames")        // error for too many frames
)

// Timecode represents timecode.
type Timecode struct {
	optp TimecodeOptionParam
	r    *rate
	HH   uint64
	MM   uint64
	SS   uint64
	FF   uint64
}

// TimecodeOptionParam represents timecode option parameter.
type TimecodeOptionParam struct {
	Sep   string
	SepDF string
}

// TimecodeOption represents timecode option.
type TimecodeOption func(*TimecodeOptionParam)

// newTimecodeOptionParam returns new TimecodeOptionParam.
func newTimecodeOptionParam() TimecodeOptionParam {
	return TimecodeOptionParam{
		Sep:   ":",
		SepDF: ":",
	}
}

// applyTimecodeOption applies TimecodeOption to TimecodeOptionParam.
func (p *TimecodeOptionParam) applyTimecodeOption(opts ...TimecodeOption) {
	for _, opt := range opts {
		opt(p)
	}
}

// newRate returns new rate.
func newRate(num, den int32) (*rate, error) {
	fps := float64(num) / float64(den)
	for _, r := range supportedRates {
		if float64(r.numerator)/float64(r.denominator) == fps {
			return r, nil
		}
	}
	return nil, ErrUnsupportedFrameRate
}

// IsSupportedFrameRate returns whether frame rate is supported.
func IsSupportedFrameRate(num, den int32) bool {
	_, err := newRate(num, den)
	return err == nil
}

// IsRepresentableFrames returns whether frames is representable.
func IsRepresentableFrames(frames uint64, num, den int32) bool {
	r, err := newRate(num, den)
	if err != nil {
		return false
	}
	return r.isRepresentableFrames(frames)
}

// NewTimecode returns new Timecode.
func NewTimecode(frames uint64, num, den int32, opts ...TimecodeOption) (*Timecode, error) {
	r, err := newRate(num, den)
	if err != nil {
		return nil, err
	}

	p := newTimecodeOptionParam()
	p.applyTimecodeOption(opts...)

	tc, err := Reset(&Timecode{r: r, optp: p}, frames)
	if err != nil {
		return nil, err
	}
	return tc, nil
}

// ParseTimecode returns new Timecode from formatted string.
func ParseTimecode(s string, num, den int32) (*Timecode, error) {
	r, err := newRate(num, den)
	if err != nil {
		return nil, err
	}

	// pattern: HH Sep1 MM Sep2 SS Sep3 FF
	// match  : 1  2    3  4    5  6    7
	match := timecodePattern.FindStringSubmatch(s)
	if len(match) != 8 || match[2] != match[4] {
		return nil, ErrInvalidTimecode
	}

	hh, _ := strconv.Atoi(match[1])
	sep := match[2]
	mm, _ := strconv.Atoi(match[3])
	ss, _ := strconv.Atoi(match[5])
	sepDF := match[6]
	ff, _ := strconv.Atoi(match[7])

	if ff < r.dropFrames && mm%10 != 0 {
		ff = r.dropFrames
	}

	return &Timecode{
		r:    r,
		optp: TimecodeOptionParam{Sep: sep, SepDF: sepDF},
		HH:   uint64(hh),
		MM:   uint64(mm),
		SS:   uint64(ss),
		FF:   uint64(ff),
	}, nil
}

// Reset returns new Timecode from Timecode and frames.
func Reset(tc *Timecode, frames uint64) (*Timecode, error) {
	if tc == nil {
		return nil, ErrNilTimecode
	}

	new := *tc

	if !new.r.isRepresentableFrames(frames) {
		return nil, ErrTooManyFrames
	}

	d := frames / uint64(new.r.framesPer10Min)
	m := frames % uint64(new.r.framesPer10Min)
	df := uint64(new.r.dropFrames)
	f := frames + 9*df*d
	if m > df {
		f += df * ((m - df) / uint64(new.r.framesPer1Min))
	}

	fps := uint64(new.r.fps)
	new.FF = f % fps
	new.SS = f / fps % 60
	new.MM = f / (fps * 60) % 60
	new.HH = f / (fps * 3600)

	return &new, nil
}

// equal returns whether rate is equal.
func (r *rate) equal(other *rate) bool {
	if r == nil || other == nil {
		return false
	}
	return r.numerator == other.numerator && r.denominator == other.denominator
}

// isRepresentableFrames returns whether frames is representable.
func (r *rate) isRepresentableFrames(frames uint64) bool {
	return frames < uint64(24*6*r.framesPer10Min)
}

// Frames returns number of frames.
func (tc *Timecode) Frames() uint64 {
	var frames uint64
	frames += tc.HH * 3600 * uint64(tc.r.fps)
	frames += tc.MM * 60 * uint64(tc.r.fps)
	frames += tc.SS * uint64(tc.r.fps)
	frames += tc.FF

	framesPer10Min := uint64(tc.r.fps) * 60 * 10
	framesPer1Min := framesPer10Min / 10

	var df uint64
	df += (frames / framesPer10Min) * uint64(tc.r.dropFrames) * 9
	df += (frames % framesPer10Min) / framesPer1Min * uint64(tc.r.dropFrames)

	return frames - df
}

// Duration returns duration from zero-origin.
func (tc *Timecode) Duration() time.Duration {
	return time.Duration((float64(tc.Frames()) * float64(tc.r.denominator) / float64(tc.r.numerator)) * float64(time.Second))
}

// Framerate denominator.
func (tc *Timecode) FramerateDenominator() int32 {
	return tc.r.denominator
}

// Framerate numerator.
func (tc *Timecode) FramerateNumerator() int32 {
	return tc.r.numerator
}

// Add Timecode and Timecode and return new Timecode.
func (tc *Timecode) Add(other *Timecode) (*Timecode, error) {
	if !tc.r.equal(other.r) {
		return nil, ErrMismatchFrameRate
	}
	return Reset(tc, tc.Frames()+other.Frames())
}

// Sub Timecode and Timecode and return new Timecode.
func (tc *Timecode) Sub(other *Timecode) (*Timecode, error) {
	if !tc.r.equal(other.r) {
		return nil, ErrMismatchFrameRate
	}
	if tc.Frames() < other.Frames() {
		return nil, ErrUnderflowFrames
	}
	return Reset(tc, tc.Frames()-other.Frames())
}

// Add Timecode and frames and return new Timecode.
func (tc *Timecode) AddFrames(frames uint64) (*Timecode, error) {
	return Reset(tc, tc.Frames()+frames)
}

// Sub Timecode and frames and return new Timecode.
func (tc *Timecode) SubFrames(frames uint64) (*Timecode, error) {
	if tc.Frames() < frames {
		return nil, ErrUnderflowFrames
	}
	return Reset(tc, tc.Frames()-frames)
}

// String returns Timecode formatted string.
// e.g. 01:23:45:28
func (tc *Timecode) String() string {
	sep := tc.optp.Sep
	lastSep := sep
	if tc.r.dropFrames > 0 {
		lastSep = tc.optp.SepDF
	}
	return fmt.Sprintf(
		"%02d%s%02d%s%02d%s%02d",
		tc.HH,
		sep,
		tc.MM,
		sep,
		tc.SS,
		lastSep,
		tc.FF,
	)
}
