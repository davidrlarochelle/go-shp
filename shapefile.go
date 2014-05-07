package goshp

import (
	"encoding/binary"
	"io"
)

type ShapeType int32

const (
	NULL        ShapeType = 0
	POINT                 = 1
	POLYLINE              = 3
	POLYGON               = 5
	MULTIPOINT            = 9
	POINTZ                = 11
	POLYLINEZ             = 13
	POLYGONZ              = 15
	MULTIPOINTZ           = 18
	POINTM                = 21
	POLYLINEM             = 23
	POLYGONM              = 25
	MULTIPOINTM           = 28
	MULTIPATCH            = 31
)

// Box structure made up from four coordinates. This type
// is used to represent bounding boxes
type Box struct {
	MinX, MinY, MaxX, MaxY float64
}

// Extends Box with coordinates from box
func (b *Box) Extend(box Box) {
	if box.MinX < b.MinX {
		b.MinX = box.MinX
	}
	if box.MinY < b.MinY {
		b.MinY = box.MinY
	}
	if box.MaxX > b.MaxX {
		b.MaxX = box.MaxX
	}
	if box.MaxY > b.MaxY {
		b.MaxY = box.MaxY
	}
}

// BBoxFromPoints returns the bounding box calculated
// from points.
func BBoxFromPoints(points []Point) (box Box) {
	for k, p := range points {
		if k == 0 {
			box = Box{p.X, p.Y, p.X, p.Y}
		} else {
			if p.X < box.MinX {
				box.MinX = p.X
			}
			if p.Y < box.MinY {
				box.MinY = p.Y
			}
			if p.X > box.MaxX {
				box.MaxX = p.X
			}
			if p.Y > box.MaxY {
				box.MaxY = p.Y
			}
		}
	}
	return
}

// Shape interface
type Shape interface {
	BBox() Box

	read(io.Reader)
	write(io.Writer)
}

// Shapefile NULL type
type Null struct {
}

// Returns the bounding box of the Null feature
func (n Null) BBox() Box {
	return Box{0.0, 0.0, 0.0, 0.0}
}

func (n *Null) read(file io.Reader) {
	binary.Read(file, binary.LittleEndian, n)
}

func (n *Null) write(file io.Writer) {
	binary.Write(file, binary.LittleEndian, n)
}

// Shapefile Point type
type Point struct {
	X, Y float64
}

// Returns the bounding box of the Point feature
func (p Point) BBox() Box {
	return Box{p.X, p.Y, p.X, p.Y}
}

func (p *Point) read(file io.Reader) {
	binary.Read(file, binary.LittleEndian, p)
}

func (p *Point) write(file io.Writer) {
	binary.Write(file, binary.LittleEndian, p)
}

// Shapefile PolyLine type
type PolyLine struct {
	Box
	NumParts  int32
	NumPoints int32
	Parts     []int32
	Points    []Point
}

// Returns the bounding box of the PolyLine feature
func (p PolyLine) BBox() Box {
	return BBoxFromPoints(p.Points)
}

func (p *PolyLine) read(file io.Reader) {
	binary.Read(file, binary.LittleEndian, &p.Box)
	binary.Read(file, binary.LittleEndian, &p.NumParts)
	binary.Read(file, binary.LittleEndian, &p.NumPoints)
	p.Parts = make([]int32, p.NumParts)
	p.Points = make([]Point, p.NumPoints)
	binary.Read(file, binary.LittleEndian, &p.Parts)
	binary.Read(file, binary.LittleEndian, &p.Points)
}

func (p *PolyLine) write(file io.Writer) {
	binary.Write(file, binary.LittleEndian, p.Box)
	binary.Write(file, binary.LittleEndian, p.NumParts)
	binary.Write(file, binary.LittleEndian, p.NumPoints)
	binary.Write(file, binary.LittleEndian, p.Parts)
	binary.Write(file, binary.LittleEndian, p.Points)
}

// Shapefile Polygon type
// The Polygon structure is identical to the PolyLine structure
type Polygon PolyLine

// Returns the bounding box of the Polygon feature
func (p Polygon) BBox() Box {
	return BBoxFromPoints(p.Points)
}

func (p *Polygon) read(file io.Reader) {
	binary.Read(file, binary.LittleEndian, &p.Box)
	binary.Read(file, binary.LittleEndian, &p.NumParts)
	binary.Read(file, binary.LittleEndian, &p.NumPoints)
	p.Parts = make([]int32, p.NumParts)
	p.Points = make([]Point, p.NumPoints)
	binary.Read(file, binary.LittleEndian, &p.Parts)
	binary.Read(file, binary.LittleEndian, &p.Points)
}

func (p *Polygon) write(file io.Writer) {
	binary.Write(file, binary.LittleEndian, p.Box)
	binary.Write(file, binary.LittleEndian, p.NumParts)
	binary.Write(file, binary.LittleEndian, p.NumPoints)
	binary.Write(file, binary.LittleEndian, p.Parts)
	binary.Write(file, binary.LittleEndian, p.Points)
}

// Field representation of a field object in the DBF file
type Field struct {
	Name      [11]byte
	Fieldtype byte
	Addr      [4]byte // not used
	Size      uint8
	Precision uint8
	Padding   [14]byte
}

// Returns a string representation of the Field. Currently
// this only returns field name.
func (f Field) String() string {
	return string(f.Name[:])
}

// Returns a StringField that can be used in SetFields to
// initialize the DBF file.
func StringField(name string, length uint8) Field {
	// TODO: Error checking
	field := Field{Fieldtype: 'C', Size: length}
	copy(field.Name[:], []byte(name))
	return field
}

// Returns a NumberField that can be used in SetFields to
// initialize the DBF file.
func NumberField(name string, length uint8) Field {
	field := Field{Fieldtype: 'N', Size: length}
	copy(field.Name[:], []byte(name))
	return field
}

// Returns a FloatField that can be used in SetFields to
// initialize the DBF file. Used to store floating points
// with precision in the DBF.
func FloatField(name string, length uint8, precision uint8) Field {
	field := Field{Fieldtype: 'F', Size: length, Precision: precision}
	copy(field.Name[:], []byte(name))
	return field
}

// Returns a DateField that can be used in SetFields to
// initialize the DBF file. Used to store Date strings
// formatted as YYYYMMDD. Data wise this is the same as
// a StringField with length 8.
func DateField(name string) Field {
	field := Field{Fieldtype: 'D', Size: 8}
	copy(field.Name[:], []byte(name))
	return field
}
