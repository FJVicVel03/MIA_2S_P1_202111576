#Calificacion Proyecto 1
#2S 2024
#Cambiar fernando-vicente -> por el usuario de su distribución de linux
#Cambiar “76” -> por los ultimos dos digitos de su carnet


#----------------- 1. mkdisk  -----------------


#----------------- mkdisk CON ERROR -----------------
# ERROR PARAMETROS
mkdisk -param=x -size=30 -path=/home/fernando-vicente/Calificacion_MIA/Discos/DiscoN.mia


#----------------- CREACION DE DISCOS -----------------
# ERROR PARAMETROS
mkdisk -tamaño=3000 -unit=K-path=/home/fernando-vicente/Calificacion_MIA/Discos/DiscoN.mia
# 50M A
mkdisk -size=50 -unit=M -fit=FF -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia
# 50M B
mkdisk -unit=K-size=51200 -fit=BF -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco2.mia
# 13M C
mkdisk -size=13 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia
# 50M D
mkdisk -size=51200 -unit=K-path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco4.mia
# 20M E
mkdisk -size=20 -unit=M -fit=WF -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia
# 50M F X
mkdisk -size=50 -unit=M -fit=FF -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco6.mia
# 50M G X
mkdisk -size=50 -unit=M -fit=FF -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco7.mia
# 50M H X
mkdisk -size=51200 -unit=K-path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco8.mia
# 50M I X
mkdisk -size=51200 -unit=K-path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco9.mia
# 50M J X
mkdisk -size=51200 -unit=K-path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco10.mia


#-----------------2. RMDISK-----------------
#ERROR DISCO NO EXISTE
rmdisk -path=/home/fernando-vicente/Calificacion_MIA/Discos/DiscoN.mia
# BORRANDO DISCO
rmdisk -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco6.mia
# BORRANDO DISCO
rmdisk -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco7.mia
# BORRANDO DISCO
rmdisk -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco8.mia
# BORRANDO DISCO
rmdisk -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco9.mia
# BORRANDO DISCO
rmdisk -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco10.mia


#-----------------3. FDISK-----------------
#-----------------CREACION DE PARTICIONES-----------------
#DISCO 1
# ERROR RUTA NO ENCONTRADA
fdisk -type=P -unit=K -name=PartErr -size=10485 -path=/home/fernando-vicente/Calificacion_MIA/Discos/DiscoN.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=K -name=Part11 -size=10485 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=K -name=Part12 -size=10240 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=M -name=Part13 -size=10 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -fit=BF
# PRIMARIA 10M
fdisk -type=P -unit=K -name=Part14 -size=10485 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -fit=BF
#ERR LMITE PARTICION PRIMARIA
#fdisk -type=P -unit=B -name=PartErr -size=10485760 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -fit=BF


# LIBRE DISCO 1: 50-4*10 = 10 -> 20%


#DISCO 3
# ERROR FALTA ESPACIO
fdisk -type=P -unit=M -name=PartErr -size=20 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia
#4M
fdisk -type=P -unit=M -name=Part31 -size=4 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia
#4M
fdisk -type=P -unit=M -name=Part32 -size=4 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia
#1M
fdisk -type=P -unit=M -name=Part33 -size=1 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia


#LIBRE DISCO 3: 13-9= 4 -> 30.77%


#DISCO 5
# 5MB
fdisk -type=E -unit=K- name=Part51 -size=5120 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=K -name=Part52 -size=1024 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 5MB
fdisk -type=P -unit=K -name=Part53 -size=5120 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=K -name=Part54 -size=1024 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=K -name=Part55 -size=1024 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -fit=BF
# 1MB
fdisk -type=L -unit=K -name=Part56 -size=1024 -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -fit=BF


# LIBRE DISCO 5: 20-10 = 5 -> 50%
# LIBRE EXTENDIDA 2: 5-4 = 1 -> 20% (por los EBR deberia ser menos)

#-----------------MOUNT-----------------
#-----------------MONTAR PARTICIONES-----------------
#DISCO 1
#761A -> 76 sus ultimos dos digitos del carnet
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -name=Part11
#762A -> 76 sus ultimos dos digitos del carnet
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -name=Part12
#ERROR PARTICION YA MONTADA
#mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco1.mia -name=Part11


#DISCO 3
#ERROR PARTCION NO EXISTE
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia -name=Part0
#763B -> 76 sus ultimos dos digitos del carnet
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia -name=Part31
#764B -> 76 sus ultimos dos digitos del carnet
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco3.mia -name=Part32


#DISCO 5
#765C -> 76 sus ultimos dos digitos del carnet
mount -path=/home/fernando-vicente/Calificacion_MIA/Discos/Disco5.mia -name=Part53

#-----------------REPORTES PARTE 1-----------------
#DISCO 1
#ERROR ID NO ENCONTRADO
rep -id=A761 -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_rE.jpg -name=mbr
#REPORTE DISK
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r1_disk.jpg -name=disk
#REPORTE MBR 
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r2_mbr.jpg -name=mbr


#DISCO 3
#ERROR ID NO ENCONTRADO
rep -id=763B -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_rE_mbr.jpg -name=mbr
#REPORTE DISK
rep -id=763B -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r3_disk.jpg -name=disk
#REPORTE MBR
rep -id=764B -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r4_disk.jpg -name=mbr


#DISCO 5
#ERROR ID NO ENCONTRADO
rep -id=IDx -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_rE_mbr.jpg -name=mbr
#REPORTE DISK
rep -id=765C -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r5_disk.jpg -name=disk
#REPORTE MBR
rep -id=765C -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p1_r6_mbr.jpg -name=mbr


#-----------------5. MKFS-----------------
mkfs -id=761A

#-----------------PARTE 3-----------------


#-----------------7. LOGIN-----------------
login -user=root -pass=123 -id=761A
#ERROR SESION INICIADA
login -user=root -pass=123 -id=761A

#-----------------15. MKDIR-----------------
mkdir -path=/bin
mkdir -path="/home/archivos/archivos 24"
mkdir -p -path=/home/archivos/user/docs/usac
mkdir -p -path=/home/archivos/carpeta1/carpeta2/carpeta3/carpeta4/carpeta5

#-----------------8. LOGOUT-----------------
logout
logout #ERROR NO HAY SESION INICIADA


#-----------------7. LOGIN-----------------
login -user=root -pass=123 -id=761A

#-----------------14. MKFILE-----------------
mkfile -path=/home/archivos/user/docs/Tarea.txt -size=75
mkfile -path=/home/archivos/user/docs/Tarea2.txt -size=768


#Para este comando hay que crear un archivo en la computadora y en cont poner su primer nombre
#Crear un archivo txt en su Escritorio llamado NAME


# Cambiar la ruta del cont por la del archivo NAME.txt que creo
mkfile -path=/home/archivos/user/docs/Tarea3.txt -size=10 -cont=/home/fernando-vicente/Calificacion_MIA/CONT/name.txt


#ERROR NO EXISTE RUTA
mkfile -path="/home/archivos/noexiste/b1.txt"

#ERROR NEGATIVO
mkfile -path="/home/archivos/b1.txt" -size=-45


#RECURSIVO
mkfile -r -path=/home/archivos/user/docs/usac/archivos/proyectos/fase1/entrada.txt


#-----------------6. CAT-----------------
cat -file1=/home/archivos/user/docs/Tarea2.txt
cat -file1=/home/archivos/user/docs/Tarea3.txt


#------------------------REPORTES PARTE 4----------------
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r1_inode.jpg -name=inode
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r2_block.pdf -name=block
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r3_bm_inode.txt -name=bm_inode
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r4_bm_block.txt -name=bm_block
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r5_sb.jpg -name=sb
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r6_file.jpg -path_file_ls=/home/archivos/user/docs/Tarea2.txt  -name=file
rep -id=761A -path=/home/fernando-vicente/Calificacion_MIA/Reportes/p4_r7_ls.jpg -path_file_ls=/home/archivos/user/docs -name=ls


#------------------------8. LOGOUT------------------------
logout
