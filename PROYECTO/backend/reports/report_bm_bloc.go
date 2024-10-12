package reports

import (
	"backend/structures"
	"fmt"
	"os"
	"strings"
)

// GenerateBmBlockReport genera el reporte bm_bloc y lo guarda en un archivo de texto
func ReportBMBlock(sb *structures.SuperBlock, partitionPath string, reportPath string) error {
	// Leer el bitmap de bloques desde el archivo
	bitmap, err := readBlockBitmap(partitionPath, sb.S_bm_block_start, sb.S_blocks_count)
	if err != nil {
		return fmt.Errorf("error al leer el bitmap de bloques: %w", err)
	}

	// Formatear el bitmap en líneas de 20 registros
	formattedBitmap := formatBitmap(bitmap, 20)

	// Escribir el resultado en un archivo de texto
	err = writeToFile(reportPath, formattedBitmap)
	if err != nil {
		return fmt.Errorf("error al escribir el reporte bm_bloc: %w", err)
	}

	return nil
}

// readBlockBitmap lee el bitmap de bloques desde el archivo
func readBlockBitmap(path string, start int32, count int32) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bitmap := make([]byte, count)
	_, err = file.ReadAt(bitmap, int64(start))
	if err != nil {
		return nil, err
	}

	return bitmap, nil
}

// formatBitmap formatea el bitmap en líneas de n registros
func formatBitmap(bitmap []byte, n int) string {
	var sb strings.Builder
	for i, b := range bitmap {
		if i > 0 && i%n == 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%d", b))
	}
	return sb.String()
}

// writeToFile escribe el contenido en un archivo de texto
func writeToFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
