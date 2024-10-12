package reports

import (
	"backend/structures"
	"backend/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ReportDisk genera un reporte del disco y lo guarda en la ruta especificada
func ReportDisk(mbr *structures.MBR, path string) error {
	// Crear las carpetas padre si no existen
	err := utils.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utils.GetFileNames(path)

	// Calcular el tamaño total del disco
	totalSize := mbr.Mbr_size

	// Iniciar el contenido DOT
	dotContent := `digraph G {
rankdir=LR;
margin=0.1;
label="Reporte de Disco";
labelloc="t";
fontsize=30;
node [shape=plaintext];
n1 [label=<
<TABLE BORDER="1" CELLBORDER="1" CELLSPACING="3" CELLPADDING="10">
<TR>
<TD ROWSPAN="2" WIDTH="140" HEIGHT="100" BGCOLOR="lightblue" ALIGN="CENTER" VALIGN="MIDDLE">MBR</TD>
`

	for _, part := range mbr.Mbr_partitions {
		if part.Part_size == -1 {
			continue
		}

		partType := rune(part.Part_type[0])

		// Si es una partición extendida
		if partType == 'E' {
			// Calcular el porcentaje que ocupa la partición extendida
			extendedPercentage := float64(part.Part_size) / float64(totalSize) * 100

			if extendedPercentage > 100 {
				return fmt.Errorf("Error: El tamaño de la partición excede el tamaño total del disco")
			}

			dotContent += fmt.Sprintf(`<TD WIDTH="140" HEIGHT="100" BGCOLOR="lightyellow" ALIGN="CENTER" VALIGN="MIDDLE">Extendida %.2f%% del Disco</TD>`, extendedPercentage)
		} else if partType == 'P' {
			// Para particiones primarias
			percentage := float64(part.Part_size) / float64(totalSize) * 100

			if percentage > 100 {
				return fmt.Errorf("Error: El tamaño de la partición excede el tamaño total del disco")
			}

			partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
			dotContent += fmt.Sprintf(`<TD WIDTH="140" HEIGHT="100" BGCOLOR="lightgreen" ALIGN="CENTER" VALIGN="MIDDLE">%s %.2f%% del Disco</TD>`, partName, percentage)
		}
	}

	// Calcular el espacio libre
	freeSpace := totalSize
	for _, part := range mbr.Mbr_partitions {
		if part.Part_size != -1 {
			freeSpace -= int32(part.Part_size)
		}
	}

	freePercentage := float64(freeSpace) / float64(totalSize) * 100
	if freePercentage < 0 {
		return fmt.Errorf("Error: El tamaño de las particiones excede el tamaño total del disco")
	}

	// Agregar el espacio libre al contenido DOT
	dotContent += fmt.Sprintf(`<TD WIDTH="140" HEIGHT="100" BGCOLOR="white" ALIGN="CENTER" VALIGN="MIDDLE">Libre %.2f%% del Disco</TD>`, freePercentage)

	// Cerrar la tabla y el contenido DOT
	dotContent += `</TR></TABLE>>];
}`

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

	fmt.Println("Imagen del reporte de disco generada:", outputImage)
	return nil
}
