package structures

import (
	"encoding/binary"
	"os"
)

type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
	Part_id     [16]byte
}

func (e *EBR) DeserializeEBR(path string, start int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Seek(start, 0)
	err = binary.Read(file, binary.LittleEndian, e)
	if err != nil {
		return err
	}

	return nil
}

func (e *EBR) SerializeEBR(path string, start int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Seek(start, 0)
	err = binary.Write(file, binary.LittleEndian, e)
	if err != nil {
		return err
	}

	return nil
}

func (e *EBR) Print() {
	println("Part_status:", string(e.Part_status))
	println("Part_fit:", string(e.Part_fit))
	println("Part_start:", e.Part_start)
	println("Part_size:", e.Part_size)
	println("Part_next:", e.Part_next)
	println("Part_name:", string(e.Part_name[:]))

}

func NewEBR() EBR {
	var eb EBR
	eb.Part_status = '0'
	eb.Part_size = 0
	eb.Part_next = -1
	return eb
}
