mkdisk -size=10 -unit=M -fit=WF -path="discos/DiscoLab.mia"
mkdisk -size=1 -unit=K -fit=BF -path="discos/DiscoLab2.mia"

#comentario
rmdisk -path="discos/DiscoLab2.mia"

fdisk -size=1 -type=P -unit=M -fit=BF -name="Particion1" -path="discos/DiscoLab.mia"
fdisk -size=10 -type=P -unit=K -fit=WF -name="Particion2" -path="discos/DiscoLab.mia"

fdisk -size=10 -type=E -unit=K -fit=BF -name="Particion3" -path="discos/DiscoLab.mia"
fdisk -size=1 -type=L -unit=B -fit=WF -name="Particion4" -path="discos/DiscoLab.mia"
fdisk -size=1 -type=L -unit=K -fit=WF -name="Particion5" -path="discos/DiscoLab.mia"

fdisk -delete=fast -name="Particion5" -path="discos/DiscoLab.mia"
fdisk -delete=full -name="Particion4" -path="discos/DiscoLab.mia"

fdisk -add=10 -unit=K -path="discos/DiscoLab.mia" -name="Particion1"
fdisk -add=1 -unit=M -path="discos/DiscoLab.mia" -name="Particion2"


mount -name="Particion1" -path="discos/DiscoLab.mia"
mount -name="Particion2" -path="discos/DiscoLab.mia"

mount -name="Particion3" -path="discos/DiscoLab.mia"
mount -name="Particion4" -path="discos/DiscoLab.mia"

unmount -id=761A


mkfs -id=761A
mkfs -id=762A -type=full
mkfs -id=763A

mkdir -path="/home"
mkfile -size=15 -path=/home/user/docs/a.txt

rep -id=763A -path="salidas/report_mbr.png" -name=mbr
rep -id=762A -path="salidas/report_inode.png" -name=inode
rep -id=761A -path="salidas/report_bm_inode.txt" -name=bm_inode
rep -id=762A -path="salidas/report_disk.png" -name=disk
rep -id=763A -path="salidas/report_sb.png" -name=sb
rep -id=761A -path="salidas/report_bm_block.txt" -name=bm_block

login -user=root -pass=123 -id=761A
mkdir -path="/home"
mkfile -size=15 -path=/home/user/docs/a.txt
mkfile -path=/home/user/docs/a.txt -cont=/home/fernando-vicente/Documentos/b.txt

cat -file1=/home/user/docs/a.txt
logout
logout
