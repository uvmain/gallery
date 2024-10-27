package types

import (
	"time"
)

var ExposureModes = map[int]string{
	0: "unknown",
	1: "Manual",
	2: "Normal program",
	3: "Aperture priority",
	4: "Shutter priority",
	5: "Creative program",
	6: "Action program",
	7: "Portrait mode",
	8: "Landscape mode",
}

var WhiteBalanceModes = map[int]string{
	0:   "Unknown",
	1:   "Daylight",
	2:   "Fluorescent",
	3:   "Tungsten",
	4:   "Flash",
	9:   "Fine Weather",
	10:  "Cloudy Weather",
	11:  "Shade",
	12:  "Daylight Fluorescent",
	13:  "Day White Fluorescent",
	14:  "Cool White Fluorescent",
	15:  "White Fluorescent",
	17:  "Standard Light A",
	18:  "Standard Light B",
	19:  "Standard Light C",
	20:  "D55",
	21:  "D65",
	22:  "D75",
	23:  "D50",
	24:  "ISO Studio Tungsten",
	255: "Other Light Source",
}

var FlashModes = map[int]string{
	0:  "No Flash",
	1:  "Fired",
	5:  "Fired, Return not detected",
	7:  "Fired, Return detected",
	8:  "On, Did not fire",
	9:  "On, Fired",
	11: "On, Return not detected",
	15: "On, Return detected",
	16: "Off, Did not fire",
	20: "Off, Did not fire, Return not detected",
	24: "Auto, Did not fire",
	25: "Auto, Fired",
	29: "Auto, Fired, Return not detected",
	31: "Auto, Fired, Return detected",
	32: "No flash function",
	48: "Off, No flash function",
	65: "Fired, Red-eye reduction",
	69: "Fired, Red-eye reduction, Return not detected",
	71: "Fired, Red-eye reduction, Return detected",
	73: "On, Red-eye reduction",
	77: "On, Red-eye reduction, Return not detected",
	79: "On, Red-eye reduction, Return detected",
	80: "Off, Red-eye reduction",
	88: "Auto, Did not fire, Red-eye reduction",
	89: "Auto, Fired, Red-eye reduction",
	93: "Auto, Fired, Red-eye reduction, Return not detected",
	95: "Auto, Fired, Red-eye reduction, Return detected",
}

type ImageMetadata struct {
	Slug             string    `json:"slug"`
	FilePath         string    `json:"filePath"`
	FileName         string    `json:"fileName"`
	Title            string    `json:"title"`
	DateTaken        time.Time `json:"dateTaken"`
	DateUploaded     time.Time `json:"dateUploaded"`
	CameraMake       string    `json:"cameraMake"`
	CameraModel      string    `json:"cameraModel"`
	LensMake         string    `json:"lensMake"`
	LensModel        string    `json:"lensModel"`
	FStop            string    `json:"fStop"`
	ExposureTime     string    `json:"exposureTime"`
	FlashStatus      string    `json:"flashStatus"`
	FocalLength      string    `json:"focalLength"`
	ISO              string    `json:"iso"`
	ExposureMode     string    `json:"exposureMode"`
	WhiteBalance     string    `json:"whiteBalance"`
	WhiteBalanceMode string    `json:"whiteBalanceMode"`
}

type MetadataFile struct {
	Slug     string
	FilePath string
	FileName string
}
