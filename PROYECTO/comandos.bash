mkdisk -size=5 -unit=M -fit=WF -path="discos/DiscoLab.mia"

fdisk -size=1 -type=P -unit=M -fit=BF -name="Particion1" -path="discos/DiscoLab.mia"
fdisk -size=10 -type=P -unit=K -fit=WF -name="Particion2" -path="discos/DiscoLab.mia"

mount -name="Particion1" -path="discos/DiscoLab.mia"
mount -name="Particion2" -path="discos/DiscoLab.mia"

mkfs -id=761A
mkfs -id=762A -type=full

rep -id=761A -path="salidas/report_mbr.png" -name=mbr
rep -id=761A -path="salidas/report_inode.png" -name=inode
rep -id=761A -path="salidas/report_bm_inode.txt" -name=bm_inode
rep -id=761A -path="salidas/report_disk.png" -name=disk