# 50M A
mkdisk -size=50 -unit=M -fit=FF -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia
# PRIMARIA 10M
fdisk -type=P -unit=K -name=Part11 -size=10 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -fit=BF

#DISCO 1
#761A -> 76 sus ultimos dos digitos del carnet
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -name=Part11
#REPORTE DISK
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r1_disk.jpg -name=disk
#REPORTE MBR
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r2_mbr.jpg -name=mbr
#-----------------5. MKFS-----------------
mkfs -id=761A
