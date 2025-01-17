package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type SuperBlock struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_inodes_count int32
	S_free_blocks_count int32
	S_mtime             float32
	S_umtime            float32
	S_mnt_count         int32
	S_magic             int32
	S_inode_size        int32
	S_block_size        int32
	S_first_ino         int32
	S_first_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
	// Total: 68 bytes
}

// Serialize escribe la estructura SuperBlock en un archivo binario en la posición especificada
func (sb *SuperBlock) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura SuperBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize lee la estructura SuperBlock desde un archivo binario en la posición especificada
func (sb *SuperBlock) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Obtener el tamaño de la estructura SuperBlock
	sbSize := binary.Size(sb)
	if sbSize <= 0 {
		return fmt.Errorf("invalid SuperBlock size: %d", sbSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura SuperBlock
	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura SuperBlock
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, sb)
	if err != nil {
		return err
	}

	return nil
}

// serializeInodeAndBitmap serializa un inodo y actualiza el bitmap de inodos
func (sb *SuperBlock) serializeInodeAndBitmap(path string, inode *Inode, offset int32) error {
	if err := inode.Serialize(path, int64(offset)); err != nil {
		return err
	}
	if err := sb.UpdateBitmapInode(path); err != nil {
		return err
	}
	return nil
}

// serializeBlockAndBitmap serializa un bloque y actualiza el bitmap de bloques
func (sb *SuperBlock) serializeBlockAndBitmap(path string, block *FolderBlock, offset int32) error {
	if err := block.Serialize(path, int64(offset)); err != nil {
		return err
	}
	if err := sb.UpdateBitmapBlock(path); err != nil {
		return err
	}
	return nil
}

// updateSuperBlockAfterInodeAndBlock actualiza los contadores del superbloque después de serializar un inodo y un bloque
func (sb *SuperBlock) updateSuperBlockAfterInodeAndBlock() {
	sb.S_inodes_count++
	sb.S_free_inodes_count--
	sb.S_first_ino += sb.S_inode_size
	sb.S_blocks_count++
	sb.S_free_blocks_count--
	sb.S_first_blo += sb.S_block_size
}

// createUsersFile crea el archivo users.txt y actualiza los inodos y bloques correspondientes
func (sb *SuperBlock) createUsersFile(path, usersText string) error {
	// Crear una instancia de Inode para el inodo raíz
	rootInode := &Inode{}
	// Deserializar el inodo raíz desde el archivo
	if err := rootInode.Deserialize(path, int64(sb.S_inode_start)); err != nil {
		return err
	}
	// Actualizar el tiempo de acceso del inodo raíz
	rootInode.I_atime = float32(time.Now().Unix())
	// Serializar el inodo raíz de vuelta al archivo
	if err := rootInode.Serialize(path, int64(sb.S_inode_start)); err != nil {
		return err
	}

	// Crear una instancia de FolderBlock para el bloque de carpeta raíz
	rootBlock := &FolderBlock{}
	// Deserializar el bloque de carpeta raíz desde el archivo
	if err := rootBlock.Deserialize(path, int64(sb.S_block_start)); err != nil {
		return err
	}
	// Actualizar el contenido del bloque de carpeta raíz para incluir users.txt
	rootBlock.B_content[2] = FolderContent{B_name: [12]byte{'u', 's', 'e', 'r', 's', '.', 't', 'x', 't'}, B_inodo: sb.S_inodes_count}
	// Serializar el bloque de carpeta raíz de vuelta al archivo
	if err := rootBlock.Serialize(path, int64(sb.S_block_start)); err != nil {
		return err
	}

	// Crear una instancia de Inode para el archivo users.txt
	usersInode := &Inode{
		I_uid:   1,
		I_gid:   1,
		I_size:  int32(len(usersText)),
		I_atime: float32(time.Now().Unix()),
		I_ctime: float32(time.Now().Unix()),
		I_mtime: float32(time.Now().Unix()),
		I_block: [15]int32{sb.S_blocks_count, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  [1]byte{'1'},
		I_perm:  [3]byte{'7', '7', '7'},
	}

	// Serializar el inodo del archivo users.txt y actualizar el bitmap de inodos
	if err := sb.serializeInodeAndBitmap(path, usersInode, sb.S_first_ino); err != nil {
		return err
	}

	// Crear una instancia de FileBlock para el contenido del archivo users.txt
	usersBlock := &FileBlock{}
	// Copiar el contenido del texto de users.txt al bloque de archivo
	copy(usersBlock.B_content[:], usersText)
	// Serializar el bloque de archivo y actualizar el bitmap de bloques
	if err := sb.serializeBlockAndBitmap(path, rootBlock, sb.S_first_blo); err != nil {
		return err
	}

	// Actualizar el superbloque después de serializar el inodo y el bloque
	sb.updateSuperBlockAfterInodeAndBlock()

	return nil
}

// CreateFolder crea una carpeta en el sistema de archivos
func (sb *SuperBlock) CreateFolder(path string, parentsDir []string, destDir string) error {
	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		return sb.createFolderInInode(path, 0, parentsDir, destDir)
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ {
		err := sb.createFolderInInode(path, i, parentsDir, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateFile crea un archivo en el sistema de archivos
func (sb *SuperBlock) CreateFile(path string, parentsDir []string, destFile string, size int, cont []string) error {

	// Si parentsDir está vacío, solo trabajar con el primer inodo que sería el raíz "/"
	if len(parentsDir) == 0 {
		return sb.createFileInInode(path, 0, parentsDir, destFile, size, cont)
	}

	// Iterar sobre cada inodo ya que se necesita buscar el inodo padre
	for i := int32(0); i < sb.S_inodes_count; i++ {
		err := sb.createFileInInode(path, i, parentsDir, destFile, size, cont)
		if err != nil {
			return err
		}
	}

	return nil
}

// Imprimir inodos
func (sb *SuperBlock) PrintInodes(path string) error {
	// Imprimir inodos
	fmt.Println("\nInodos\n----------------")
	// Iterar sobre cada inodo
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Imprimir el inodo
		fmt.Printf("\nInodo %d:\n", i)
		inode.Print()
	}

	return nil
}

// Impriir bloques
func (sb *SuperBlock) PrintBlocks(path string) error {
	// Imprimir bloques
	fmt.Println("\nBloques\n----------------")
	// Iterar sobre cada inodo
	for i := int32(0); i < sb.S_inodes_count; i++ {
		inode := &Inode{}
		// Deserializar el inodo
		err := inode.Deserialize(path, int64(sb.S_inode_start+(i*sb.S_inode_size)))
		if err != nil {
			return err
		}
		// Iterar sobre cada bloque del inodo (apuntadores)
		for _, blockIndex := range inode.I_block {
			// Si el bloque no existe, salir
			if blockIndex == -1 {
				break
			}
			// Si el inodo es de tipo carpeta
			if inode.I_type[0] == '0' {
				block := &FolderBlock{}
				// Deserializar el bloque
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Imprimir el bloque
				fmt.Printf("\nBloque %d:\n", blockIndex)
				block.Print()
				continue

				// Si el inodo es de tipo archivo
			} else if inode.I_type[0] == '1' {
				block := &FileBlock{}
				// Deserializar el bloque
				err := block.Deserialize(path, int64(sb.S_block_start+(blockIndex*sb.S_block_size))) // 64 porque es el tamaño de un bloque
				if err != nil {
					return err
				}
				// Imprimir el bloque
				fmt.Printf("\nBloque %d:\n", blockIndex)
				block.Print()
				continue
			}

		}
	}

	return nil
}
