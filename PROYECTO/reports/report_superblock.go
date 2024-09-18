package reports

import (
	"PROYECTO/structures"
	"PROYECTO/utils"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ReportSuperBlock genera un reporte del SuperBloque y lo guarda en la ruta especificada
func ReportSuperBlock(sb *structures.SuperBlock, path string, diskName string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Iniciar el contenido DOT para crear la tabla del SuperBloque
	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0" cellpadding="4">
                <tr><td colspan="2" bgcolor="green" color="white">Reporte de SUPERBLOQUE</td></tr>
                <tr><td>sb_nombre_hd</td><td>%s</td></tr>
                <tr><td>sb_arbol_virtual_count</td><td>%d</td></tr>
                <tr><td>sb_detalle_directorio_count</td><td>%d</td></tr>
                <tr><td>sb_inodos_count</td><td>%d</td></tr>
                <tr><td>sb_bloques_count</td><td>%d</td></tr>
                <tr><td>sb_arbol_virtual_free</td><td>%d</td></tr>
                <tr><td>sb_detalle_directorio_free</td><td>%d</td></tr>
                <tr><td>sb_inodos_free</td><td>%d</td></tr>
                <tr><td>sb_bloques_free</td><td>%d</td></tr>
                <tr><td>sb_date_creacion</td><td>%s</td></tr>
                <tr><td>sb_date_ultimo_montaje</td><td>%s</td></tr>
                <tr><td>sb_montajes_count</td><td>%d</td></tr>
                <tr><td>sb_ap_bitmap_arbol_directorio</td><td>%d</td></tr>
                <tr><td>sb_ap_arbol_directorio</td><td>%d</td></tr>
                <tr><td>sb_ap_bitmap_detalle_directorio</td><td>%d</td></tr>
                <tr><td>sb_ap_detalle_directorio</td><td>%d</td></tr>
                <tr><td>sb_ap_bitmap_inodos</td><td>%d</td></tr>
                <tr><td>sb_ap_inodos</td><td>%d</td></tr>
                <tr><td>sb_ap_bitmap_bloques</td><td>%d</td></tr>
                <tr><td>sb_ap_bloques</td><td>%d</td></tr>
                <tr><td>sb_ap_log</td><td>%d</td></tr>
            </table>>] }`,
		diskName, // Ahora se usa el nombre del disco proporcionado
		sb.S_inodes_count, sb.S_blocks_count, sb.S_free_inodes_count, sb.S_free_blocks_count,
		sb.S_blocks_count, sb.S_free_inodes_count, sb.S_free_blocks_count, sb.S_blocks_count,
		time.Unix(int64(sb.S_mtime), 0).Format("2006-01-02 15:04:05"),
		time.Unix(int64(sb.S_umtime), 0).Format("2006-01-02 15:04:05"),
		sb.S_mnt_count, sb.S_bm_inode_start, sb.S_inode_start, sb.S_bm_block_start, sb.S_block_start,
		sb.S_bm_inode_start, sb.S_inode_start, sb.S_bm_block_start, sb.S_block_start, sb.S_magic)

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Ejecutar el comando Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	fmt.Println("Imagen del reporte de SuperBloque generada:", outputImage)
	return nil
}
